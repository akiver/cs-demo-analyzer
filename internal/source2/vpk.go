// Package source2 reads the Source 2 file formats needed to extract CS2 map collision geometry:
// VPK archives, compiled resources (_c files) and binary KeyValues3.
package source2

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	vpkSignature        = 0x55aa1234
	vpkDirArchiveIndex  = 0x7fff
	vpkEntryTerminator  = 0xffff
	vpkHeaderSizeV1     = 12
	vpkHeaderSizeV2     = 28
	vpkEntryFixedLength = 18
)

// ErrVPKFileNotFound is returned by VPK.ReadFile when the archive has no entry for the requested path.
var ErrVPKFileNotFound = errors.New("file not found in VPK")

type vpkEntry struct {
	preload      []byte
	archiveIndex uint16
	offset       uint32
	length       uint32
}

// VPK is a Valve pak archive.
// Only single-file archives (data stored in the _dir file itself, such as maps/<map>.vpk) and multi-chunk archives
// whose chunks are next to the directory file (pak01_dir.vpk + pak01_000.vpk...) are supported.
type VPK struct {
	path       string
	dataOffset int64
	entries    map[string]vpkEntry
}

// OpenVPK parses the directory tree of a VPK file.
func OpenVPK(path string) (*VPK, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var header [vpkHeaderSizeV2]byte
	if _, err := io.ReadFull(reader, header[:vpkHeaderSizeV1]); err != nil {
		return nil, fmt.Errorf("read VPK header: %w", err)
	}

	if binary.LittleEndian.Uint32(header[0:]) != vpkSignature {
		return nil, fmt.Errorf("%s is not a VPK file", path)
	}

	version := binary.LittleEndian.Uint32(header[4:])
	treeSize := binary.LittleEndian.Uint32(header[8:])
	headerSize := vpkHeaderSizeV1
	switch version {
	case 1:
	case 2:
		if _, err := io.ReadFull(reader, header[vpkHeaderSizeV1:]); err != nil {
			return nil, fmt.Errorf("read VPK v2 header: %w", err)
		}
		headerSize = vpkHeaderSizeV2
	default:
		return nil, fmt.Errorf("unsupported VPK version %d", version)
	}

	tree := make([]byte, treeSize)
	if _, err := io.ReadFull(reader, tree); err != nil {
		return nil, fmt.Errorf("read VPK tree: %w", err)
	}

	vpk := &VPK{
		path:       path,
		dataOffset: int64(headerSize) + int64(treeSize),
		entries:    make(map[string]vpkEntry),
	}

	if err := vpk.parseTree(tree); err != nil {
		return nil, err
	}

	return vpk, nil
}

func (vpk *VPK) parseTree(tree []byte) error {
	position := 0
	readString := func() (string, error) {
		end := position
		for end < len(tree) && tree[end] != 0 {
			end++
		}
		if end >= len(tree) {
			return "", fmt.Errorf("unterminated string in VPK tree")
		}
		value := string(tree[position:end])
		position = end + 1

		return value, nil
	}

	for {
		extension, err := readString()
		if err != nil {
			return err
		}
		if extension == "" {
			return nil
		}

		for {
			directory, err := readString()
			if err != nil {
				return err
			}
			if directory == "" {
				break
			}

			for {
				name, err := readString()
				if err != nil {
					return err
				}
				if name == "" {
					break
				}

				if position+vpkEntryFixedLength > len(tree) {
					return fmt.Errorf("truncated VPK entry %s", name)
				}

				preloadSize := int(binary.LittleEndian.Uint16(tree[position+4:]))
				entry := vpkEntry{
					archiveIndex: binary.LittleEndian.Uint16(tree[position+6:]),
					offset:       binary.LittleEndian.Uint32(tree[position+8:]),
					length:       binary.LittleEndian.Uint32(tree[position+12:]),
				}
				if binary.LittleEndian.Uint16(tree[position+16:]) != vpkEntryTerminator {
					return fmt.Errorf("invalid VPK entry terminator for %s", name)
				}
				position += vpkEntryFixedLength

				if position+preloadSize > len(tree) {
					return fmt.Errorf("truncated VPK preload data for %s", name)
				}
				entry.preload = tree[position : position+preloadSize]
				position += preloadSize

				vpk.entries[vpkEntryPath(directory, name, extension)] = entry
			}
		}
	}
}

func vpkEntryPath(directory string, name string, extension string) string {
	var builder strings.Builder
	if directory != " " {
		builder.WriteString(directory)
		builder.WriteByte('/')
	}
	if name != " " {
		builder.WriteString(name)
	}
	if extension != " " {
		builder.WriteByte('.')
		builder.WriteString(extension)
	}

	return strings.ToLower(builder.String())
}

// Files returns the paths of all files in the archive.
func (vpk *VPK) Files() []string {
	files := make([]string, 0, len(vpk.entries))
	for path := range vpk.entries {
		files = append(files, path)
	}

	return files
}

// ReadFile returns the content of a file stored in the archive.
func (vpk *VPK) ReadFile(path string) ([]byte, error) {
	entry, ok := vpk.entries[strings.ToLower(strings.ReplaceAll(path, "\\", "/"))]
	if !ok {
		return nil, fmt.Errorf("%w: %s in %s", ErrVPKFileNotFound, path, vpk.path)
	}

	data := make([]byte, len(entry.preload)+int(entry.length))
	copy(data, entry.preload)
	if entry.length == 0 {
		return data, nil
	}

	archivePath := vpk.path
	offset := int64(entry.offset)
	if entry.archiveIndex == vpkDirArchiveIndex {
		offset += vpk.dataOffset
	} else {
		dir, base := filepath.Split(vpk.path)
		base = strings.TrimSuffix(base, "_dir.vpk")
		archivePath = filepath.Join(dir, fmt.Sprintf("%s_%03d.vpk", base, entry.archiveIndex))
	}

	file, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if _, err := file.ReadAt(data[len(entry.preload):], offset); err != nil {
		return nil, fmt.Errorf("read %s from %s: %w", path, archivePath, err)
	}

	return data, nil
}
