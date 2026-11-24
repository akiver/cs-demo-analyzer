package source2

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/klauspost/compress/zstd"
)

// Binary KeyValues3 decoder.
// Reference: https://github.com/ValveResourceFormat/ValveResourceFormat/blob/master/ValveResourceFormat/Resource/ResourceTypes/BinaryKV3.cs

const (
	kv3Magic         = 0x4b563300
	kv3BlockTrailer  = 0xffeedd00
	kv3MinVersion    = 1
	kv3MaxVersion    = 5
	kv3GUIDSize      = 16
	compressionNone  = 0
	compressionLZ4   = 1
	compressionZstd  = 2
	lz4MinMatch      = 4
	kv3FlagByteMask  = 0x80
	kv3ExtraByteMask = 0x40
)

type kv3NodeType byte

const (
	kv3Null                   kv3NodeType = 1
	kv3Boolean                kv3NodeType = 2
	kv3Int64                  kv3NodeType = 3
	kv3UInt64                 kv3NodeType = 4
	kv3Double                 kv3NodeType = 5
	kv3String                 kv3NodeType = 6
	kv3BinaryBlob             kv3NodeType = 7
	kv3Array                  kv3NodeType = 8
	kv3Object                 kv3NodeType = 9
	kv3ArrayTyped             kv3NodeType = 10
	kv3Int32                  kv3NodeType = 11
	kv3UInt32                 kv3NodeType = 12
	kv3BooleanTrue            kv3NodeType = 13
	kv3BooleanFalse           kv3NodeType = 14
	kv3Int64Zero              kv3NodeType = 15
	kv3Int64One               kv3NodeType = 16
	kv3DoubleZero             kv3NodeType = 17
	kv3DoubleOne              kv3NodeType = 18
	kv3Float                  kv3NodeType = 19
	kv3Int16                  kv3NodeType = 20
	kv3UInt16                 kv3NodeType = 21
	kv3Int8                   kv3NodeType = 22
	kv3UInt8                  kv3NodeType = 23
	kv3ArrayTypeByteLength    kv3NodeType = 24
	kv3ArrayTypeAuxiliaryBuff kv3NodeType = 25
	kv3UInt32Zero             kv3NodeType = 26
	kv3UInt32One              kv3NodeType = 27
)

var errKV3OutOfData = errors.New("kv3: unexpected end of data")

type kv3Header struct {
	version                       int
	compressionMethod             uint32
	compressionFrameSize          uint16
	countBytes1                   int
	countBytes2                   int
	countBytes4                   int
	countBytes8                   int
	countTypes                    int
	sizeUncompressedTotal         int
	sizeCompressedTotal           int
	countBlocks                   int
	sizeBinaryBlobsBytes          int
	sizeBlockCompressedSizesBytes int
	sizeUncompressedBuffer1       int
	sizeCompressedBuffer1         int
	sizeUncompressedBuffer2       int
	sizeCompressedBuffer2         int
	countBytes1Buffer2            int
	countBytes2Buffer2            int
	countBytes4Buffer2            int
	countBytes8Buffer2            int
	countObjectsBuffer2           int
}

// Before version 5, zstd stores the buffer and the binary blobs as two consecutive frames decompressed together.
func (header kv3Header) sizeDecompressedBuffer1() int {
	if header.version < 5 && header.compressionMethod == compressionZstd {
		return header.sizeUncompressedBuffer1 + header.sizeBinaryBlobsBytes
	}

	return header.sizeUncompressedBuffer1
}

type kv3Lanes struct {
	bytes1 []byte
	bytes2 []byte
	bytes4 []byte
	bytes8 []byte
}

type kv3Context struct {
	version         int
	types           []byte
	objectLengths   []byte
	binaryBlobs     [][]byte
	nextBinaryBlob  int
	strings         []string
	buffer          *kv3Lanes
	auxiliaryBuffer *kv3Lanes
}

type byteReader struct {
	data     []byte
	position int
	err      error
}

