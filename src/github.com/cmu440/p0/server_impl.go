// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"fmt"
	"io"
	"net"
	"sync"
)

const MAX_CLIENT_NUM = 10
const MAX_MESSAGE_NUM = 100
const RECV_BUF_LEN = 1024

type multiEchoServer struct {
	// TODO: implement this!
	clients [MAX_CLIENT_NUM] *client
	mutex sync.Mutex
	cnt int
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!
	return &multiEchoServer{
		cnt : 0,
	}
}

func (mes *multiEchoServer) Start(port int) error {
	// TODO: implement this!
	listener, err := net.Listen("tcp", "0.0.0.0:" + fmt.Sprintf("%d", port))
	
	if err != nil {	
		panic("error listening:" + err.Error())
	}
	fmt.Println("Starting the server")

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("Error accept:" + err.Error())
		}
		fmt.Println("Accepted the Connection :", conn.RemoteAddr())
		
		//init client
		mes.mutex.Lock()
		
		cl := &client{
			conn: conn,
			connIndex: mes.cnt,
			startIndex: 0,
			endIndex: 1,
		}

		mes.clients[mes.cnt] = cl
		mes.cnt ++

		mes.mutex.Unlock()

		go mes.EchoServer(cl)
	}
}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!
	for _, cl := range mes.clients {
		mes.closeClient(cl)
	}
	// go routines should be signaled to return

	return
}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!
	return mes.cnt
}

// TODO: add additional methods/functions below!
type client struct{
	conn net.Conn
	connIndex int

	Messages [MAX_MESSAGE_NUM][] byte
	mutex sync.Mutex
	startIndex int
	endIndex int
}

func (mes *multiEchoServer) sendToAll(s []byte) {
	for _, cl := range mes.clients {
		// not full
		cl.mutex.Lock()

		if (cl.endIndex - cl.startIndex) < MAX_MESSAGE_NUM {
			cl.Messages[cl.endIndex % MAX_MESSAGE_NUM] = s
			cl.endIndex ++
		}

		cl.mutex.Unlock()
	}
}

func (mes *multiEchoServer) closeClient(cl *client) {
	cl.conn.Close()

	// del client from array
	index := cl.connIndex
	mes.mutex.Lock()
	for i := index; i < mes.cnt - 1; i ++ {
		mes.clients[i] = mes.clients[i + 1]
    }
	mes.clients[mes.cnt - 1] = nil
	mes.cnt --
	mes.mutex.Unlock()

	//free struct client
}

func (mes *multiEchoServer) EchoServer(cl *client) {
	fmt.Printf("Client %d established a connection \n", cl.connIndex)
	conn := cl.conn
	buf := make([]byte, RECV_BUF_LEN)

	defer mes.closeClient(cl)

	go func(){
		//check Messages
		for {
			//only one procedure call
			//do not need lock here
			if cl.endIndex != cl.startIndex {
				conn.Write(cl.Messages[cl.startIndex])
				cl.startIndex ++
			}
		}
	}()

	for {
		n, err := conn.Read(buf)
		switch err {
		case nil:
			fmt.Printf("Client %d send a message: %s \n", cl.connIndex, buf)
			mes.sendToAll(buf[0:n])
		case io.EOF:
			fmt.Printf("Warning: End of data: %s \n", err)
			return
		default:
			fmt.Printf("Error: Reading data : %s \n", err)
			return
		}
	}
}