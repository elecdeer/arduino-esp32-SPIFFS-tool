package cmd

import (
	"fmt"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/cmd/options"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/operation/common"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/operation/write"

	"github.com/spf13/cobra"
)

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write a image file to ESP32 flash memory",
	Args:  cobra.ExactArgs(1),
	RunE:  RunWrite,
}

func RunWrite(cmd *cobra.Command, args []string) error {
	common.ApplyLogStyle()

	writeImagePath := args[0]

	if !common.IsFileExists(writeImagePath) {
		return fmt.Errorf("image file not found: %s", writeImagePath)
	}

	if !common.IsFileExists(writeOptions.writeToolPath) {
		return fmt.Errorf("write tool not found: %s", writeOptions.writeToolPath)
	}

	partition, err := options.ResolvePartitionOptions()
	if err != nil {
		return fmt.Errorf("Error resolving partition scheme: %s\n", err)
	}

	cmd.Printf("Writing image file: %s\n", writeImagePath)
	err = write.WriteImageWithSerial(writeOptions.writeToolPath, writeImagePath, write.WriteImageOptions{
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

type ParsedWriteOptions struct {
	writeToolPath       string
	partitionSchemaPath string
	chip                string
	serialPort          string
	serialBaud          int
	flashMode           string
	flashFreq           string
	flashSize           string
	noVerify            bool
}

var (
	writeOptions ParsedWriteOptions
)

func init() {
	rootCmd.AddCommand(writeCmd)

	writeCmd.Flags().StringVarP(&writeOptions.writeToolPath, "tool", "t", "esptool.py", "Path to esptool.py")
	options.AddPartitionOptions(writeCmd)
	AddWriteCommandOptions(writeCmd)
}

func AddWriteCommandOptions(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&writeOptions.chip, "chip", "c", "auto", "ESP32 chip type")
	cmd.Flags().StringVarP(&writeOptions.serialPort, "port", "P", "", "Serial port")
	cmd.Flags().IntVarP(&writeOptions.serialBaud, "baud", "b", 0, "Serial baud rate")
	cmd.MarkFlagRequired("baud")
	cmd.Flags().StringVarP(&writeOptions.flashMode, "flash_mode", "m", "keep", "Flash mode")
	cmd.Flags().StringVarP(&writeOptions.flashFreq, "flash_freq", "f", "keep", "Flash frequency")
	cmd.Flags().StringVarP(&writeOptions.flashSize, "flash_size", "s", "detect", "Flash size")
	cmd.Flags().BoolVar(&writeOptions.noVerify, "no-verify", false, "don't verify flash")
}