func (reader *byteReader) take(count int) []byte {
	if reader.err != nil {
		return nil
	}
	if count < 0 || reader.position+count > len(reader.data) {
		reader.err = errKV3OutOfData
		return nil
	}
	value := reader.data[reader.position : reader.position+count]
	reader.position += count

	return value
}

func (reader *byteReader) uint16() uint16 {
	bytes := reader.take(2)
	if bytes == nil {
		return 0
	}

	return binary.LittleEndian.Uint16(bytes)
}

func (reader *byteReader) uint32() uint32 {
	bytes := reader.take(4)
	if bytes == nil {
		return 0
	}

	return binary.LittleEndian.Uint32(bytes)
}

func (reader *byteReader) int32() int {
	return int(int32(reader.uint32()))
}

// DecodeKV3 decodes a binary KeyValues3 block and returns its root value.
func DecodeKV3(data []byte) (any, error) {
	reader := &byteReader{data: data}
	magic := reader.uint32()
	if reader.err != nil {
		return nil, reader.err
	}

	version := int(magic & 0xff)
	if magic&0xffffff00 != kv3Magic {
		return nil, fmt.Errorf("kv3: unsupported signature 0x%08x", magic)
	}
	if version < kv3MinVersion || version > kv3MaxVersion {
		return nil, fmt.Errorf("kv3: unsupported version %d", version)
	}

	header := readKV3Header(reader, version, len(data))
	if reader.err != nil {
		return nil, fmt.Errorf("kv3: read header: %w", reader.err)
	}

	return readKV3Buffers(reader, header)
}

func readKV3Header(reader *byteReader, version int, blockEnd int) kv3Header {
	header := kv3Header{version: version}
	reader.take(kv3GUIDSize)
	header.compressionMethod = reader.uint32()

	if version == 1 {
		header.countBytes1 = reader.int32()
		header.countBytes4 = reader.int32()
		header.countBytes8 = reader.int32()
		header.sizeUncompressedTotal = reader.int32()
		header.sizeCompressedTotal = blockEnd - reader.position
	} else {
		compressionDictionaryID := reader.uint16()
		if compressionDictionaryID != 0 && reader.err == nil {
			reader.err = fmt.Errorf("unsupported compression dictionary %d", compressionDictionaryID)
		}
		header.compressionFrameSize = reader.uint16()
		header.countBytes1 = reader.int32()
		header.countBytes4 = reader.int32()
		header.countBytes8 = reader.int32()
		header.countTypes = reader.int32()
		reader.uint16() // countObjects
		reader.uint16() // countArrays
		header.sizeUncompressedTotal = reader.int32()
		header.sizeCompressedTotal = reader.int32()
		header.countBlocks = reader.int32()
		header.sizeBinaryBlobsBytes = reader.int32()
	}

	if version >= 4 {
		header.countBytes2 = reader.int32()
		header.sizeBlockCompressedSizesBytes = reader.int32()
	}

	if version >= 5 {
		header.sizeUncompressedBuffer1 = reader.int32()
		header.sizeCompressedBuffer1 = reader.int32()
		header.sizeUncompressedBuffer2 = reader.int32()
		header.sizeCompressedBuffer2 = reader.int32()
		header.countBytes1Buffer2 = reader.int32()
		header.countBytes2Buffer2 = reader.int32()
		header.countBytes4Buffer2 = reader.int32()
		header.countBytes8Buffer2 = reader.int32()
		reader.int32() // countNodes
		header.countObjectsBuffer2 = reader.int32()
		reader.int32() // countArrays
		reader.int32() // countArrayElements
	} else {
		header.sizeCompressedBuffer1 = header.sizeCompressedTotal
		header.sizeUncompressedBuffer1 = header.sizeUncompressedTotal
	}

	return header
}

