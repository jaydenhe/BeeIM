package s2c

import (

)

type Packet struct{
	packetType uint32
	packetData []byte
}

func NewPacket() (packet *Packet) {
	packet = &Packet{0, nil}
	return packet
}

func (self *Packet) SetType(t uint32) {
	self.packetType = t
}

func (self *Packet) GetType() (t uint32) {
	return self.packetType
}

func (self *Packet) SetData(data []byte) {
	self.packetData = data
}

func (self *Packet) GetData() (data []byte) {
	return self.packetData
}

