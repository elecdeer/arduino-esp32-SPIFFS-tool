package cmd

import (
	"fmt"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/operation/common"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/operation/pack"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/operation/write"

	"github.com/spf13/cobra"
)

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write a image file to ESP32 flash memory",
	Args:  cobra.ExactArgs(1),
	Run:   RunWrite,
}

func RunWrite(cmd *cobra.Command, args []string) {
	common.ApplyLogStyle()

	writeParam, err := ConstructWriteParams(args, &writeOptions)
	if err != nil {
		cmd.Printf("Error: %s\n", err)
		return
	}

	cmd.Printf("Writing image file: %s\n", writeParam.imagePath)
	err = write.WriteImageWithSerial(writeParam.writeToolPath, writeParam.imagePath, writeParam.options)
	if err != nil {
		cmd.Printf("Error: %s\n", err)
		return
	}
	cmd.PrintErrf("Write image successfully\n")
}

type WriteParam struct {
	options       write.WriteImageOptions
	writeToolPath string
	imagePath     string
}

func ConstructWriteParams(args []string, options *ParsedWriteOptions) (WriteParam, error) {
	writeImagePath := args[0]

	if !common.IsFileExists(writeImagePath) {
		return WriteParam{}, fmt.Errorf("image file not found: %s", writeImagePath)
	}

	if !common.IsFileExists(options.writeToolPath) {
		return WriteParam{}, fmt.Errorf("write tool not found: %s", options.writeToolPath)
	}

	partition, err := pack.ReadPartitionSchemeFile(options.partitionSchemaPath)
	if err != nil {
		return WriteParam{}, fmt.Errorf("Error reading partition scheme file: %s\n", err)
	}

	return WriteParam{
		options: write.WriteImageOptions{
			Partition:  partition,
			Chip:       options.chip,
			SerialBaud: options.serialBaud,
			SerialPort: options.serialPort,
			FlashMode:  options.flashMode,
			FlashFreq:  options.flashFreq,
			FlashSize:  options.flashSize,
			NoVerify:   options.noVerify,
		},
		writeToolPath: options.writeToolPath,
		imagePath:     writeImagePath,
	}, nil
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
	writeCmd.Flags().StringVarP(&writeOptions.partitionSchemaPath, "partition", "p", "", "Partition scheme file path")
	writeCmd.MarkFlagRequired("partition")
	writeCmd.Flags().StringVarP(&writeOptions.chip, "chip", "c", "auto", "ESP32 chip type")
	writeCmd.Flags().StringVarP(&writeOptions.serialPort, "port", "P", "", "Serial port")
	writeCmd.Flags().IntVarP(&writeOptions.serialBaud, "baud", "b", 0, "Serial baud rate")
	writeCmd.MarkFlagRequired("baud")
	writeCmd.Flags().StringVarP(&writeOptions.flashMode, "flash_mode", "m", "keep", "Flash mode")
	writeCmd.Flags().StringVarP(&writeOptions.flashFreq, "flash_freq", "f", "keep", "Flash frequency")
	writeCmd.Flags().StringVarP(&writeOptions.flashSize, "flash_size", "s", "detect", "Flash size")
	writeCmd.Flags().BoolVar(&writeOptions.noVerify, "no-verify", false, "don't verify flash")
}