func readKV3BufferData(reader *byteReader, compressionMethod uint32, size int, compressedSize int) ([]byte, error) {
	switch compressionMethod {
	case compressionNone:
		data := reader.take(size)
		if reader.err != nil {
			return nil, reader.err
		}

		return data, nil
	case compressionLZ4:
		input := reader.take(compressedSize)
		if reader.err != nil {
			return nil, reader.err
		}
		output := make([]byte, 0, size)
		output, err := decodeLZ4Block(output, input)
		if err != nil {
			return nil, err
		}
		if len(output) != size {
			return nil, fmt.Errorf("kv3: LZ4 decompressed %d bytes, expected %d", len(output), size)
		}

		return output, nil
	case compressionZstd:
		input := reader.take(compressedSize)
		if reader.err != nil {
			return nil, reader.err
		}

		return decodeZstd(input, size)
	default:
		return nil, fmt.Errorf("kv3: unknown compression method %d", compressionMethod)
	}
}

func decodeZstd(input []byte, size int) ([]byte, error) {
	if len(input) == 0 && size == 0 {
		return []byte{}, nil
	}

	decoder, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(1), zstd.WithDecoderMaxMemory(1<<31))
	if err != nil {
		return nil, err
	}
	defer decoder.Close()

	output, err := decoder.DecodeAll(input, make([]byte, 0, size))
	if err != nil {
		return nil, fmt.Errorf("kv3: zstd: %w", err)
	}
	if len(output) != size {
		return nil, fmt.Errorf("kv3: zstd decompressed %d bytes, expected %d", len(output), size)
	}

	return output, nil
}

// decodeLZ4Block decodes a raw LZ4 block and appends it to dst.
// Matches may reference bytes already present in dst, which is how chained LZ4 frames are decoded.
func decodeLZ4Block(dst []byte, src []byte) ([]byte, error) {
	position := 0
	readLength := func(length int) (int, error) {
		if length != 15 {
			return length, nil
		}
		for {
			if position >= len(src) {
				return 0, errKV3OutOfData
			}
			value := src[position]
			position++
			length += int(value)
			if value != 255 {
				return length, nil
			}
		}
	}

	for position < len(src) {
		token := src[position]
		position++

		literalLength, err := readLength(int(token >> 4))
		if err != nil {
			return nil, err
		}
		if position+literalLength > len(src) {
			return nil, errKV3OutOfData
		}
		dst = append(dst, src[position:position+literalLength]...)
		position += literalLength

		if position >= len(src) {
			break
		}

		if position+2 > len(src) {
			return nil, errKV3OutOfData
		}
		offset := int(binary.LittleEndian.Uint16(src[position:]))
		position += 2
		if offset == 0 || offset > len(dst) {
			return nil, fmt.Errorf("kv3: invalid LZ4 match offset %d", offset)
		}

		matchLength, err := readLength(int(token & 0x0f))
		if err != nil {
			return nil, err
		}
		matchLength += lz4MinMatch

		start := len(dst) - offset
		for i := range matchLength {
			dst = append(dst, dst[start+i])
		}
	}

	return dst, nil
}

func align(offset int, alignment int) int {
	return (offset + alignment - 1) &^ (alignment - 1)
}

func slice(buffer []byte, start int, end int) ([]byte, error) {
	if start < 0 || end > len(buffer) || start > end {
		return nil, errKV3OutOfData
	}

	return buffer[start:end], nil
}

