package v47

import mcprotocol "github.com/chhongzh/go-mc-protocol"

type Handshake struct {
	Version int32
	Addr    string
	Port    uint16
	Intent  int32
}

func NewHandshake(buf *mcprotocol.Buffer) (*Handshake, error) {
	version, err := buf.ReadVarInt()
	if err != nil {
		return nil, err
	}
	addr, err := buf.ReadString()
	if err != nil {
		return nil, err
	}
	port, err := buf.ReadUShort()
	if err != nil {
		return nil, err
	}
	intent, err := buf.ReadVarInt()
	if err != nil {
		return nil, err
	}

	return &Handshake{version, addr, port, intent}, nil
}
