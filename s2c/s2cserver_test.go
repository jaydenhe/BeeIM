package s2c

import (
	"testing"
	"net"
	"fmt"
	"bytes"
	"time"
)

func TestLogin(t *testing.T) {

	s2c_server := CreateServer()

	fmt.Println("Running on \n")
	go s2c_server.Start("localhost:1114")
	fmt.Println("s2c_server start!")

	conn, err := net.Dial("tcp4", "localhost:1114")

	if err != nil {
		t.Error("Dial error:", err)
	}

	time.Sleep(time.Second * 1)

	packet := NewPacket()
	packet.SetType(testPacketType)
	packet.SetData(testPacketData)

	packetWR := NewPacketReadWriter(NewPacketReader(conn), NewPacketWriter(conn))


	err = packetWR.WriteAPacket(packet)

	if err != nil {
		t.Error("packetWriter error:", err)
	}

	rPacket, err := packetWR.ReadAPacket()

	if err != nil {
		t.Error("packetReader error:", err)
	}

	if rPacket.GetType() != testPacketType && bytes.Compare(rPacket.GetData(), testPacketData) != 0 {
		t.Error("rPacket != packet")
	}

	conn.Close()

	s2c_server.Stop()
}



















