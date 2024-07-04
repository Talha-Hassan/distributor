package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	buffer := make([]byte, 1024)

	n, err := r.Read(buffer)

	if err != nil {
		return err
	}

	msg.Payload = buffer[:n]
	// fmt.Printf("Message: %+v\n", buffer[:n])

	return nil
}