func readKV3Buffers(reader *byteReader, header kv3Header) (any, error) {
	context := &kv3Context{version: header.version}

	buffer1Raw, err := readKV3BufferData(reader, header.compressionMethod, header.sizeDecompressedBuffer1(), header.sizeCompressedBuffer1)
	if err != nil {
		return nil, fmt.Errorf("kv3: buffer 1: %w", err)
	}
	if len(buffer1Raw) < header.sizeUncompressedBuffer1 {
		return nil, errKV3OutOfData
	}
	buffer1Data := buffer1Raw[:header.sizeUncompressedBuffer1]

	var bufferWithBinaryBlobSizes []byte

	buffer1 := &kv3Lanes{}
	offset := 0
	if header.countBytes1 > 0 {
		if buffer1.bytes1, err = slice(buffer1Data, offset, offset+header.countBytes1); err != nil {
			return nil, err
		}
		offset += header.countBytes1
	}
	if header.countBytes2 > 0 {
		offset = align(offset, 2)
		if buffer1.bytes2, err = slice(buffer1Data, offset, offset+header.countBytes2*2); err != nil {
			return nil, err
		}
		offset += header.countBytes2 * 2
	}
	if header.countBytes4 > 0 {
		offset = align(offset, 4)
		if buffer1.bytes4, err = slice(buffer1Data, offset, offset+header.countBytes4*4); err != nil {
			return nil, err
		}
		offset += header.countBytes4 * 4
	}
	if header.countBytes8 > 0 {
		offset = align(offset, 8)
		if buffer1.bytes8, err = slice(buffer1Data, offset, offset+header.countBytes8*8); err != nil {
			return nil, err
		}
		offset += header.countBytes8 * 8
	} else if header.version < 5 {
		offset = align(offset, 8)
	}

	if len(buffer1.bytes4) < 4 {
		return nil, fmt.Errorf("kv3: missing strings count")
	}
	countStrings := int(int32(binary.LittleEndian.Uint32(buffer1.bytes4)))
	buffer1.bytes4 = buffer1.bytes4[4:]
	if countStrings < 0 {
		return nil, fmt.Errorf("kv3: invalid strings count %d", countStrings)
	}
	context.auxiliaryBuffer = buffer1

	if header.version >= 5 {
		context.strings, buffer1.bytes1, err = readKV3Strings(buffer1.bytes1, countStrings)
		if err != nil {
			return nil, err
		}
	} else {
		context.buffer = buffer1
		stringsStartOffset := offset
		var remaining []byte
		context.strings, remaining, err = readKV3Strings(buffer1Data[offset:], countStrings)
		if err != nil {
			return nil, err
		}
		offset = len(buffer1Data) - len(remaining)

		var typesLength int
		if header.version == 1 {
			typesLength = header.sizeUncompressedTotal - offset - 4
		} else {
			typesLength = header.countTypes - offset + stringsStartOffset
		}
		if context.types, err = slice(buffer1Data, offset, offset+typesLength); err != nil {
			return nil, err
		}
		offset += typesLength

		if header.countBlocks == 0 {
			if offset+4 > len(buffer1Data) || binary.LittleEndian.Uint32(buffer1Data[offset:]) != kv3BlockTrailer {
				return nil, fmt.Errorf("kv3: invalid buffer trailer")
			}
		} else {
			bufferWithBinaryBlobSizes = buffer1Data[offset:]
		}
	}

	if header.version >= 5 {
		buffer2Data, err := readKV3BufferData(reader, header.compressionMethod, header.sizeUncompressedBuffer2, header.sizeCompressedBuffer2)
		if err != nil {
			return nil, fmt.Errorf("kv3: buffer 2: %w", err)
		}

		buffer2 := &kv3Lanes{}
		context.buffer = buffer2

		end := header.countObjectsBuffer2 * 4
		if context.objectLengths, err = slice(buffer2Data, 0, end); err != nil {
			return nil, err
		}
		offset := end
		if header.countBytes1Buffer2 > 0 {
			if buffer2.bytes1, err = slice(buffer2Data, offset, offset+header.countBytes1Buffer2); err != nil {
				return nil, err
			}
			offset += header.countBytes1Buffer2
		}
		if header.countBytes2Buffer2 > 0 {
			offset = align(offset, 2)
			if buffer2.bytes2, err = slice(buffer2Data, offset, offset+header.countBytes2Buffer2*2); err != nil {
				return nil, err
			}
			offset += header.countBytes2Buffer2 * 2
		}
		if header.countBytes4Buffer2 > 0 {
			offset = align(offset, 4)
			if buffer2.bytes4, err = slice(buffer2Data, offset, offset+header.countBytes4Buffer2*4); err != nil {
				return nil, err
			}
			offset += header.countBytes4Buffer2 * 4
		}
		if header.countBytes8Buffer2 > 0 {
			offset = align(offset, 8)
			if buffer2.bytes8, err = slice(buffer2Data, offset, offset+header.countBytes8Buffer2*8); err != nil {
				return nil, err
			}
			offset += header.countBytes8Buffer2 * 8
		}

		if context.types, err = slice(buffer2Data, offset, offset+header.countTypes); err != nil {
			return nil, err
		}
		offset += header.countTypes

		if header.countBlocks == 0 {
			if offset+4 > len(buffer2Data) || binary.LittleEndian.Uint32(buffer2Data[offset:]) != kv3BlockTrailer {
				return nil, fmt.Errorf("kv3: invalid buffer 2 trailer")
			}
		} else {
			bufferWithBinaryBlobSizes = buffer2Data[offset:]
		}
	}

	if header.countBlocks > 0 {
		if err := readKV3BinaryBlobs(reader, header, context, bufferWithBinaryBlobSizes, buffer1Raw); err != nil {
			return nil, err
		}
	}

	rootType, err := context.readType()
	if err != nil {
		return nil, err
	}

	return context.readValue(rootType, context.buffer)
}

