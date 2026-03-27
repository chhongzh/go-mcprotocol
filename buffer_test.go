package mcprotocol

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
)

func TestBuffer_VarInt(t *testing.T) {
	tests := []struct {
		val      int32
		expected []byte
	}{
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{2, []byte{0x02}},
		{127, []byte{0x7F}},
		{128, []byte{0x80, 0x01}},
		{255, []byte{0xFF, 0x01}},
		{25565, []byte{0xDD, 0xC7, 0x01}},
		{2097151, []byte{0xFF, 0xFF, 0x7F}},
		{2147483647, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x07}},
		{-1, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x0F}},
		{-2147483648, []byte{0x80, 0x80, 0x80, 0x80, 0x08}},
	}

	for _, tc := range tests {
		buf := NewBuffer(nil)
		buf.WriteVarInt(tc.val)
		if !bytes.Equal(buf.Bytes(), tc.expected) {
			t.Errorf("WriteVarInt(%d) = %v, expected %v", tc.val, buf.Bytes(), tc.expected)
		}

		bufRead := NewBuffer(tc.expected)
		val, err := bufRead.ReadVarInt()
		if err != nil {
			t.Errorf("ReadVarInt(%v) error: %v", tc.expected, err)
		}
		if val != tc.val {
			t.Errorf("ReadVarInt(%v) = %d, expected %d", tc.expected, val, tc.val)
		}
	}
}

func TestBuffer_VarLong(t *testing.T) {
	tests := []struct {
		val      int64
		expected []byte
	}{
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{127, []byte{0x7F}},
		{128, []byte{0x80, 0x01}},
		{2147483647, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x07}},
		{9223372036854775807, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F}},
		{-1, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01}},
		{-9223372036854775808, []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01}},
	}

	for _, tc := range tests {
		buf := NewBuffer(nil)
		buf.WriteVarLong(tc.val)
		if !bytes.Equal(buf.Bytes(), tc.expected) {
			t.Errorf("WriteVarLong(%d) = %v, expected %v", tc.val, buf.Bytes(), tc.expected)
		}

		bufRead := NewBuffer(tc.expected)
		val, err := bufRead.ReadVarLong()
		if err != nil {
			t.Errorf("ReadVarLong(%v) error: %v", tc.expected, err)
		}
		if val != tc.val {
			t.Errorf("ReadVarLong(%v) = %d, expected %d", tc.expected, val, tc.val)
		}
	}
}

func TestBuffer_String(t *testing.T) {
	s := "Hello Minecraft!"
	buf := NewBuffer(nil)
	buf.WriteString(s)

	readBuf := NewBuffer(buf.Bytes())
	readS, err := readBuf.ReadString()
	if err != nil {
		t.Fatalf("ReadString error: %v", err)
	}
	if readS != s {
		t.Errorf("ReadString() = %q, expected %q", readS, s)
	}
}

func TestBuffer_BasicTypes(t *testing.T) {
	buf := NewBuffer(nil)
	buf.WriteBoolean(true)
	buf.WriteByte(-123)
	buf.WriteUByte(250)
	buf.WriteShort(-30000)
	buf.WriteUShort(60000)
	buf.WriteInt(-1000000)
	buf.WriteLong(-999999999999)
	buf.WriteFloat(1.234)
	buf.WriteDouble(5.678)
	id := uuid.New()
	buf.WriteUUID(id)

	readBuf := NewBuffer(buf.Bytes())
	if b, _ := readBuf.ReadBoolean(); b != true {
		t.Error("ReadBoolean failed")
	}
	if v, _ := readBuf.ReadByte(); v != -123 {
		t.Error("ReadByte failed")
	}
	if v, _ := readBuf.ReadUByte(); v != 250 {
		t.Error("ReadUByte failed")
	}
	if v, _ := readBuf.ReadShort(); v != -30000 {
		t.Error("ReadShort failed")
	}
	if v, _ := readBuf.ReadUShort(); v != 60000 {
		t.Error("ReadUShort failed")
	}
	if v, _ := readBuf.ReadInt(); v != -1000000 {
		t.Error("ReadInt failed")
	}
	if v, _ := readBuf.ReadLong(); v != -999999999999 {
		t.Error("ReadLong failed")
	}
	if v, _ := readBuf.ReadFloat(); v != 1.234 {
		t.Error("ReadFloat failed")
	}
	if v, _ := readBuf.ReadDouble(); v != 5.678 {
		t.Error("ReadDouble failed")
	}
	if v, _ := readBuf.ReadUUID(); v != id {
		t.Error("ReadUUID failed")
	}
}
