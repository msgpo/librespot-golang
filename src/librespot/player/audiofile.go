package player

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"librespot/connection"
	"math"
)

const kChunkSize = 32768

type AudioFile struct {
	Size    uint32
	Chunks  map[int]bool
	Data    []byte
	FileId  []byte
	Channel *Channel
	Stream  connection.PacketStream
}

func NewAudioFile(fileId []byte, channel *Channel, stream connection.PacketStream) *AudioFile {
	return &AudioFile{
		Channel: channel,
		FileId:  fileId,
		Stream:  stream,
	}
}

func (a *AudioFile) Load() error {
	// Request audio data

	for i := 0; i < a.TotalChunks(); i++ {
		err := a.Stream.SendPacket(0x8, buildAudioChunkRequest(a.Channel.num, a.FileId, uint32(i*kChunkSize), uint32((i+1)*kChunkSize)))

		if err != nil {
			return err
		}

	}

	return nil
}

func (a *AudioFile) HasChunk(index int) bool {
	has, ok := a.Chunks[index]
	return has && ok
}

func (a *AudioFile) TotalChunks() int {
	return int(math.Ceil(float64(a.Size) / float64(kChunkSize)))
}

func (a *AudioFile) PutEncryptedChunk(index int, data []byte) {
	byteIndex := index * kChunkSize
	decryptedData := decryptAudio(data)

	copy(a.Data[byteIndex:], decryptedData)
	a.Chunks[index] = true
}

func buildAudioChunkRequest(channel uint16, fileId []byte, start uint32, end uint32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, channel)
	binary.Write(buf, binary.BigEndian, uint8(0x0))
	binary.Write(buf, binary.BigEndian, uint8(0x1))
	binary.Write(buf, binary.BigEndian, uint16(0x0000))
	binary.Write(buf, binary.BigEndian, uint32(0x00000000))
	binary.Write(buf, binary.BigEndian, uint32(0x00009C40))
	binary.Write(buf, binary.BigEndian, uint32(0x00020000))
	buf.Write(fileId)
	binary.Write(buf, binary.BigEndian, start)
	binary.Write(buf, binary.BigEndian, end)

	fmt.Printf("%x", buf.Bytes())

	return buf.Bytes()
}