func readKV3Strings(buffer []byte, count int) ([]string, []byte, error) {
	strings := make([]string, count)
	position := 0
	for i := range count {
		end := position
		for end < len(buffer) && buffer[end] != 0 {
			end++
		}
		if end >= len(buffer) {
			return nil, nil, fmt.Errorf("kv3: unterminated string")
		}
		strings[i] = string(buffer[position:end])
		position = end + 1
	}

	return strings, buffer[position:], nil
}

func readKV3BinaryBlobs(reader *byteReader, header kv3Header, context *kv3Context, sizes []byte, buffer1Raw []byte) error {
	end := header.countBlocks * 4
	if end+4 > len(sizes) {
		return errKV3OutOfData
	}

	blobLengths := make([]int, header.countBlocks)
	total := 0
	for i := range blobLengths {
		blobLengths[i] = int(int32(binary.LittleEndian.Uint32(sizes[i*4:])))
		if blobLengths[i] < 0 {
			return fmt.Errorf("kv3: invalid binary blob length %d", blobLengths[i])
		}
		total += blobLengths[i]
	}
	if total != header.sizeBinaryBlobsBytes {
		return fmt.Errorf("kv3: binary blob lengths %d do not match total %d", total, header.sizeBinaryBlobsBytes)
	}
	if binary.LittleEndian.Uint32(sizes[end:]) != kv3BlockTrailer {
		return fmt.Errorf("kv3: invalid binary blob sizes trailer")
	}
	sizes = sizes[end+4:]

	var blobsData []byte
	switch header.compressionMethod {
	case compressionNone:
		blobsData = reader.take(total)
		if reader.err != nil {
			return reader.err
		}
	case compressionLZ4:
		// Each blob is split into frames of up to the compression frame size. Frames are chained, so decoding them
		// into one contiguous buffer lets matches reference the previous blobs.
		blobsData = make([]byte, 0, total)
		for _, length := range blobLengths {
			blobEnd := len(blobsData) + length
			for len(blobsData) < blobEnd {
				if len(sizes) < 2 {
					return errKV3OutOfData
				}
				compressedSize := int(binary.LittleEndian.Uint16(sizes))
				sizes = sizes[2:]
				input := reader.take(compressedSize)
				if reader.err != nil {
					return reader.err
				}
				before := len(blobsData)
				var err error
				blobsData, err = decodeLZ4Block(blobsData, input)
				if err != nil {
					return err
				}
				if len(blobsData) == before {
					return fmt.Errorf("kv3: empty LZ4 blob frame")
				}
			}
		}
	case compressionZstd:
		if header.version >= 5 {
			compressedSize := header.sizeCompressedTotal - header.sizeCompressedBuffer1 - header.sizeCompressedBuffer2
			input := reader.take(compressedSize)
			if reader.err != nil {
				return reader.err
			}
			var err error
			blobsData, err = decodeZstd(input, total)
			if err != nil {
				return err
			}
		} else {
			if header.sizeUncompressedBuffer1+total > len(buffer1Raw) {
				return errKV3OutOfData
			}
			blobsData = buffer1Raw[header.sizeUncompressedBuffer1 : header.sizeUncompressedBuffer1+total]
		}
	default:
		return fmt.Errorf("kv3: unknown compression method %d", header.compressionMethod)
	}

	if len(blobsData) != total {
		return fmt.Errorf("kv3: decoded %d binary blob bytes, expected %d", len(blobsData), total)
	}

	context.binaryBlobs = make([][]byte, header.countBlocks)
	offset := 0
	for i, length := range blobLengths {
		context.binaryBlobs[i] = blobsData[offset : offset+length : offset+length]
		offset += length
	}

	if reader.uint32() != kv3BlockTrailer {
		return fmt.Errorf("kv3: invalid binary blobs trailer")
	}

	return nil
}

