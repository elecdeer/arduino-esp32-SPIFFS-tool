package cmd

import (
	"fmt"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/cmd/options"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/operation/common"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/operation/pack"
	"github.com/spf13/cobra"
	"os"
)

// packCmd represents the pack command
var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "Create a SPIFFS image from files",
	Args:  cobra.ExactArgs(2),
	RunE:  RunPack,
}

func RunPack(cmd *cobra.Command, args []string) error {
	common.ApplyLogStyle()
	distPath := args[0]
	sourceDir := args[1]

	if !common.IsDirExists(sourceDir) {
		return fmt.Errorf("Source directory not found: %s\n", sourceDir)
	}

	partition, err := options.ResolvePartitionOptions()
	if err != nil {
		return fmt.Errorf("Error resolving partition scheme: %s\n", err)
	}

	tempDir, err := common.MakeTempDir()
	if err != nil {
		return fmt.Errorf("failed to make temp dir\n%s\n", err)
	}
	defer os.RemoveAll(tempDir)
	cmd.Printf("use temp dir: %s\n", tempDir)

	cmd.Printf("copy files to temp dir\n")
	err = pack.MakeTargetDir(sourceDir, tempDir, packOptions.ignoreDotfile, packOptions.ignorePattern)
	if err != nil {
		return fmt.Errorf("failed to make target dir\n%s\n", err)
	}

	pack.PrintDirFiles(tempDir)

	err = pack.MakeSpiffsImage(packOptions.packToolPath, distPath, sourceDir, pack.MakeSpiffsImageOptions{
		PageSize:  packOptions.pageSize,
		BlockSize: packOptions.blockSize,
		Partition: &partition,
	})
	if err != nil {
		return fmt.Errorf("failed to make spiffs image\n%s\n", err)
	}
	cmd.Printf("SPIFFS image created: %s\n", distPath)

	return nil
}

type PackParam struct {
	partition     *options.PartitionSchema
	mkspiffsPath  string
	distPath      string
	sourceDir     string
	ignoreDotfile bool
	ignorePattern string
	pageSize      int
	blockSize     int
}

//type FileSystem string
//
//const (
//	LittleFS = FileSystem("littlefs")
//	Spiffs   = FileSystem("spiffs")
//)

type ParsedPackOptions struct {
	packToolPath        string
	partitionSchemePath string
	//fileSystem          FileSystem
	ignoreDotfile bool
	ignorePattern string
	blockSize     int
	pageSize      int
}

var (
	packOptions ParsedPackOptions
)

func init() {
	rootCmd.AddCommand(packCmd)

	packCmd.Flags().StringVarP(&packOptions.packToolPath, "tool", "t", "", "Pack tool path")
	packCmd.MarkFlagRequired("tool")

	options.AddPartitionOptions(packCmd)
	AddPackCommandOptions(packCmd)
}

func AddPackCommandOptions(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&packOptions.ignoreDotfile, "dotfile", true, "Ignore dotfile")
	cmd.Flags().StringVarP(&packOptions.ignorePattern, "ignore", "i", "", "Ignore file regexp pattern")

	cmd.Flags().IntVar(&packOptions.blockSize, "block", 4096, "Block size")
	cmd.Flags().IntVar(&packOptions.pageSize, "page", 256, "Page size")
}
