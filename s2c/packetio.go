/*
author:jaydenhe
email:jayden.he@qq.com
 */
package s2c

import (
	"io"
	"bufio"
	"encoding/binary"
	"errors"
)

/*
packet struct
--------------------------------------------------------------------------------
4bytes(header[0xDEADBEEF])  + 4bytes(packetType[uint32]) + 4bytes(length[uint32{limit:<uint32/2}]) + n bytes(data)


*/


const (
	defaultPacketBufSize = 2048
	maxPacketLength      = 0xEFFFFFFF
)

var (
	ErrInvalidPacketHat    = errors.New("pakcetio:Invalid PacketHat")
	ErrInvalidPacketLength = errors.New("packetio:invalid Packet length")
)

var (
	PacketHat = []byte{byte(0xDE), byte(0xAD), byte(0xBE), byte(0xEF)}
)

// Reader implements buffering for an io.Reader object.
type PacketReader struct {
	br *bufio.Reader
}

// NewReader returns a new Reader whose buffer has the default size.
func NewPacketReader(rd io.Reader) *PacketReader {
	r := &PacketReader{}
	r.br = bufio.NewReaderSize(rd, defaultPacketBufSize)
	return r
}

//check packet hat
func (r *PacketReader) checkPacketHat() (bool, error) {

	//read 4 bytes from bufio.reader
	for v := range PacketHat {
		b ,err := r.br.ReadByte()

		if err != nil {
			return false,err
		}

		if byte(v) != b {
			return false,nil
		}

	}

	return true, nil
}

//read packetType from io stream
func (r *PacketReader) readPacketType() (packetType uint32, err error) {
	buf := make([]byte, 4)

	hasRead := int(0)
	for {
		n, err := r.br.Read(buf[hasRead:])
		if err != nil {
			return 0, err
		}

		hasRead += n
		if hasRead >= len(buf) {
			break
		}

	}

	packetType = binary.BigEndian.Uint32(buf)

	return packetType, nil
}

//read data length
func (r *PacketReader) readPacketLength() (packetSize uint32, err error) {

	buf := make([]byte, 4)

	//read 4 bytes
	hasRead := int(0)
	for {
		n, err := r.br.Read(buf[hasRead:])

		if err != nil {
			return 0, err
		}

		hasRead += n
		if hasRead >= len(buf) {
			break
		}
	}

	//decoding packet size
	packetSize = binary.BigEndian.Uint32(buf)

	return packetSize, nil
}

func (r *PacketReader) readPacketData(data []byte) (err error) {

	//read len(data) bytes from io stream
	hasRead := uint32(0)
	for {
		n, err := r.br.Read(data[hasRead:])
		if err != nil {
			return err
		}

		hasRead += uint32(n)

		if hasRead >= uint32(len(data)) {
			break
		}
	}

	return nil
}

func (r *PacketReader) ReadAPacket() (packet *Packet, err error) {

	packet = NewPacket()

	//check packet hat
	hasHat, err := r.checkPacketHat()

	if (!hasHat) {
		return nil, ErrInvalidPacketHat
	}

	//check packet type
	packetType, err := r.readPacketType()

	if err != nil {
		return nil, err
	}

	packet.SetType(packetType)

	//check packet length
	packetLength, err := r.readPacketLength()

	if err != nil {
		return nil, err
	}

	if packetLength > maxPacketLength {
		return nil, ErrInvalidPacketLength
	}

	//alloc packet data buffer
	packetData := make([]byte, packetLength)

	//read packet data
	err = r.readPacketData(packetData)

	if err != nil {
		return nil, err
	}

	return packet, nil
}

// Writer implements buffering for an io.Reader object.
type PacketWriter struct{
	bw *bufio.Writer
}

// NewReader returns a new Reader whose buffer has the default size.
func NewPacketWriter(w io.Writer) *PacketWriter {
	wt := &PacketWriter{}
	wt.bw = bufio.NewWriterSize(w, defaultPacketBufSize)
	return wt
}

func (w *PacketWriter) WriteAPacket(packet *Packet) (err error) {

	buf := make([]byte, 4)

	//write packet hat
	_, err = w.bw.Write(PacketHat)
	if err != nil {
		return err
	}

	//write packet type
	binary.BigEndian.PutUint32(buf, packet.GetType())
	_, err = w.bw.Write(buf)
	if err != nil {
		return err
	}

	//write packet length
	binary.BigEndian.PutUint32(buf, uint32(len(packet.GetData())))
	_, err = w.bw.Write(buf)
	if err != nil {
		return err
	}

	//write packet data
	_, err = w.bw.Write(packet.GetData())
	if err != nil {
		return err
	}

	return nil
}

func (w *PacketWriter) Flush() (error) {
	return w.bw.Flush()
}

// buffered input and output

// ReadWriter stores pointers to a Reader and a Writer.
// It implements io.ReadWriter.
type PacketReadWriter struct {
	*PacketReader
	*PacketWriter
}

// NewReadWriter allocates a new ReadWriter that dispatches to r and w.
func NewPacketReadWriter(r *PacketReader, w *PacketWriter) *PacketReadWriter {
	return &PacketReadWriter{r, w}
}
