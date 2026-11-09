package api

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/golang/geo/r2"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

// positionsSpool stores the positions of the rounds that are over in a temporary file, so that only the positions of
// the current round are in memory while a demo is analyzed. Positions are created for every entity on every frame
// and make up most of the memory of an analysis when they are included, see AnalyzeAndExportDemo.
//
// The positions of a round are written per kind (players, grenades...) as a segment once the round is over, a very
// long round is written in several segments, see Analyzer.spillPositions(). Rounds may be discarded afterwards
// (match restart, round restart, incomplete round) so segments are only forgotten and their bytes left in the file,
// the file is truncated when the whole match is reset.
// Exporters read the segments of a kind back in writing order, a segment at a time, see Match.eachPlayerPositions().
type positionsSpool struct {
	file     *os.File
	writer   *bufio.Writer
	size     int64 // Bytes written to the file so far.
	segments []spoolSegment
	encoder  spoolEncoder
	tables   []byte // Scratch for the tables written after the records of a segment.
	err      error  // First write error, reported at the end of the analysis.
}

// spoolRecordsBlockSize is how many bytes of records are encoded before they are written to the file.
const spoolRecordsBlockSize = 64 << 10

type positionKind uint8

const (
	kindPlayerPositions positionKind = iota
	kindGrenadePositions
	kindInfernoPositions
	kindChickenPositions
	kindHostagePositions
)

type spoolSegment struct {
	kind          positionKind
	roundNumber   int
	offset        int64
	length        int
	recordsLength int // The records come first, the tables follow.
}

func newPositionsSpool() (*positionsSpool, error) {
	file, err := os.CreateTemp("", "cs-demo-analyzer-*.positions")
	if err != nil {
		return nil, fmt.Errorf("failed to create the positions temporary file: %w", err)
	}

	return &positionsSpool{
		file:   file,
		writer: bufio.NewWriterSize(file, 1<<20),
	}, nil
}

// close deletes the temporary file.
func (spool *positionsSpool) close() error {
	err := spool.file.Close()
	if removeErr := os.Remove(spool.file.Name()); err == nil {
		err = removeErr
	}

	return err
}

// writeRound writes the positions currently in the match as the positions of the round, the match positions must
// be emptied by the caller afterwards.
func (spool *positionsSpool) writeRound(roundNumber int, match *Match) {
	writeSpoolSegment(spool, kindPlayerPositions, roundNumber, match.PlayerPositions, encodePlayerPosition)
	writeSpoolSegment(spool, kindGrenadePositions, roundNumber, match.GrenadePositions, encodeGrenadePosition)
	writeSpoolSegment(spool, kindInfernoPositions, roundNumber, match.InfernoPositions, encodeInfernoPosition)
	writeSpoolSegment(spool, kindChickenPositions, roundNumber, match.ChickenPositions, encodeChickenPosition)
	writeSpoolSegment(spool, kindHostagePositions, roundNumber, match.HostagePositions, encodeHostagePosition)

	// Flushed after each round so that the file is complete whenever segments are read back.
	if err := spool.writer.Flush(); err != nil {
		spool.fail(err)
	}
}

// deleteRound forgets the segments of the round, like Match.resetRound() does with the positions in memory.
func (spool *positionsSpool) deleteRound(roundNumber int) {
	kept := spool.segments[:0]
	for _, segment := range spool.segments {
		if segment.roundNumber != roundNumber {
			kept = append(kept, segment)
		}
	}
	spool.segments = kept
}

// reset forgets all the segments and empties the file, like Match.reset() does with the positions in memory.
func (spool *positionsSpool) reset() {
	spool.segments = spool.segments[:0]
	spool.size = 0
	spool.writer.Reset(spool.file)
	if err := spool.file.Truncate(0); err != nil {
		spool.fail(err)
		return
	}
	if _, err := spool.file.Seek(0, io.SeekStart); err != nil {
		spool.fail(err)
	}
}

func (spool *positionsSpool) fail(err error) {
	if spool.err == nil {
		spool.err = fmt.Errorf("failed to write positions to the temporary file: %w", err)
	}
}

