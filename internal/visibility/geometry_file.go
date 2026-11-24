package visibility

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

// Map geometry file layout (little-endian, the whole file is gzip compressed):
//
//	magic        [8]byte "CSDAGEO\x00"
//	version      uint32
//	vertexCount  uint32
//	indexCount   uint32
//	nodeCount    uint32
//	vertices     vertexCount * 3 * float32
//	indices      indexCount * uint32 (3 per triangle, in BVH leaf order)
//	nodes        nodeCount * (6 * float32 + 2 * uint32)
const (
	GeometryFileExtension = ".tri"
	// Bumped when the layout changes, readers reject other versions.
	GeometryFileVersion = 1
	vertexSize          = 12
	indexSize           = 4
	nodeSize            = 32
	// Upper bound of each section's element count, far above the largest official map (about 500k indices), to limit
	// the memory used when reading a corrupted file.
	maxGeometryElements = 1 << 24
	// Sections are decoded by chunks of this many elements.
	geometryChunkElements = 4096
	// Maximum number of elements allocated from the count declared in the header before any of them has been read.
	// It covers the official maps, bigger sections grow with the data actually read.
	geometryPreallocatedElements = 1 << 20
)

var geometryFileMagic = [8]byte{'C', 'S', 'D', 'A', 'G', 'E', 'O', 0}

// ErrInvalidGeometryFile is returned when a map geometry file is not valid.
var ErrInvalidGeometryFile = errors.New("invalid map geometry file")

// WriteGeometry writes the BVH in the map geometry file format.
func WriteGeometry(writer io.Writer, bvh *BVH) error {
	if len(bvh.Vertices) > maxGeometryElements || len(bvh.Indices) > maxGeometryElements || len(bvh.Nodes) > maxGeometryElements {
		return fmt.Errorf("geometry too large to be written (%d vertices, %d indices, %d nodes)", len(bvh.Vertices), len(bvh.Indices), len(bvh.Nodes))
	}

	gzipWriter, err := gzip.NewWriterLevel(writer, gzip.BestCompression)
	if err != nil {
		return err
	}

	bufferedWriter := bufio.NewWriterSize(gzipWriter, 1<<20)
	var scratch [nodeSize]byte
	write := func(data []byte) {
		if err == nil {
			_, err = bufferedWriter.Write(data)
		}
	}
	writeUint32 := func(value uint32) {
		binary.LittleEndian.PutUint32(scratch[:4], value)
		write(scratch[:4])
	}
	writeVec3 := func(vector Vec3, buffer []byte) {
		for axis := range 3 {
			binary.LittleEndian.PutUint32(buffer[axis*4:], math.Float32bits(vector[axis]))
		}
	}

	write(geometryFileMagic[:])
	writeUint32(GeometryFileVersion)
	writeUint32(uint32(len(bvh.Vertices)))
	writeUint32(uint32(len(bvh.Indices)))
	writeUint32(uint32(len(bvh.Nodes)))
	for _, vertex := range bvh.Vertices {
		writeVec3(vertex, scratch[:12])
		write(scratch[:12])
	}
	for _, index := range bvh.Indices {
		writeUint32(index)
	}
	for _, node := range bvh.Nodes {
		writeVec3(node.Min, scratch[0:12])
		writeVec3(node.Max, scratch[12:24])
		binary.LittleEndian.PutUint32(scratch[24:], node.LeftOrFirst)
		binary.LittleEndian.PutUint32(scratch[28:], node.Count)
		write(scratch[:nodeSize])
	}

	if err != nil {
		return err
	}
	if err := bufferedWriter.Flush(); err != nil {
		return err
	}

	return gzipWriter.Close()
}

