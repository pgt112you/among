package conn

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func sendData(conn *net.Conn, data []byte) (int, error) {
	packLen := len(data)
	n := 0
	for n < packLen {
		tempN, err := (*conn).Write(data[n:])
		if err != nil {
			return n, err
		}
		n += tempN
	}
	return n, nil
}

func convertByte2UInt(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func getPackLen(b []byte) uint32 {
	rawLen := convertByte2UInt(b)
	var mask uint32 = 16777215
	len := rawLen & mask
	return len
}

func recvPack(conn *net.Conn, num uint32) ([]byte, int, error) {
	buf := make([]byte, num)
	var n int
	for uint32(n) < num {
		packLen, err := (*conn).Read(buf[n:])
		if err != nil {
			if err == io.EOF {
				fmt.Println("recv eof, close conn")
				(*conn).Close()
				return buf, n, err
			}
			fmt.Printf("recv package error %s\n", err)
			return nil, 0, err
		}
		n += packLen
	}
	return buf, n, nil
}

func recvData(conn *net.Conn) ([]byte, int, error) {
	var hlen uint32 = 4
	headerBuf, _, err := recvPack(conn, hlen)
	if err != nil {
		fmt.Println("recvdata read headerbuf error", err)
		return nil, 0, err
	}
	plen := getPackLen(headerBuf)
	payloadBuf, _, err := recvPack(conn, plen)
	if err != nil {
		fmt.Println("recvdata read payloadbuf error", err)
		return nil, 0, err
	}
	packLen := hlen + plen
	packBuf := make([]byte, packLen)
	copy(packBuf, headerBuf)
	copy(packBuf[hlen:], payloadBuf)
	return packBuf, int(packLen), err
}