// writeSpoolSegment writes the items as a segment of the round: the count of items and their records, then the
// tables of the strings and weapon lists the records refer to. The records are written to the file by blocks as
// they are encoded so that the encoded segment is never in memory, the tables are complete only once all the
// records are encoded so they come last.
func writeSpoolSegment[T any](spool *positionsSpool, kind positionKind, roundNumber int, items []T, encode func(*spoolEncoder, T)) {
	if len(items) == 0 || spool.err != nil {
		return
	}

	encoder := &spool.encoder
	encoder.reset()
	recordsLength := 0
	writeRecords := func() bool {
		n, err := spool.writer.Write(encoder.records)
		recordsLength += n
		encoder.records = encoder.records[:0]
		if err != nil {
			spool.fail(err)
			return false
		}
		return true
	}

	encoder.records = binary.AppendUvarint(encoder.records, uint64(len(items)))
	for _, item := range items {
		encode(encoder, item)
		if len(encoder.records) >= spoolRecordsBlockSize && !writeRecords() {
			return
		}
	}
	if !writeRecords() {
		return
	}

	tables := spool.tables[:0]
	tables = binary.AppendUvarint(tables, uint64(len(encoder.stringList)))
	for _, s := range encoder.stringList {
		tables = binary.AppendUvarint(tables, uint64(len(s)))
		tables = append(tables, s...)
	}
	tables = binary.AppendUvarint(tables, uint64(len(encoder.weaponsList)))
	for _, weapons := range encoder.weaponsList {
		for _, list := range weapons.lists() {
			tables = appendSliceLength(tables, list == nil, len(list))
			for _, name := range list {
				tables = binary.AppendUvarint(tables, uint64(encoder.strings[string(name)]))
			}
		}
	}
	spool.tables = tables

	if _, err := spool.writer.Write(tables); err != nil {
		spool.fail(err)
		return
	}

	length := recordsLength + len(tables)
	spool.segments = append(spool.segments, spoolSegment{
		kind:          kind,
		roundNumber:   roundNumber,
		offset:        spool.size,
		length:        length,
		recordsLength: recordsLength,
	})
	spool.size += int64(length)
}

// readSpoolSegments calls fn with the decoded items of each segment of the kind, in writing order (a round or a
// part of a round at a time). The spool must not be written to anymore.
func readSpoolSegments[T any](spool *positionsSpool, kind positionKind, decode func(*spoolDecoder) T, fn func([]T) error) error {
	if spool.err != nil {
		return spool.err
	}

	var buffer []byte
	for _, segment := range spool.segments {
		if segment.kind != kind {
			continue
		}

		if cap(buffer) < segment.length {
			buffer = make([]byte, segment.length)
		}
		buffer = buffer[:segment.length]
		if _, err := spool.file.ReadAt(buffer, segment.offset); err != nil {
			return fmt.Errorf("failed to read positions from the temporary file: %w", err)
		}

		decoder := newSpoolDecoder(buffer[:segment.recordsLength], buffer[segment.recordsLength:])
		count := int(decoder.uvarint())
		if decoder.err != nil || count < 0 || count > segment.length {
			return errCorruptedSpool
		}
		items := make([]T, 0, count)
		for i := 0; i < count && decoder.err == nil; i++ {
			items = append(items, decode(decoder))
		}
		if decoder.err != nil {
			return decoder.err
		}

		if err := fn(items); err != nil {
			return err
		}
	}

	return nil
}

// appendSliceLength encodes a slice length keeping nil apart from empty, nil and empty slices are marshalled
// differently to JSON.
func appendSliceLength(buffer []byte, isNil bool, length int) []byte {
	if isNil {
		return binary.AppendUvarint(buffer, 0)
	}

	return binary.AppendUvarint(buffer, uint64(length)+1)
}

// lists returns the weapon lists in the order they are encoded.
func (w *PlayerWeapons) lists() [6][]constants.WeaponName {
	return [6][]constants.WeaponName{w.Equipments, w.Grenades, w.Pistols, w.SMGs, w.Rifles, w.Heavy}
}

