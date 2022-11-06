package pack

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PartitionSchema struct {
	Offset uint64
	Size   uint64
}

func ReadPartitionSchemeFile(path string) (PartitionSchema, error) {
	file, err := os.Open(path)
	if err != nil {
		return PartitionSchema{}, fmt.Errorf("failed to open partition scheme file: %w", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return PartitionSchema{}, fmt.Errorf("failed to read partition scheme file: %w", err)
	}

	for _, record := range records {
		if record[0] == "spiffs" {
			offset, err := strconv.ParseUint(strings.TrimSpace(record[3]), 0, 32)
			if err != nil {
				return PartitionSchema{}, fmt.Errorf("failed to parse Offset: %w", err)
			}

			size, err := strconv.ParseUint(strings.TrimSpace(record[4]), 0, 32)
			if err != nil {
				return PartitionSchema{}, fmt.Errorf("failed to parse Size: %w", err)
			}

			return PartitionSchema{
				Offset: offset,
				Size:   size,
			}, nil
		}
	}

	return PartitionSchema{}, fmt.Errorf("no SPIFFS partition found in partition scheme file")
}
