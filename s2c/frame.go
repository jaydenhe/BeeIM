package s2c

import (
	"io"
	"encoding/binary"
	"bytes"
	"log"
)


type Frame struct
{
	messageType uint16
	protobufData []byte
}

func (f *Frame) Marshal()([]byte,error){
	
	frameLen := uint32(len(f.protobufData) + 4 + 2)

	buffer := new(bytes.Buffer)
	
	binary.Write(buffer,binary.BigEndian,frameLen)
	binary.Write(buffer,binary.BigEndian,f.messageType)
	buffer.Write(f.protobufData)
	
	return buffer.Bytes(),nil
}




func ReadAFrame(in io.Reader) (frame *Frame,err error) {
	
	frame = new(Frame)
	frameLenAndType := make([]byte,6)
		
	//header
	n,err := io.ReadFull(in,frameLenAndType)
	if n == 0 && err == io.EOF {
		return nil,err
	}else if err != nil {
		log.Println("error read framelen:",err)
		return nil,err
	}
	frameLen:=binary.BigEndian.Uint32(frameLenAndType[:4])
	frame.messageType = binary.BigEndian.Uint16(frameLenAndType[4:])

	//data
	frame.protobufData = make([]byte,(frameLen - 6))
	
	n,err = io.ReadFull(in,frame.protobufData)

	if err != nil {
		log.Println("error receiving msg:",err)
		return nil,err
	}
	return frame,nil
}




