// spoolEncoder encodes the records of a segment. Strings and weapon lists are replaced by their index in tables
// written before the records: names are repeated in every position and the weapon lists are shared by the
// consecutive positions of a player.
type spoolEncoder struct {
	records     []byte
	strings     map[string]uint32
	stringList  []string
	weapons     map[*PlayerWeapons]uint32
	weaponsList []*PlayerWeapons
}

func (encoder *spoolEncoder) reset() {
	encoder.records = encoder.records[:0]
	encoder.strings = make(map[string]uint32)
	encoder.stringList = encoder.stringList[:0]
	encoder.weapons = make(map[*PlayerWeapons]uint32)
	encoder.weaponsList = encoder.weaponsList[:0]
}

func (encoder *spoolEncoder) internString(s string) uint32 {
	index, exists := encoder.strings[s]
	if !exists {
		index = uint32(len(encoder.stringList))
		encoder.strings[s] = index
		encoder.stringList = append(encoder.stringList, s)
	}

	return index
}

func (encoder *spoolEncoder) internWeapons(weapons *PlayerWeapons) uint32 {
	index, exists := encoder.weapons[weapons]
	if !exists {
		index = uint32(len(encoder.weaponsList))
		encoder.weapons[weapons] = index
		encoder.weaponsList = append(encoder.weaponsList, weapons)
		// The names are interned now so that the string table is complete when the weapons table is written.
		for _, list := range weapons.lists() {
			for _, name := range list {
				encoder.internString(string(name))
			}
		}
	}

	return index
}

func (encoder *spoolEncoder) int(value int) {
	encoder.records = binary.AppendVarint(encoder.records, int64(value))
}

func (encoder *spoolEncoder) uint64(value uint64) {
	encoder.records = binary.LittleEndian.AppendUint64(encoder.records, value)
}

func (encoder *spoolEncoder) float64(value float64) {
	encoder.records = binary.LittleEndian.AppendUint64(encoder.records, math.Float64bits(value))
}

func (encoder *spoolEncoder) float32(value float32) {
	encoder.records = binary.LittleEndian.AppendUint32(encoder.records, math.Float32bits(value))
}

func (encoder *spoolEncoder) byte(value byte) {
	encoder.records = append(encoder.records, value)
}

func (encoder *spoolEncoder) bool(value bool) {
	if value {
		encoder.records = append(encoder.records, 1)
	} else {
		encoder.records = append(encoder.records, 0)
	}
}

func (encoder *spoolEncoder) string(value string) {
	encoder.records = binary.AppendUvarint(encoder.records, uint64(encoder.internString(value)))
}

func (encoder *spoolEncoder) playerWeapons(weapons *PlayerWeapons) {
	encoder.records = binary.AppendUvarint(encoder.records, uint64(encoder.internWeapons(weapons)))
}

func (encoder *spoolEncoder) points(points []r2.Point) {
	encoder.records = appendSliceLength(encoder.records, points == nil, len(points))
	for _, point := range points {
		encoder.float64(point.X)
		encoder.float64(point.Y)
	}
}

var errCorruptedSpool = errors.New("corrupted positions temporary file")

// spoolDecoder decodes a segment, the tables are read when it's created and the records are then decoded from data.
// Decoding errors are sticky, err must be checked once the records are decoded.
type spoolDecoder struct {
	data        []byte
	stringList  []string
	weaponsList []*PlayerWeapons
	err         error
}

func newSpoolDecoder(records []byte, tables []byte) *spoolDecoder {
	decoder := &spoolDecoder{data: tables}
	defer func() {
		decoder.data = records
	}()

	stringCount := int(decoder.uvarint())
	if decoder.err != nil || stringCount < 0 || stringCount > len(tables) {
		decoder.err = errCorruptedSpool
		return decoder
	}
	decoder.stringList = make([]string, 0, stringCount)
	for i := 0; i < stringCount && decoder.err == nil; i++ {
		length := int(decoder.uvarint())
		if length < 0 || length > len(decoder.data) {
			decoder.err = errCorruptedSpool
			break
		}
		decoder.stringList = append(decoder.stringList, string(decoder.data[:length]))
		decoder.data = decoder.data[length:]
	}

	weaponsCount := int(decoder.uvarint())
	if decoder.err != nil || weaponsCount < 0 || weaponsCount > len(tables) {
		decoder.err = errCorruptedSpool
		return decoder
	}
	decoder.weaponsList = make([]*PlayerWeapons, 0, weaponsCount)
	for i := 0; i < weaponsCount && decoder.err == nil; i++ {
		decoder.weaponsList = append(decoder.weaponsList, &PlayerWeapons{
			Equipments: decoder.weaponNames(),
			Grenades:   decoder.weaponNames(),
			Pistols:    decoder.weaponNames(),
			SMGs:       decoder.weaponNames(),
			Rifles:     decoder.weaponNames(),
			Heavy:      decoder.weaponNames(),
		})
	}

	return decoder
}

