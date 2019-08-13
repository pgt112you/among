package tool

import (
	"encoding/binary"
	"fmt"
	"net"
)

func convertByte2UInt(b []byte) uint32 {
	//bytesBuffer := bytes.NewBuffer(b)
	//var x uint32
	//binary.Read(bytesBuffer, binary.BigEndian, &x)
	//return x
	//return binary.BigEndian.Uint32(b)
	return binary.LittleEndian.Uint32(b)
}

func getPackLen(b []byte) uint32 {
	rawLen := convertByte2UInt(b)
	fmt.Println("rawlen is", rawLen)
	var mask uint32 = 16777215
	len := rawLen & mask
	fmt.Println("len is ", len)
	return len
}

func GetData(c net.Conn, ff string) []byte {
	headerBuf := make([]byte, 4)
	packLen, err := c.Read(headerBuf)
	if err != nil {
		fmt.Println("getdata read lenbuf error", ff, err)
		return nil
	}
	if packLen <= 0 {
		fmt.Println("getdata pack len <=0")
		return nil
	}
	len := getPackLen(headerBuf)
	payloadBuf := make([]byte, len)
	packLen, err = c.Read(payloadBuf)
	if err != nil {
		fmt.Println("getdata read databuf error", ff, err)
		return nil
	}
	if packLen <= 0 {
		fmt.Println("getdata data len <=0")
		return nil
	}
	packBuf := make([]byte, 4+len)
	fmt.Printf("%s package len is %d\n", ff, 4+len)
	copy(packBuf, headerBuf)
	copy(packBuf[4:], payloadBuf)
	fmt.Printf("%s packbuf is %v\n", ff, packBuf)
	return packBuf
}