func (context *kv3Context) readType() (kv3NodeType, error) {
	if len(context.types) == 0 {
		return 0, errKV3OutOfData
	}
	dataByte := context.types[0]
	context.types = context.types[1:]

	if context.version >= 3 {
		if dataByte&kv3FlagByteMask != 0 {
			if len(context.types) == 0 {
				return 0, errKV3OutOfData
			}
			context.types = context.types[1:] // flag, unused
		}
		if dataByte&kv3ExtraByteMask != 0 {
			if len(context.types) == 0 {
				return 0, errKV3OutOfData
			}
			context.types = context.types[1:]
		}
		dataByte &= 0x3f
	} else {
		if dataByte&kv3FlagByteMask != 0 {
			if len(context.types) == 0 {
				return 0, errKV3OutOfData
			}
			context.types = context.types[1:]
		}
		dataByte &= 0x7f
	}

	return kv3NodeType(dataByte), nil
}

func (context *kv3Context) string(id int) string {
	if id < 0 || id >= len(context.strings) {
		return ""
	}

	return context.strings[id]
}

func readLane(lane *[]byte, size int) ([]byte, error) {
	if len(*lane) < size {
		return nil, errKV3OutOfData
	}
	value := (*lane)[:size]
	*lane = (*lane)[size:]

	return value, nil
}

func readInt32Lane(lane *[]byte) (int, error) {
	bytes, err := readLane(lane, 4)
	if err != nil {
		return 0, err
	}

	return int(int32(binary.LittleEndian.Uint32(bytes))), nil
}