func (decoder *spoolDecoder) uvarint() uint64 {
	value, n := binary.Uvarint(decoder.data)
	if n <= 0 {
		decoder.err = errCorruptedSpool
		return 0
	}
	decoder.data = decoder.data[n:]

	return value
}

func (decoder *spoolDecoder) int() int {
	value, n := binary.Varint(decoder.data)
	if n <= 0 {
		decoder.err = errCorruptedSpool
		return 0
	}
	decoder.data = decoder.data[n:]

	return int(value)
}

func (decoder *spoolDecoder) bytes(n int) []byte {
	if len(decoder.data) < n {
		decoder.err = errCorruptedSpool
		return make([]byte, n)
	}
	bytes := decoder.data[:n]
	decoder.data = decoder.data[n:]

	return bytes
}

func (decoder *spoolDecoder) uint64() uint64 {
	return binary.LittleEndian.Uint64(decoder.bytes(8))
}

func (decoder *spoolDecoder) float64() float64 {
	return math.Float64frombits(decoder.uint64())
}

func (decoder *spoolDecoder) float32() float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(decoder.bytes(4)))
}

func (decoder *spoolDecoder) byte() byte {
	return decoder.bytes(1)[0]
}

func (decoder *spoolDecoder) bool() bool {
	return decoder.byte() != 0
}

func (decoder *spoolDecoder) string() string {
	index := decoder.uvarint()
	if index >= uint64(len(decoder.stringList)) {
		decoder.err = errCorruptedSpool
		return ""
	}

	return decoder.stringList[index]
}

func (decoder *spoolDecoder) playerWeapons() *PlayerWeapons {
	index := decoder.uvarint()
	if index >= uint64(len(decoder.weaponsList)) {
		decoder.err = errCorruptedSpool
		return &PlayerWeapons{}
	}

	return decoder.weaponsList[index]
}

// sliceLength decodes a length written by appendSliceLength, -1 for a nil slice.
func (decoder *spoolDecoder) sliceLength() int {
	length := decoder.uvarint()
	if length == 0 {
		return -1
	}
	if length-1 > uint64(len(decoder.data)) {
		decoder.err = errCorruptedSpool
		return -1
	}

	return int(length - 1)
}

func (decoder *spoolDecoder) weaponNames() []constants.WeaponName {
	length := decoder.sliceLength()
	if length < 0 {
		return nil
	}
	names := make([]constants.WeaponName, length)
	for i := range names {
		names[i] = constants.WeaponName(decoder.string())
	}

	return names
}

func (decoder *spoolDecoder) points() []r2.Point {
	length := decoder.sliceLength()
	if length < 0 {
		return nil
	}
	points := make([]r2.Point, length)
	for i := range points {
		points[i].X = decoder.float64()
		points[i].Y = decoder.float64()
	}

	return points
}

func encodePlayerPosition(encoder *spoolEncoder, position *PlayerPosition) {
	encoder.int(position.Tick)
	encoder.int(position.RoundNumber)
	encoder.bool(position.IsAlive)
	encoder.string(position.Name)
	encoder.uint64(position.SteamID64)
	encoder.float64(position.X)
	encoder.float64(position.Y)
	encoder.float64(position.Z)
	encoder.float32(position.Yaw)
	encoder.float32(position.Pitch)
	encoder.float64(position.FlashDurationRemaining)
	encoder.byte(byte(position.Side))
	encoder.int(position.Money)
	encoder.int(position.Health)
	encoder.int(position.Armor)
	encoder.bool(position.HasHelmet)
	encoder.bool(position.HasBomb)
	encoder.bool(position.HasDefuseKit)
	encoder.bool(position.IsDucking)
	encoder.bool(position.IsAirborne)
	encoder.bool(position.IsScoping)
	encoder.bool(position.IsDefusing)
	encoder.bool(position.IsPlanting)
	encoder.bool(position.IsGrabbingHostage)
	encoder.string(string(position.ActiveWeaponName))
	encoder.playerWeapons(position.PlayerWeapons)
	encoder.int(position.Ping)
}

