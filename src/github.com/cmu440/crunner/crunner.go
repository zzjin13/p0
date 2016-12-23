package main

import (
	"fmt"
	"net"
	"time"
)

const (
	defaultHost  = "localhost"
	defaultPort  = 9999
	RECV_BUF_LEN = 1024
)

// To test your server implementation, you might find it helpful to implement a
// simple 'client runner' program. The program could be very simple, as long as
// it is able to connect with and send messages to your server and is able to
// read and print out the server's echoed response to standard output. Whether or
// not you add any code to this file will not affect your grade.
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:" + fmt.Sprintf("%d", defaultPort))
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	buf := make([]byte, RECV_BUF_LEN)
	for i := 0; i < 5; i++ {
		//send message
		msg := fmt.Sprintf("Hello World, %03d", i)
		n, err := conn.Write([]byte(msg))
		if err != nil {
			println("Write Buffer Error:", err.Error())
			break
		}
		fmt.Printf("client send: %s\n", msg)
		
		//receive message from server
		n, err = conn.Read(buf)
		if err != nil {
			println("Read Buffer Error:", err.Error())
			break
		}
		fmt.Printf("client receive: %s\n", string(buf[0:n]))
		
		//wait
		time.Sleep(time.Second)
	}
}
