package linkedveil

import(
	"fmt"
	"log"

	"linkedveil/internal/transport"
)

func main() { 
	udp, err := transport.Listen(":9999")
	if err != nil { 
		log.Fatal(err)
	}
	defer udp.Close()

	fmt.Println("listening on :9999")

	buf := make([]byte, 2048)
	for {
		n, addr, err := udp.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("received %q from %s\n", buf[:n], addr)
	}
}