func decodePlayerPosition(decoder *spoolDecoder) *PlayerPosition {
	return &PlayerPosition{
		Tick:                   decoder.int(),
		RoundNumber:            decoder.int(),
		IsAlive:                decoder.bool(),
		Name:                   decoder.string(),
		SteamID64:              decoder.uint64(),
		X:                      decoder.float64(),
		Y:                      decoder.float64(),
		Z:                      decoder.float64(),
		Yaw:                    decoder.float32(),
		Pitch:                  decoder.float32(),
		FlashDurationRemaining: decoder.float64(),
		Side:                   common.Team(decoder.byte()),
		Money:                  decoder.int(),
		Health:                 decoder.int(),
		Armor:                  decoder.int(),
		HasHelmet:              decoder.bool(),
		HasBomb:                decoder.bool(),
		HasDefuseKit:           decoder.bool(),
		IsDucking:              decoder.bool(),
		IsAirborne:             decoder.bool(),
		IsScoping:              decoder.bool(),
		IsDefusing:             decoder.bool(),
		IsPlanting:             decoder.bool(),
		IsGrabbingHostage:      decoder.bool(),
		ActiveWeaponName:       constants.WeaponName(decoder.string()),
		PlayerWeapons:          decoder.playerWeapons(),
		Ping:                   decoder.int(),
	}
}

func encodeGrenadePosition(encoder *spoolEncoder, position *GrenadePosition) {
	encoder.int(position.Tick)
	encoder.int(position.RoundNumber)
	encoder.string(position.GrenadeID)
	encoder.string(position.ProjectileID)
	encoder.float64(position.X)
	encoder.float64(position.Y)
	encoder.float64(position.Z)
	encoder.uint64(position.ThrowerSteamID64)
	encoder.string(position.ThrowerName)
	encoder.byte(byte(position.ThrowerSide))
	encoder.string(position.ThrowerTeamName)
	encoder.float64(position.ThrowerVelocityX)
	encoder.float64(position.ThrowerVelocityY)
	encoder.float64(position.ThrowerVelocityZ)
	encoder.float32(position.ThrowerPitch)
	encoder.float32(position.ThrowerYaw)
	encoder.string(string(position.GrenadeName))
}

func decodeGrenadePosition(decoder *spoolDecoder) *GrenadePosition {
	return &GrenadePosition{
		Tick:             decoder.int(),
		RoundNumber:      decoder.int(),
		GrenadeID:        decoder.string(),
		ProjectileID:     decoder.string(),
		X:                decoder.float64(),
		Y:                decoder.float64(),
		Z:                decoder.float64(),
		ThrowerSteamID64: decoder.uint64(),
		ThrowerName:      decoder.string(),
		ThrowerSide:      common.Team(decoder.byte()),
		ThrowerTeamName:  decoder.string(),
		ThrowerVelocityX: decoder.float64(),
		ThrowerVelocityY: decoder.float64(),
		ThrowerVelocityZ: decoder.float64(),
		ThrowerPitch:     decoder.float32(),
		ThrowerYaw:       decoder.float32(),
		GrenadeName:      constants.WeaponName(decoder.string()),
	}
}

func encodeInfernoPosition(encoder *spoolEncoder, position *InfernoPosition) {
	encoder.int(position.Tick)
	encoder.int(position.RoundNumber)
	encoder.uint64(position.ThrowerSteamID64)
	encoder.string(position.ThrowerName)
	encoder.string(position.UniqueID)
	encoder.float64(position.X)
	encoder.float64(position.Y)
	encoder.float64(position.Z)
	encoder.points(position.ConvexHull2D)
}

