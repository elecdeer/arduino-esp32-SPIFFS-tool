package cmd

import (
	"fmt"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/cmd/options"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/operation/common"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/operation/pack"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/operation/write"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Create as SPIFFS image and write to ESP32 flash memory",
	Args:  cobra.ExactArgs(1),
	RunE:  RunUpload,
}

func RunUpload(cmd *cobra.Command, args []string) error {
	common.ApplyLogStyle()
	sourceDir := args[0]

	if !common.IsFileExists(packOptions.packToolPath) {
		return fmt.Errorf("Pack tool not found: %s\n", packOptions.packToolPath)
	}
	if !common.IsFileExists(writeOptions.writeToolPath) {
		return fmt.Errorf("write tool not found: %s", writeOptions.writeToolPath)
	}

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
	sourceTempDir := filepath.Join(tempDir, "source")
	err = pack.MakeTargetDir(sourceDir, sourceTempDir, packOptions.ignoreDotfile, packOptions.ignorePattern)
	if err != nil {
		return fmt.Errorf("failed to make target dir\n%s\n", err)
	}

	pack.PrintDirFiles(sourceTempDir)
	distImagePath := filepath.Join(tempDir, "image.bin")
	err = pack.MakeSpiffsImage(packOptions.packToolPath, distImagePath, sourceDir, pack.MakeSpiffsImageOptions{
		PageSize:  packOptions.pageSize,
		BlockSize: packOptions.blockSize,
		Partition: &partition,
	})
	if err != nil {
		return fmt.Errorf("failed to make spiffs image\n%s\n", err)
	}
	cmd.Printf("SPIFFS image created: %s\n", distImagePath)

	cmd.Printf("Writing image file: %s\n", distImagePath)
	err = write.WriteImageWithSerial(writeOptions.writeToolPath, distImagePath, write.WriteImageOptions{
		Partition:  partition,
		Chip:       writeOptions.chip,
		SerialBaud: writeOptions.serialBaud,
		SerialPort: writeOptions.serialPort,
		FlashMode:  writeOptions.flashMode,
		FlashFreq:  writeOptions.flashFreq,
		FlashSize:  writeOptions.flashSize,
		NoVerify:   writeOptions.noVerify,
	})
	if err != nil {
		return fmt.Errorf("Error writing image file: %s\n", err)
	}
	cmd.PrintErrf("Write image successfully\n")

	return nil
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVar(&packOptions.packToolPath, "pack-tool", "", "Pack tool path")
	uploadCmd.MarkFlagRequired("pack-tool")

	uploadCmd.Flags().StringVar(&writeOptions.writeToolPath, "write-tool", "", "Write tool path")
	uploadCmd.MarkFlagRequired("write-tool")

	options.AddPartitionOptions(uploadCmd)
	AddWriteCommandOptions(uploadCmd)
	AddPackCommandOptions(uploadCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
