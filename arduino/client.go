package arduino

import (
	"github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/arduino/arduino-cli/rpc/cc/arduino/cli/settings/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Client struct {
	commands commands.ArduinoCoreServiceClient
	settings settings.SettingsServiceClient
	conn     *grpc.ClientConn
}

func DialGrpcDaemon(dialTarget string) *Client {
	log.Printf("dialing arduino-cli daemon at %s\n", dialTarget)

	conn, err := grpc.Dial(dialTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connecting to arduino-cli rpc server")
	}
	log.Printf("connect established")

	client := commands.NewArduinoCoreServiceClient(conn)
	settingsClient := settings.NewSettingsServiceClient(conn)

	return &Client{
		commands: client,
		settings: settingsClient,
		conn:     conn,
	}
}

func (client *Client) Close() {
	err := client.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}

//
//func RunArduinoCli() string {
//	ch := make(chan DaemonAddress)
//	go startupArduinoDaemon(ch)
//
//	address := <-ch
//
//	return fmt.Sprintf("%s:%s", address.IP, address.Port)
//}
//
//func callVersion(client commands.ArduinoCoreServiceClient) {
//	versionResp, err := client.Version(context.Background(), &commands.VersionRequest{})
//	if err != nil {
//		log.Fatalf("Error getting version: %s", err)
//	}
//
//	log.Printf("arduino-cli version: %v", versionResp.GetVersion())
//}