func decodeInfernoPosition(decoder *spoolDecoder) *InfernoPosition {
	return &InfernoPosition{
		Tick:             decoder.int(),
		RoundNumber:      decoder.int(),
		ThrowerSteamID64: decoder.uint64(),
		ThrowerName:      decoder.string(),
		UniqueID:         decoder.string(),
		X:                decoder.float64(),
		Y:                decoder.float64(),
		Z:                decoder.float64(),
		ConvexHull2D:     decoder.points(),
	}
}

func encodeChickenPosition(encoder *spoolEncoder, position *ChickenPosition) {
	encoder.int(position.chickenID)
	encoder.int(position.Tick)
	encoder.int(position.RoundNumber)
	encoder.float64(position.X)
	encoder.float64(position.Y)
	encoder.float64(position.Z)
}

func decodeChickenPosition(decoder *spoolDecoder) *ChickenPosition {
	return &ChickenPosition{
		chickenID:   decoder.int(),
		Tick:        decoder.int(),
		RoundNumber: decoder.int(),
		X:           decoder.float64(),
		Y:           decoder.float64(),
		Z:           decoder.float64(),
	}
}

func encodeHostagePosition(encoder *spoolEncoder, position *HostagePosition) {
	encoder.int(position.entityID)
	encoder.int(position.serialNum)
	encoder.int(position.Tick)
	encoder.int(position.RoundNumber)
	encoder.float64(position.X)
	encoder.float64(position.Y)
	encoder.float64(position.Z)
	encoder.byte(byte(position.State))
}

func decodeHostagePosition(decoder *spoolDecoder) *HostagePosition {
	return &HostagePosition{
		entityID:    decoder.int(),
		serialNum:   decoder.int(),
		Tick:        decoder.int(),
		RoundNumber: decoder.int(),
		X:           decoder.float64(),
		Y:           decoder.float64(),
		Z:           decoder.float64(),
		State:       common.HostageState(decoder.byte()),
	}
}

// The each*Positions methods call fn with the positions of the match in frame order, a chunk at a time. A chunk
// holds whole rounds: the positions of a round are never split across chunks. There is a chunk per round when the
// positions were spooled during the analysis (see positionsSpool), otherwise a single chunk with the positions in
// memory. Exporters go through positions this way so that they never need them all in memory.

func (match *Match) eachPlayerPositions(fn func(positions []*PlayerPosition) error) error {
	return eachPositions(match, kindPlayerPositions, decodePlayerPosition, match.PlayerPositions, fn)
}

func (match *Match) eachGrenadePositions(fn func(positions []*GrenadePosition) error) error {
	return eachPositions(match, kindGrenadePositions, decodeGrenadePosition, match.GrenadePositions, fn)
}

func (match *Match) eachInfernoPositions(fn func(positions []*InfernoPosition) error) error {
	return eachPositions(match, kindInfernoPositions, decodeInfernoPosition, match.InfernoPositions, fn)
}

func (match *Match) eachChickenPositions(fn func(positions []*ChickenPosition) error) error {
	return eachPositions(match, kindChickenPositions, decodeChickenPosition, match.ChickenPositions, fn)
}

func (match *Match) eachHostagePositions(fn func(positions []*HostagePosition) error) error {
	return eachPositions(match, kindHostagePositions, decodeHostagePosition, match.HostagePositions, fn)
}

func eachPositions[T any](match *Match, kind positionKind, decode func(*spoolDecoder) T, inMemory []T, fn func([]T) error) error {
	if match.positionsSpool != nil {
		if err := readSpoolSegments(match.positionsSpool, kind, decode, fn); err != nil {
			return err
		}
	}

	// The positions still in memory are the ones of the rounds that haven't been spooled, they come after the spooled
	// ones.
	if len(inMemory) > 0 {
		return fn(inMemory)
	}

	return nil
}

// closePositionsSpool deletes the temporary file of the positions spooled during the analysis, if any.
func (match *Match) closePositionsSpool() error {
	if match.positionsSpool == nil {
		return nil
	}

	err := match.positionsSpool.close()
	match.positionsSpool = nil

	return err
}
