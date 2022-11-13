package options

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type PartitionSchema struct {
	Offset uint64
	Size   uint64
}

var (
	partitionSchemaPath string
	offset              uint64
	size                uint64
)

func AddPartitionOptions(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&partitionSchemaPath, "partition-schema", "p", "", "partition schema file path")

	// MaxUint64を無指定値として使っている
	cmd.Flags().Uint64Var(&offset, "offset", math.MaxUint64, "partition offset")
	cmd.Flags().Uint64Var(&size, "size", math.MaxUint64, "partition size")
}

func ResolvePartitionOptions() (PartitionSchema, error) {
	if partitionSchemaPath != "" {
		log.Printf("read partition schema file: %s", partitionSchemaPath)
		schema, err := readPartitionSchemeFile(partitionSchemaPath)
		if err != nil {
			log.Printf("failed to read partition schema file: %s", err)
			return PartitionSchema{}, err
		}
		log.Printf("partition schema: %v", schema)

		//個別に指定されたオプションを優先する
		if offset == math.MaxUint64 {
			log.Printf("use offset from partition schema file")
			offset = schema.Offset
		}
		if size == math.MaxUint64 {
			log.Printf("use size from partition schema file")
			size = schema.Size
		}
	}

	if offset == math.MaxUint64 || size == math.MaxUint64 {
		return PartitionSchema{}, fmt.Errorf("offset and size must be specified")
	}

	return PartitionSchema{
		Offset: offset,
		Size:   size,
	}, nil
}

func readPartitionSchemeFile(path string) (PartitionSchema, error) {
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
