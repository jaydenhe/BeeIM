package s2c

import (
	"testing"
	"bytes"
)

var (
	testPacketType  = uint32(2)
	testPacketData  = []byte{0xDD, 0xFF, 0xEE, 0xEE, 0xAB, 0x23, 0x33}
	testPacketSteam = []byte{/*hat*/0xde, 0xad, 0xbe, 0xef, /*type*/0x00, 0x00, 0x00, 0x02,
	/*length*/0x00, 0x00, 0x00, 0x07, /*data*/0xDD, 0xFF, 0xEE, 0xEE, 0xAB, 0x23, 0x33}
)

func TestPacketWriteAndRead(t *testing.T) {

	io := bytes.NewBuffer(nil)
	packet := NewPacket()
	packet.SetType(testPacketType)
	packet.SetData(testPacketData)


	packetWriter := NewPacketWriter(io)
	err := packetWriter.WriteAPacket(packet)
	if err != nil {
		t.Error("packetWriter write error:", err)
	}

	if err = packetWriter.Flush(); err != nil {
		t.Error("packetWrter flush error", err)
	}


	if bytes.Compare(testPacketSteam, io.Bytes()) != 0 {
		t.Error("\ninvalid packet stream buf:", io.Bytes(), "\n            defaultStream:", testPacketSteam)
	}else {
		t.Log("io_buffer:", io.Bytes())
	}

	packetReader := NewPacketReader(io)

	rPacket, err := packetReader.ReadAPacket()

	if rPacket.GetType() != testPacketType && bytes.Compare(rPacket.GetData(), testPacketData) != 0 {
		t.Error("packetReader error:", err)
	}

}
