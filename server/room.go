package server

import (
	"encoding/binary"
	"net"
	"sync"
)

const arrSize = 1000

type Room struct {
	Members map[string]Member
}

type Member struct {
	id             string
	inputChunk     []float32
	inputChunkLock sync.RWMutex
	micBuffer      chan []float32
	connection     *net.TCPConn
}

func (member *Member) ListenMic() {
	// continuously listens from the connection and send it to micBuffer
	var tempMemChunk []float32
	for {
		readErr := binary.Read(member.connection, binary.BigEndian, tempMemChunk)
		if readErr != nil {
			panic(readErr)
		}
		member.micBuffer <- tempMemChunk
	}
}

func (member *Member) MicChunkRegister() {
	// read from the buffer and discard multiple chunks, keeping the latest one only
	var tempMemChunk []float32
	for {
		tempMemChunk = <-member.micBuffer
	}
	member.inputChunkLock.Lock()
	member.inputChunk = tempMemChunk
	member.inputChunkLock.Unlock()
}