// ReadGeometry reads a BVH written by WriteGeometry.
func ReadGeometry(reader io.Reader) (*BVH, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidGeometryFile, err)
	}
	defer gzipReader.Close()

	bufferedReader := bufio.NewReaderSize(gzipReader, 1<<20)
	var header [24]byte
	if _, err := io.ReadFull(bufferedReader, header[:]); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidGeometryFile, err)
	}
	if [8]byte(header[:8]) != geometryFileMagic {
		return nil, fmt.Errorf("%w: bad magic", ErrInvalidGeometryFile)
	}
	if version := binary.LittleEndian.Uint32(header[8:]); version != GeometryFileVersion {
		return nil, fmt.Errorf("%w: unsupported version %d", ErrInvalidGeometryFile, version)
	}

	vertexCount := binary.LittleEndian.Uint32(header[12:])
	indexCount := binary.LittleEndian.Uint32(header[16:])
	nodeCount := binary.LittleEndian.Uint32(header[20:])
	if vertexCount > maxGeometryElements || indexCount > maxGeometryElements || nodeCount > maxGeometryElements || indexCount%3 != 0 {
		return nil, fmt.Errorf("%w: invalid counts", ErrInvalidGeometryFile)
	}

	readVec3 := func(data []byte) Vec3 {
		return Vec3{
			math.Float32frombits(binary.LittleEndian.Uint32(data)),
			math.Float32frombits(binary.LittleEndian.Uint32(data[4:])),
			math.Float32frombits(binary.LittleEndian.Uint32(data[8:])),
		}
	}

	vertices, err := readGeometrySection(bufferedReader, vertexCount, vertexSize, func(data []byte, _ int) (Vec3, error) {
		return readVec3(data), nil
	})
	if err != nil {
		return nil, err
	}

	indices, err := readGeometrySection(bufferedReader, indexCount, indexSize, func(data []byte, _ int) (uint32, error) {
		index := binary.LittleEndian.Uint32(data)
		if index >= vertexCount {
			return 0, fmt.Errorf("%w: vertex index out of range", ErrInvalidGeometryFile)
		}

		return index, nil
	})
	if err != nil {
		return nil, err
	}

	triangleCount := indexCount / 3
	nodes, err := readGeometrySection(bufferedReader, nodeCount, nodeSize, func(data []byte, i int) (Node, error) {
		node := Node{
			Min:         readVec3(data),
			Max:         readVec3(data[12:]),
			LeftOrFirst: binary.LittleEndian.Uint32(data[24:]),
			Count:       binary.LittleEndian.Uint32(data[28:]),
		}
		if node.Count > 0 {
			if uint64(node.LeftOrFirst)+uint64(node.Count) > uint64(triangleCount) {
				return Node{}, fmt.Errorf("%w: leaf out of range", ErrInvalidGeometryFile)
			}
		} else if uint64(node.LeftOrFirst)+1 >= uint64(nodeCount) || node.LeftOrFirst <= uint32(i) {
			// Children are always stored after their parent, which also guarantees the traversal terminates.
			return Node{}, fmt.Errorf("%w: node child out of range", ErrInvalidGeometryFile)
		}

		return node, nil
	})
	if err != nil {
		return nil, err
	}

	return &BVH{Vertices: vertices, Indices: indices, Nodes: nodes}, nil
}

// readGeometrySection decodes the count elements of elementSize bytes of a section, decode receives the bytes of an
// element and its index in the section.
// The section is read by chunks and the returned slice grows with the data actually read, so that a corrupted header
// declaring huge counts can't trigger a huge allocation.
func readGeometrySection[T any](reader io.Reader, count uint32, elementSize int, decode func(data []byte, index int) (T, error)) ([]T, error) {
	total := int(count)
	elements := make([]T, 0, min(total, geometryPreallocatedElements))
	chunk := make([]byte, min(total, geometryChunkElements)*elementSize)
	for len(elements) < total {
		data := chunk[:min(total-len(elements), geometryChunkElements)*elementSize]
		if _, err := io.ReadFull(reader, data); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidGeometryFile, err)
		}

		for offset := 0; offset < len(data); offset += elementSize {
			element, err := decode(data[offset:offset+elementSize], len(elements))
			if err != nil {
				return nil, err
			}
			elements = append(elements, element)
		}
	}

	return elements, nil
}

// LoadGeometryFile reads a map geometry file from disk.
// It also returns the SHA-256 of the file as a hex string, which identifies the version of the map geometry. It's the
// value listed for the file in the manifest generated by cmd/map-geometry.
func LoadGeometryFile(path string) (*BVH, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	bvh, err := ReadGeometry(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}
	checksum := sha256.Sum256(data)

	return bvh, hex.EncodeToString(checksum[:]), nil
}
