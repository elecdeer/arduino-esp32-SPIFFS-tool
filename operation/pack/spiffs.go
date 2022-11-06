package pack

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type MakeSpiffsImageOptions struct {
	Partition *PartitionSchema
	PageSize  int
	BlockSize int
}

func MakeSpiffsImage(toolPath string, dist string, source string, options MakeSpiffsImageOptions) error {
	cmdArgs := []string{
		"-c", source,
		"-p", strconv.Itoa(options.PageSize),
		"-b", strconv.Itoa(options.BlockSize),
		"-s", strconv.Itoa(int(options.Partition.Size)),
		dist,
	}

	log.Printf("exec cmd: mkspiffs %s\n", strings.Join(cmdArgs, " "))

	cmd := exec.Command(toolPath, cmdArgs...)
	toolResult, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to exec mkspiffs: %w", err)
	}

	log.Printf("include files:\n")
	for _, line := range strings.SplitAfter(string(toolResult), "\n") {
		log.Printf("  %s", line)
	}

	return nil
}