func readPackedArray(elementType kv3NodeType, lane *kv3Lanes, count int) (any, bool, error) {
	switch elementType {
	case kv3Boolean:
		bytes, err := readLane(&lane.bytes1, count)
		if err != nil {
			return nil, true, err
		}
		values := make([]bool, count)
		for i, value := range bytes {
			values[i] = value != 0
		}

		return values, true, nil
	case kv3Int8:
		bytes, err := readLane(&lane.bytes1, count)
		if err != nil {
			return nil, true, err
		}
		values := make([]int32, count)
		for i, value := range bytes {
			values[i] = int32(int8(value))
		}

		return values, true, nil
	case kv3UInt8:
		bytes, err := readLane(&lane.bytes1, count)
		if err != nil {
			return nil, true, err
		}
		values := make([]uint32, count)
		for i, value := range bytes {
			values[i] = uint32(value)
		}

		return values, true, nil
	case kv3Int16:
		bytes, err := readLane(&lane.bytes2, count*2)
		if err != nil {
			return nil, true, err
		}
		values := make([]int16, count)
		for i := range values {
			values[i] = int16(binary.LittleEndian.Uint16(bytes[i*2:]))
		}

		return values, true, nil
	case kv3UInt16:
		bytes, err := readLane(&lane.bytes2, count*2)
		if err != nil {
			return nil, true, err
		}
		values := make([]uint16, count)
		for i := range values {
			values[i] = binary.LittleEndian.Uint16(bytes[i*2:])
		}

		return values, true, nil
	case kv3Int32:
		bytes, err := readLane(&lane.bytes4, count*4)
		if err != nil {
			return nil, true, err
		}
		values := make([]int32, count)
		for i := range values {
			values[i] = int32(binary.LittleEndian.Uint32(bytes[i*4:]))
		}

		return values, true, nil
	case kv3UInt32:
		bytes, err := readLane(&lane.bytes4, count*4)
		if err != nil {
			return nil, true, err
		}
		values := make([]uint32, count)
		for i := range values {
			values[i] = binary.LittleEndian.Uint32(bytes[i*4:])
		}

		return values, true, nil
	case kv3Float:
		bytes, err := readLane(&lane.bytes4, count*4)
		if err != nil {
			return nil, true, err
		}
		values := make([]float32, count)
		for i := range values {
			values[i] = math.Float32frombits(binary.LittleEndian.Uint32(bytes[i*4:]))
		}

		return values, true, nil
	case kv3Int64:
		bytes, err := readLane(&lane.bytes8, count*8)
		if err != nil {
			return nil, true, err
		}
		values := make([]int64, count)
		for i := range values {
			values[i] = int64(binary.LittleEndian.Uint64(bytes[i*8:]))
		}

		return values, true, nil
	case kv3UInt64:
		bytes, err := readLane(&lane.bytes8, count*8)
		if err != nil {
			return nil, true, err
		}
		values := make([]uint64, count)
		for i := range values {
			values[i] = binary.LittleEndian.Uint64(bytes[i*8:])
		}

		return values, true, nil
	case kv3Double:
		bytes, err := readLane(&lane.bytes8, count*8)
		if err != nil {
			return nil, true, err
		}
		values := make([]float64, count)
		for i := range values {
			values[i] = math.Float64frombits(binary.LittleEndian.Uint64(bytes[i*8:]))
		}

		return values, true, nil
	default:
		return nil, false, nil
	}
}

