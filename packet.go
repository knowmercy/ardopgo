package main

var PacketType = map[string]byte{
	"PING":       0x3E,
	"PINGACK":    0x3E,
	"BREAK":      0x23,
	"IDLEFRAME":  0x24,
	"OCONREQ200": 0x25,
}

// Frame Types
var BREAK = Packet{
	messageType: PacketType["BREAK"],
}

type END struct {
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

type Packet struct {
	messageType byte
	sessionId   byte
}

func (packet *Packet) serialize() []byte {
	return []byte{}
}

// This will build a ping packet using the source and destination call sign.
// The result is reutrned as a byte array.
func ping() []byte {
	pingPacket := Packet{
		messageType: PacketType["PING"],
		sessionId:   0xFF,
	}

	return pingPacket.serialize()
}
