package mcprotocol

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

// Conn wraps a net.Conn and provides Minecraft-style packet reading/writing.
type Conn struct {
	net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

// NewConn creates a new Conn wrapping the given net.Conn.
func NewConn(conn net.Conn) *Conn {
	return &Conn{
		Conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

// ReadNext reads the next Minecraft packet from the connection.
// A packet starts with a VarInt length (including packet ID) followed by packet ID and data.
func (c *Conn) ReadNext(ctx context.Context) (int32, *Buffer, error) {
	// Context support using a deadline trigger.
	if ctx.Err() != nil {
		return 0, nil, ctx.Err()
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			c.Conn.SetReadDeadline(time.Now())
		case <-done:
		}
	}()

	// Read packet length (VarInt)
	length, err := c.readVarInt()
	if err != nil {
		if ctx.Err() != nil {
			return 0, nil, ctx.Err()
		}
		return 0, nil, err
	}

	if length < 0 {
		return 0, nil, fmt.Errorf("negative packet length: %d", length)
	}

	// Read packet data (including ID)
	data := make([]byte, length)
	_, err = io.ReadFull(c.reader, data)
	if err != nil {
		if ctx.Err() != nil {
			return 0, nil, ctx.Err()
		}
		return 0, nil, err
	}

	// Create a temporary buffer to parse Packet ID
	tmpBuf := NewBuffer(data)
	packetID, err := tmpBuf.ReadVarInt()
	if err != nil {
		return 0, nil, err
	}

	// Reset deadline after operation
	c.Conn.SetReadDeadline(time.Time{})

	// Return Packet ID and the rest of the data as a new Buffer
	return packetID, NewBuffer(data[tmpBuf.Pos:]), nil
}

// WriteNext writes a Minecraft packet to the connection.
// It automatically prepends the VarInt length (including packet ID).
func (c *Conn) WriteNext(ctx context.Context, packetID int32, buf *Buffer) error {
	// Context support using a deadline trigger.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			c.Conn.SetWriteDeadline(time.Now())
		case <-done:
		}
	}()

	// Create a temporary buffer for the packet body (ID + Data)
	packetBody := NewBuffer(nil)
	packetBody.WriteVarInt(packetID)
	packetBody.Data = append(packetBody.Data, buf.Data...)

	// Write total length
	err := c.writeVarInt(int32(len(packetBody.Data)))
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}

	// Write packet body
	_, err = c.writer.Write(packetBody.Data)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}

	// Flush to ensure it's sent
	err = c.writer.Flush()
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}

	// Reset deadline after operation
	c.Conn.SetWriteDeadline(time.Time{})

	return nil
}

// readVarInt reads a VarInt from the buffered reader.
func (c *Conn) readVarInt() (int32, error) {
	var value uint32
	var position uint
	for {
		currentByte, err := c.reader.ReadByte()
		if err != nil {
			return 0, err
		}
		value |= uint32(currentByte&0x7F) << position
		if (currentByte & 0x80) == 0 {
			break
		}
		position += 7
		if position >= 32 {
			return 0, errors.New("VarInt is too big")
		}
	}
	return int32(value), nil
}

// writeVarInt writes a VarInt to the buffered writer.
func (c *Conn) writeVarInt(v int32) error {
	uValue := uint32(v)
	for {
		if (uValue & ^uint32(0x7F)) == 0 {
			return c.writer.WriteByte(byte(uValue))
		}
		err := c.writer.WriteByte(byte((uValue & 0x7F) | 0x80))
		if err != nil {
			return err
		}
		uValue >>= 7
	}
}
