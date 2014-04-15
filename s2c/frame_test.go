package s2c

import (
	"testing"
	"bytes"
	"encoding/binary"
)


func TestFrameMarshal(t *testing.T){
	f := new(Frame)
	f.messageType = 25
	f.protobufData = []byte{0x1,0x2,0x3}
	

	bytesFrame,err := f.Marshal()
	
	if err != nil {
		t.Fatal("Frame marshal error:",err)
	}

	buffer := bytes.NewBuffer(bytesFrame)
	
	var frameLen uint32
	var frameType uint16
	binary.Read(buffer,binary.BigEndian,&frameLen)
	binary.Read(buffer,binary.BigEndian,&frameType)
	
	if frameLen != 4 + 2 + 3 || frameType != 25{
		t.Fatal("frame Marshal error (9)framelen = ",frameLen,"\n(25)frameTyp=",frameType,"\nf =",f,"\nbytesFrame=",bytesFrame)
	}
	

	conn := bytes.NewBuffer(bytesFrame)
	unmarshalFrame,err :=ReadAFrame(buffer)
	
	if err != nil {
		t.Errorf("ReadAFrame error :",unmarshalFrame)
	}
	
	if unmarshalFrame.messageType != 25 || unmarshalFrame.protobufData != []byte{0x1,0x2,0x3} {
		t.Errorf("ReadAFrame error : (25)messageType=",unmarshalFrame.messageType,"\n protobufData = ",unmarshalFrame.protobufData)
	}

}




















