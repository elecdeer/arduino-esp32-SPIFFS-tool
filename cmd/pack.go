package cmd

import (
	"fmt"
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
	Run:   RunPack,
}

func RunPack(cmd *cobra.Command, args []string) {
	common.ApplyLogStyle()
	param, err := ConstructPackParams(args)
	if err != nil {
		cmd.PrintErrf("failed to construct pack params\n%s\n", err)
		return
	}

	tempDir, err := common.MakeTempDir()
	if err != nil {
		cmd.PrintErrf("failed to make temp dir\n%s\n", err)
		return
	}
	defer os.RemoveAll(tempDir)
	cmd.Printf("use temp dir: %s\n", tempDir)

	cmd.Printf("copy files to temp dir\n")
	err = pack.MakeTargetDir(param.sourceDir, tempDir, param.ignoreDotfile, param.ignorePattern)
	if err != nil {
		cmd.PrintErrf("failed to make target dir\n%s\n", err)
		return
	}

	pack.PrintDirFiles(tempDir)

	err = pack.MakeSpiffsImage(param.mkspiffsPath, param.distPath, param.sourceDir, pack.MakeSpiffsImageOptions{
		PageSize:  param.pageSize,
		BlockSize: param.blockSize,
		Partition: param.partition,
	})
	if err != nil {
		cmd.PrintErrf("failed to make spiffs image\n%s\n", err)
		return
	}
	cmd.Printf("SPIFFS image created: %s\n", param.distPath)
}

type PackParam struct {
	partition     *pack.PartitionSchema
	mkspiffsPath  string
	distPath      string
	sourceDir     string
	ignoreDotfile bool
	ignorePattern string
	pageSize      int
	blockSize     int
}

func ConstructPackParams(args []string) (PackParam, error) {
	distPath := args[0]
	sourceDir := args[1]

	if !common.IsFileExists(partitionSchemePath) {
		return PackParam{}, fmt.Errorf("Partition scheme file not found: %s\n", partitionSchemePath)
	}

	if !common.IsDirExists(sourceDir) {
		return PackParam{}, fmt.Errorf("Source directory not found: %s\n", sourceDir)
	}

	partition, err := pack.ReadPartitionSchemeFile(partitionSchemePath)
	if err != nil {
		return PackParam{}, fmt.Errorf("Error reading partition scheme file: %s\n", err)
	}

	return PackParam{
		distPath:      distPath,
		sourceDir:     sourceDir,
		partition:     &partition,
		mkspiffsPath:  mkspiffsPath,
		ignoreDotfile: true,
		ignorePattern: "",
		pageSize:      256,
		blockSize:     4096,
	}, nil
}

var (
	partitionSchemePath string
	mkspiffsPath        string
)

func init() {
	rootCmd.AddCommand(packCmd)

	packCmd.Flags().StringVarP(&partitionSchemePath, "partition-scheme", "p", "", "Partition scheme csv file path")
	packCmd.MarkFlagRequired("partition-scheme")

	packCmd.Flags().StringVarP(&mkspiffsPath, "mkspiffs", "m", "", "mkspiffs path")
}
