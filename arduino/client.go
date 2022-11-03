package arduino

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os/exec"
)

func DialClient() {
	address := RunArduinoCli()
	log.Printf("dialing arduino-cli daemon at %s\n", address)

	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connecting to arduino-cli rpc server, you can start it by running `arduino-cli daemon`")
	}
	defer conn.Close()
	log.Printf("connect established")

	client := commands.NewArduinoCoreServiceClient(conn)
	//settingsClient := settings.NewSettingsServiceClient(conn)

	callVersion(client)
}

func RunArduinoCli() string {
	ch := make(chan DaemonAddress)
	go startupArduinoDaemon(ch)

	address := <-ch

	return fmt.Sprintf("%s:%s", address.IP, address.Port)
}

type DaemonAddress struct {
	IP   string
	Port string
}

func startupArduinoDaemon(done chan<- DaemonAddress) {
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

func callVersion(client commands.ArduinoCoreServiceClient) {
	versionResp, err := client.Version(context.Background(), &commands.VersionRequest{})
	if err != nil {
		log.Fatalf("Error getting version: %s", err)
	}

	log.Printf("arduino-cli version: %v", versionResp.GetVersion())
}
