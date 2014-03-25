// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type multiEchoServer struct {
	// TODO: implement this!
    ln net.Listener
	clients       []*multiEchoClient
	addClientChan chan net.Conn
	readChan      chan []byte
	closeChan     chan bool
}

type multiEchoClient struct {
	conn      net.Conn
	writeChan chan []byte
	closeChan chan bool
}


// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!
	mes := new(multiEchoServer)
	mes.addClientChan = make(chan net.Conn)
	mes.readChan = make(chan []byte)
	mes.closeChan = make(chan bool)
	return mes
}

func (mes *multiEchoServer) Start(port int) error {
	// TODO: implement this!
	fmt.Println("Start listening.")
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
    mes.ln = ln
	go mes.serverRoutine()
	return err
}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!
	mes.closeChan <- true
    mes.ln.Close()
}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!
	return len(mes.clients)
}

// TODO: add additional methods/functions below!
func (mes *multiEchoServer) serverRoutine() {
	go mes.handleAccept(mes.ln)

	for {
		select {
		case line := <-mes.readChan:
			for _, client := range mes.clients {
				client.writeChan <- line
			}
		case conn := <-mes.addClientChan:
			client := new(multiEchoClient)
			client.conn = conn
            client.writeChan = make(chan []byte, 100)
            client.closeChan = make(chan bool)
			mes.clients = append(mes.clients, client)
			go client.clientRoutine()
			go mes.handleRead(client)
		case <-mes.closeChan:
			for _, client := range mes.clients {
				client.closeChan <- true
			}
			return
		}

	}
}

func (client *multiEchoClient) clientRoutine() {
	for {
		select {
		case line := <-client.writeChan:
			client.conn.Write(line)
		case <-client.closeChan:
			client.conn.Close()
			return
		}
	}
}

func (mes *multiEchoServer) handleAccept(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			break
		}
		mes.addClientChan <- conn
	}
}

func (mes *multiEchoServer) handleRead(client *multiEchoClient) {
	reader := bufio.NewReader(client.conn)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {

		}

		mes.readChan <- line
	}
}

