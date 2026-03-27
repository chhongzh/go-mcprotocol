package mcprotocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/google/uuid"
)

// Buffer is a helper to read/write Minecraft data types.
// It uses a byte slice as the underlying storage.
type Buffer struct {
	Data []byte
	Pos  int
}

// NewBuffer creates a new Buffer with the given data.
func NewBuffer(data []byte) *Buffer {
	return &Buffer{Data: data, Pos: 0}
}

// Reset clears the buffer and resets the position.
func (b *Buffer) Reset() {
	b.Data = b.Data[:0]
	b.Pos = 0
}

// Len returns the number of bytes in the buffer.
func (b *Buffer) Len() int {
	return len(b.Data)
}

// Bytes returns the underlying data slice.
func (b *Buffer) Bytes() []byte {
	return b.Data
}

// --- Basic Types ---

// WriteBoolean writes a boolean (1 byte) to the buffer.
func (b *Buffer) WriteBoolean(v bool) {
	if v {
		b.Data = append(b.Data, 0x01)
	} else {
		b.Data = append(b.Data, 0x00)
	}
}

// ReadBoolean reads a boolean (1 byte) from the buffer.
func (b *Buffer) ReadBoolean() (bool, error) {
	if b.Pos >= len(b.Data) {
		return false, io.ErrUnexpectedEOF
	}
	v := b.Data[b.Pos] != 0x00
	b.Pos++
	return v, nil
}

// WriteByte writes a signed byte (1 byte) to the buffer.
func (b *Buffer) WriteByte(v int8) {
	b.Data = append(b.Data, byte(v))
}

// ReadByte reads a signed byte (1 byte) from the buffer.
func (b *Buffer) ReadByte() (int8, error) {
	if b.Pos >= len(b.Data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := int8(b.Data[b.Pos])
	b.Pos++
	return v, nil
}

// WriteUByte writes an unsigned byte (1 byte) to the buffer.
func (b *Buffer) WriteUByte(v uint8) {
	b.Data = append(b.Data, v)
}

// ReadUByte reads an unsigned byte (1 byte) from the buffer.
func (b *Buffer) ReadUByte() (uint8, error) {
	if b.Pos >= len(b.Data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := b.Data[b.Pos]
	b.Pos++
	return v, nil
}

// WriteShort writes a signed short (2 bytes, big-endian) to the buffer.
func (b *Buffer) WriteShort(v int16) {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], uint16(v))
	b.Data = append(b.Data, buf[:]...)
}

// ReadShort reads a signed short (2 bytes, big-endian) from the buffer.
func (b *Buffer) ReadShort() (int16, error) {
	if b.Pos+2 > len(b.Data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := int16(binary.BigEndian.Uint16(b.Data[b.Pos : b.Pos+2]))
	b.Pos += 2
	return v, nil
}

// WriteUShort writes an unsigned short (2 bytes, big-endian) to the buffer.
func (b *Buffer) WriteUShort(v uint16) {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], v)
	b.Data = append(b.Data, buf[:]...)
}

// ReadUShort reads an unsigned short (2 bytes, big-endian) from the buffer.
func (b *Buffer) ReadUShort() (uint16, error) {
	if b.Pos+2 > len(b.Data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := binary.BigEndian.Uint16(b.Data[b.Pos : b.Pos+2])
	b.Pos += 2
	return v, nil
}

// WriteInt writes a signed int (4 bytes, big-endian) to the buffer.
func (b *Buffer) WriteInt(v int32) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(v))
	b.Data = append(b.Data, buf[:]...)
}

// ReadInt reads a signed int (4 bytes, big-endian) from the buffer.
func (b *Buffer) ReadInt() (int32, error) {
	if b.Pos+4 > len(b.Data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := int32(binary.BigEndian.Uint32(b.Data[b.Pos : b.Pos+4]))
	b.Pos += 4
	return v, nil
}

// WriteLong writes a signed long (8 bytes, big-endian) to the buffer.
func (b *Buffer) WriteLong(v int64) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(v))
	b.Data = append(b.Data, buf[:]...)
}

// ReadLong reads a signed long (8 bytes, big-endian) from the buffer.
func (b *Buffer) ReadLong() (int64, error) {
	if b.Pos+8 > len(b.Data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := int64(binary.BigEndian.Uint64(b.Data[b.Pos : b.Pos+8]))
	b.Pos += 8
	return v, nil
}

// WriteFloat writes a float (4 bytes, big-endian) to the buffer.
func (b *Buffer) WriteFloat(v float32) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], math.Float32bits(v))
	b.Data = append(b.Data, buf[:]...)
}

