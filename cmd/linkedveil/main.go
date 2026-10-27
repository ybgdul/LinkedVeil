package linkedveil

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

const socketPath = "/var/run/linkedveil.cok"

func main() { 
	if len(os.Args) < 2 { 
		fmt.Println("Usage linkedveil: <status|reload>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command{
	case "status":
		sendReq("/status")
	case "reload":
		sendReq("/reload")
	default: 
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func sendReq(endpoint string) {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", socketPath)	
			},
		},
	}

	resp, err := client.Get("http://unix" + endpoint)

	if err != nil { 
		fmt.Printf("Error connecting to linkedveil: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}