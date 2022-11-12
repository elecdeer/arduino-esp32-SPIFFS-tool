package write

import (
	"fmt"
	"github.com/elecdeer/arduino-esp32-SPIFFS-tool/operation/pack"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type WriteImageOptions struct {
	Chip       string               //default: auto
	Partition  pack.PartitionSchema //必須
	SerialBaud int                  //必須
	SerialPort string               //指定しない場合は自動検出
	FlashMode  string               //default: keep
	FlashFreq  string               //default: keep
	FlashSize  string               //default: detect
	NoVerify   bool                 //default: true
}

func (o WriteImageOptions) ToCmdOptionsText(imagePath string) []string {
	var args []string
	toolArgs := []string{
		"--chip", o.Chip,
		"--baud", strconv.Itoa(o.SerialBaud),
	}
	if o.SerialPort != "" {
		toolArgs = append(toolArgs, "--port", o.SerialPort)
	}
	args = append(args, toolArgs...)

	cmdArgs := []string{
		"write_flash",
		"--flash_freq", o.FlashFreq,
		"--flash_mode", o.FlashMode,
		"--flash_size", o.FlashSize,
	}
	if o.NoVerify {
		cmdArgs = append(cmdArgs, "--verify", "false")
	}
	args = append(args, cmdArgs...)

	posArgs := []string{
		"0x" + strconv.FormatUint(o.Partition.Offset, 16),
		imagePath,
	}
	args = append(args, posArgs...)

	return args
}

func WriteImageWithSerial(toolPath string, imagePath string, options WriteImageOptions) error {
	args := options.ToCmdOptionsText(imagePath)

	log.Printf("esptool.py: %s", toolPath)
	log.Printf("exec cmd: esptool.py %s\n", strings.Join(args, " "))

	fmt.Printf("========================================\n")
	cmd := exec.Command(toolPath, args...)
	//書き込みログを標準出力と標準エラーに出す
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	fmt.Printf("========================================\n")
	if err != nil {
		return fmt.Errorf("failed to exec esptool.py: %w", err)
	}

	return nil
}