// readValue reads a value. Primitives come from the given lane (the auxiliary buffer for ARRAY_TYPE_AUXILIARY_BUFFER
// elements), strings, blobs and container lengths always come from the main buffer.
func (context *kv3Context) readValue(nodeType kv3NodeType, lane *kv3Lanes) (any, error) {
	buffer := context.buffer

	switch nodeType {
	case kv3Null:
		return nil, nil
	case kv3BooleanTrue:
		return true, nil
	case kv3BooleanFalse:
		return false, nil
	case kv3Int64Zero:
		return int64(0), nil
	case kv3Int64One:
		return int64(1), nil
	case kv3DoubleZero:
		return float64(0), nil
	case kv3DoubleOne:
		return float64(1), nil
	case kv3UInt32Zero:
		return uint64(0), nil
	case kv3UInt32One:
		return uint64(1), nil
	case kv3Boolean:
		bytes, err := readLane(&lane.bytes1, 1)
		if err != nil {
			return nil, err
		}

		return bytes[0] != 0, nil
	case kv3Int8:
		bytes, err := readLane(&lane.bytes1, 1)
		if err != nil {
			return nil, err
		}

		return int64(int8(bytes[0])), nil
	case kv3UInt8:
		bytes, err := readLane(&lane.bytes1, 1)
		if err != nil {
			return nil, err
		}

		return uint64(bytes[0]), nil
	case kv3Int16:
		bytes, err := readLane(&lane.bytes2, 2)
		if err != nil {
			return nil, err
		}

		return int64(int16(binary.LittleEndian.Uint16(bytes))), nil
	case kv3UInt16:
		bytes, err := readLane(&lane.bytes2, 2)
		if err != nil {
			return nil, err
		}

		return uint64(binary.LittleEndian.Uint16(bytes)), nil
	case kv3Int32:
		bytes, err := readLane(&lane.bytes4, 4)
		if err != nil {
			return nil, err
		}

		return int64(int32(binary.LittleEndian.Uint32(bytes))), nil
	case kv3UInt32:
		bytes, err := readLane(&lane.bytes4, 4)
		if err != nil {
			return nil, err
		}

		return uint64(binary.LittleEndian.Uint32(bytes)), nil
	case kv3Float:
		bytes, err := readLane(&lane.bytes4, 4)
		if err != nil {
			return nil, err
		}

		return float64(math.Float32frombits(binary.LittleEndian.Uint32(bytes))), nil
	case kv3Int64:
		bytes, err := readLane(&lane.bytes8, 8)
		if err != nil {
			return nil, err
		}

		return int64(binary.LittleEndian.Uint64(bytes)), nil
	case kv3UInt64:
		bytes, err := readLane(&lane.bytes8, 8)
		if err != nil {
			return nil, err
		}

		return binary.LittleEndian.Uint64(bytes), nil
	case kv3Double:
		bytes, err := readLane(&lane.bytes8, 8)
		if err != nil {
			return nil, err
		}

		return math.Float64frombits(binary.LittleEndian.Uint64(bytes)), nil
	case kv3String:
		id, err := readInt32Lane(&buffer.bytes4)
		if err != nil {
			return nil, err
		}

		return context.string(id), nil
	case kv3BinaryBlob:
		if context.version < 2 {
			length, err := readInt32Lane(&buffer.bytes4)
			if err != nil {
				return nil, err
			}
			if length <= 0 {
				return []byte{}, nil
			}

			return readLane(&buffer.bytes1, length)
		}
		if context.nextBinaryBlob >= len(context.binaryBlobs) {
			return nil, errKV3OutOfData
		}
		blob := context.binaryBlobs[context.nextBinaryBlob]
		context.nextBinaryBlob++

		return blob, nil
	case kv3Array:
		length, err := readInt32Lane(&buffer.bytes4)
		if err != nil {
			return nil, err
		}
		if length < 0 {
			return nil, fmt.Errorf("kv3: invalid array length %d", length)
		}
		array := make([]any, 0, length)
		for range length {
			elementType, err := context.readType()
			if err != nil {
				return nil, err
			}
			value, err := context.readValue(elementType, context.buffer)
			if err != nil {
				return nil, err
			}
			array = append(array, value)
		}

		return array, nil
	case kv3ArrayTyped, kv3ArrayTypeByteLength, kv3ArrayTypeAuxiliaryBuff:
		var length int
		if nodeType == kv3ArrayTyped {
			value, err := readInt32Lane(&buffer.bytes4)
			if err != nil {
				return nil, err
			}
			length = value
		} else {
			bytes, err := readLane(&buffer.bytes1, 1)
			if err != nil {
				return nil, err
			}
			length = int(bytes[0])
		}
		if length < 0 {
			return nil, fmt.Errorf("kv3: invalid typed array length %d", length)
		}

		elementType, err := context.readType()
		if err != nil {
			return nil, err
		}
		elementLane := buffer
		if nodeType == kv3ArrayTypeAuxiliaryBuff {
			elementLane = context.auxiliaryBuffer
		}

		packed, ok, err := readPackedArray(elementType, elementLane, length)
		if err != nil {
			return nil, err
		}
		if ok {
			return packed, nil
		}

		array := make([]any, 0, length)
		for range length {
			value, err := context.readValue(elementType, elementLane)
			if err != nil {
				return nil, err
			}
			array = append(array, value)
		}

		return array, nil
	case kv3Object:
		var length int
		var err error
		if context.version >= 5 {
			length, err = readInt32Lane(&context.objectLengths)
		} else {
			length, err = readInt32Lane(&buffer.bytes4)
		}
		if err != nil {
			return nil, err
		}
		if length < 0 {
			return nil, fmt.Errorf("kv3: invalid object length %d", length)
		}

		object := make(map[string]any, length)
		for range length {
			memberType, err := context.readType()
			if err != nil {
				return nil, err
			}
			nameID, err := readInt32Lane(&context.buffer.bytes4)
			if err != nil {
				return nil, err
			}
			value, err := context.readValue(memberType, context.buffer)
			if err != nil {
				return nil, err
			}
			object[context.string(nameID)] = value
		}

		return object, nil
	default:
		return nil, fmt.Errorf("kv3: unknown node type %d", nodeType)
	}
}
