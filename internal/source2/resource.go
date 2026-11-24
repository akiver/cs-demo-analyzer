package source2

import (
	"encoding/binary"
	"fmt"
)

const resourceHeaderSize = 16

// ResourceBlock returns the raw content of the first block of the given type (e.g. "DATA", "PHYS") of a compiled
// Source 2 resource (_c file).
func ResourceBlock(resource []byte, blockType string) ([]byte, error) {
	if len(resource) < resourceHeaderSize {
		return nil, fmt.Errorf("resource too small")
	}

	blockOffset := int(binary.LittleEndian.Uint32(resource[8:]))
	blockCount := int(binary.LittleEndian.Uint32(resource[12:]))
	position := 8 + blockOffset
	for range blockCount {
		if position+12 > len(resource) {
			return nil, fmt.Errorf("truncated resource block table")
		}

		currentType := string(resource[position : position+4])
		offset := int(binary.LittleEndian.Uint32(resource[position+4:]))
		size := int(binary.LittleEndian.Uint32(resource[position+8:]))
		start := position + 4 + offset
		if currentType == blockType {
			if start < 0 || start+size > len(resource) {
				return nil, fmt.Errorf("resource block %s is out of bounds", blockType)
			}

			return resource[start : start+size], nil
		}

		position += 12
	}

	return nil, fmt.Errorf("resource block %s not found", blockType)
}
