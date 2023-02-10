package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

var (
	tnc     = TNC{}
	version = "0.0.1"
	state   = []string{"Unknown", "Offline", "Disconnected", "ISS", "IRS", "Idle", "FECSend", "FECReceive"}
)

// Frame Types
type END struct {
}
type BREAK struct {
}
type NAK struct {
}
type ACK struct {
}
type DATANAK struct {
}
type DATAACK struct {
}

type CONREQ200 struct {
}
type CONREQ500 struct {
}
type CONREQ1000 struct {
}
type CONREQ2000 struct {
}

type CONACK200 struct {
}
type CONACK500 struct {
}
type CONACK1000 struct {
}
type CONACK2000 struct {
}

// types
type TNC struct {
	control      io.ReadWriteCloser
	dataConn     *net.TCPConn
	state        string
	version      string
	myCall       string
	gridSquare   string
	protocol     string
	arqTimeout   uint16 // seconds
	arqBandwidth uint16
	busy         bool
	connected    bool
	closed       bool
}

func (tnc *TNC) State() string {
	return tnc.state
}

type State struct {
}

// PAT Message handling
func handleRequest(message []byte, conn net.Conn) {
	splitRequest := bytes.Split(message, []byte("\r"))

	msgParts := bytes.Split(splitRequest[0], []byte(" "))
	// This should be the message type
	msgType := string(msgParts[0])

	switch msgType {
	case "INITIALIZE":
		tnc = TNC{
			arqBandwidth: 500,
			version:      version,
		}
		tnc.state = state[2]
		resp := fmt.Sprintf("INITIALIZE\r")
		io.WriteString(conn, resp)
	case "LISTEN":
		resp := fmt.Sprintf("%s\r", msgType)
		io.WriteString(conn, resp)
	case "CWID":
		resp := fmt.Sprintf("%s\r", msgType)
		io.WriteString(conn, resp)
	case "VERSION":
		resp := fmt.Sprintf("%s %s\r", msgType, tnc.version)
		io.WriteString(conn, resp)
	case "MYCALL":
		myCall := string(splitRequest[1])
		tnc.myCall = myCall
		resp := fmt.Sprintf("%s %s\r", msgType, tnc.myCall)
		io.WriteString(conn, resp)
	case "GRIDSQUARE":
		grid := string(splitRequest[1])
		tnc.gridSquare = grid
		resp := fmt.Sprintf("%s %s\r", msgType, tnc.gridSquare)
		io.WriteString(conn, resp)
	case "ARQBW":
		resp := fmt.Sprintf("%s %d\r", msgType, tnc.arqBandwidth)
		io.WriteString(conn, resp)
	case "STATE":
		resp := fmt.Sprintf("STATE %s\r", tnc.State())
		io.WriteString(conn, resp)
	case "PROTOCOLMODE":
		resp := fmt.Sprintf("PROTOCOLMODE\r")
		tnc.protocol = string(splitRequest[1])
		io.WriteString(conn, resp)
	case "ARQTIMEOUT":
		data := binary.BigEndian.Uint16(msgParts[1])
		tnc.arqTimeout = data
		resp := fmt.Sprintf("ARQTIMEOUT %d\r", data)
		io.WriteString(conn, resp)
	default:
		resp := fmt.Sprintf("%s a\r", msgType)
		io.WriteString(conn, resp)
		fmt.Printf("Unknown packet type [%s]\n", msgType)

	}
}

// When PAT makes a connection into this client, it will land here.
func handleConnection(connection net.Conn) {
	defer connection.Close()

	readBuffer := make([]byte, 128)
	for {
		if _, err := connection.Read(readBuffer); err != nil {
			fmt.Printf("ERROR while reading from PAT [%s]\n", err)
		}

		fmt.Printf("Handling connection from [%s] on [%s]\n", connection.RemoteAddr(), connection.LocalAddr())
		fmt.Printf("Got [%d] bytes from PAT [%x]\n", len(readBuffer), readBuffer)

		handleRequest(readBuffer, connection)
	}

}

func startListener(host string) {
	fmt.Printf("Starting up listener on: %s\n", host)

	listener, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Printf("Error establishing listener: [%s]\n", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error Accepting connections on listener: [%s]\n", err)
			return
		}
		go handleConnection(conn)
	}
}

func main() {
	version := "0.0.1"
	fmt.Printf("Starting ARDOP client version [%s]\n", version)

	var listenerWG sync.WaitGroup
	listeners := []string{
		"localhost:8515",
		"localhost:8516",
	}

	for _, listener := range listeners {
		listenerWG.Add(1)
		go startListener(listener)
	}

	listenerWG.Wait()
}