// ReadFloat reads a float (4 bytes, big-endian) from the buffer.
func (b *Buffer) ReadFloat() (float32, error) {
	if b.Pos+4 > len(b.Data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := math.Float32frombits(binary.BigEndian.Uint32(b.Data[b.Pos : b.Pos+4]))
	b.Pos += 4
	return v, nil
}

// WriteDouble writes a double (8 bytes, big-endian) to the buffer.
func (b *Buffer) WriteDouble(v float64) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(v))
	b.Data = append(b.Data, buf[:]...)
}

// ReadDouble reads a double (8 bytes, big-endian) from the buffer.
func (b *Buffer) ReadDouble() (float64, error) {
	if b.Pos+8 > len(b.Data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := math.Float64frombits(binary.BigEndian.Uint64(b.Data[b.Pos : b.Pos+8]))
	b.Pos += 8
	return v, nil
}

// --- Variable-length Types ---

// WriteVarInt writes a variable-length integer (up to 5 bytes) to the buffer.
func (b *Buffer) WriteVarInt(v int32) {
	uValue := uint32(v)
	for {
		if (uValue & ^uint32(0x7F)) == 0 {
			b.Data = append(b.Data, byte(uValue))
			return
		}
		b.Data = append(b.Data, byte((uValue&0x7F)|0x80))
		uValue >>= 7
	}
}

// ReadVarInt reads a variable-length integer (up to 5 bytes) from the buffer.
func (b *Buffer) ReadVarInt() (int32, error) {
	var value uint32
	var position uint
	for {
		if b.Pos >= len(b.Data) {
			return 0, io.ErrUnexpectedEOF
		}
		currentByte := b.Data[b.Pos]
		b.Pos++
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

// WriteVarLong writes a variable-length long (up to 10 bytes) to the buffer.
func (b *Buffer) WriteVarLong(v int64) {
	uValue := uint64(v)
	for {
		if (uValue & ^uint64(0x7F)) == 0 {
			b.Data = append(b.Data, byte(uValue))
			return
		}
		b.Data = append(b.Data, byte((uValue&0x7F)|0x80))
		uValue >>= 7
	}
}

// ReadVarLong reads a variable-length long (up to 10 bytes) from the buffer.
func (b *Buffer) ReadVarLong() (int64, error) {
	var value uint64
	var position uint
	for {
		if b.Pos >= len(b.Data) {
			return 0, io.ErrUnexpectedEOF
		}
		currentByte := b.Data[b.Pos]
		b.Pos++
		value |= uint64(currentByte&0x7F) << position
		if (currentByte & 0x80) == 0 {
			break
		}
		position += 7
		if position >= 64 {
			return 0, errors.New("VarLong is too big")
		}
	}
	return int64(value), nil
}

// --- String and Data ---

// WriteString writes a UTF-8 string prefixed by its length as a VarInt.
func (b *Buffer) WriteString(v string) {
	b.WriteVarInt(int32(len(v)))
	b.Data = append(b.Data, []byte(v)...)
}

// ReadString reads a UTF-8 string prefixed by its length as a VarInt.
func (b *Buffer) ReadString() (string, error) {
	length, err := b.ReadVarInt()
	if err != nil {
		return "", err
	}
	if length < 0 {
		return "", fmt.Errorf("negative string length: %d", length)
	}
	if b.Pos+int(length) > len(b.Data) {
		return "", io.ErrUnexpectedEOF
	}
	v := string(b.Data[b.Pos : b.Pos+int(length)])
	b.Pos += int(length)
	return v, nil
}

// WriteUUID writes a 128-bit UUID to the buffer.
func (b *Buffer) WriteUUID(v uuid.UUID) {
	b.Data = append(b.Data, v[:]...)
}

// ReadUUID reads a 128-bit UUID from the buffer.
func (b *Buffer) ReadUUID() (uuid.UUID, error) {
	if b.Pos+16 > len(b.Data) {
		return uuid.Nil, io.ErrUnexpectedEOF
	}
	var v uuid.UUID
	copy(v[:], b.Data[b.Pos:b.Pos+16])
	b.Pos += 16
	return v, nil
}

// WriteByteArray writes a byte array prefixed by its length as a VarInt.
func (b *Buffer) WriteByteArray(v []byte) {
	b.WriteVarInt(int32(len(v)))
	b.Data = append(b.Data, v...)
}

// ReadByteArray reads a byte array prefixed by its length as a VarInt.
func (b *Buffer) ReadByteArray() ([]byte, error) {
	length, err := b.ReadVarInt()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("negative byte array length: %d", length)
	}
	if b.Pos+int(length) > len(b.Data) {
		return nil, io.ErrUnexpectedEOF
	}
	v := make([]byte, length)
	copy(v, b.Data[b.Pos:b.Pos+int(length)])
	b.Pos += int(length)
	return v, nil
}
