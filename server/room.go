package server

import (
	"encoding/binary"
	"net"
	"sync"
	"time"
)

const arrSize = 1000

type Room struct {
	endSignal chan bool
	ticker    time.Ticker
	Members   map[string]*Member
}

func (room *Room) AppendMember(member *Member) error {
	if member == nil {
		panic("whoa")
	}
	if _, found := room.Members[member.id]; !found {
		room.Members[member.id] = member
	} else {
		panic("whoa")
	}
	go member.ListenMic()
	return nil
}

func (room *Room) RemoveMember(member *Member) error {
	delete(room.Members, member.id)
	return nil
}

func (room *Room) Serve() error {
	go func() {
		for {
			select {
			case <-room.endSignal:
				//close the room
			case <-room.ticker.C:
				req := make(map[string][]float32)
				for _, m := range room.Members {
					req[m.id] = m.inputChunk
				}
				resp := Multiplex(req)
				for _, m := range room.Members {
					_ = binary.Write(m.connection, binary.BigEndian, resp[m.id])
				}
			}
		}
	}()
	return nil
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
	// call this function in every tick
	var tempMemChunk []float32
	for {
		tempMemChunk = <-member.micBuffer
	}
	member.inputChunkLock.Lock()
	member.inputChunk = tempMemChunk
	member.inputChunkLock.Unlock()
}
