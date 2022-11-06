package arduino

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type DaemonAddress struct {
	IP   string
	Port string
}

func (address DaemonAddress) ToString() string {
	return fmt.Sprintf("%s:%s", address.IP, address.Port)
}

func StartupArduinoDaemon(done chan<- DaemonAddress) {
	cmd := exec.Command("arduino-cli", "daemon", "--format", "JSONMini", "--log-file", "./arduino-cli.log", "--log-level", "trace", "--log-format", "json")

	//指定しないとdaemonがすぐに閉じてしまう
	_, err := cmd.StdinPipe()

	log.Println("starting arduino-cli daemon")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	log.Printf("process started pid %d", cmd.Process.Pid)

	var address DaemonAddress
	if err := json.NewDecoder(stdout).Decode(&address); err != nil {
		log.Fatal(err)
	}

	log.Printf("arduino-cli daemon started on %s:%s", address.IP, address.Port)
	done <- address

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	log.Printf("process finished with error = %v", err)
}
