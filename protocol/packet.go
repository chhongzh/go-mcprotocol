package protocol

type Packet interface {
	GetPacketID() int
}
