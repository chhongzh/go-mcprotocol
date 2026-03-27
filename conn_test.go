package mcprotocol

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestConn_ReadWriteNext(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	conn1 := NewConn(c1)
	conn2 := NewConn(c2)

	testData := []byte{0x00, 0x01, 0x02, 0x03}
	buf := NewBuffer(testData)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Write in a goroutine
	go func() {
		err := conn1.WriteNext(ctx, 0x42, buf)
		if err != nil {
			t.Errorf("WriteNext error: %v", err)
		}
	}()

	// Read from conn2
	packetID, readBuf, err := conn2.ReadNext(ctx)
	if err != nil {
		t.Fatalf("ReadNext error: %v", err)
	}

	if packetID != 0x42 {
		t.Errorf("Packet ID = %x, expected %x", packetID, 0x42)
	}

	if len(readBuf.Data) != len(testData) {
		t.Fatalf("Read data length = %d, expected %d", len(readBuf.Data), len(testData))
	}

	for i := range testData {
		if readBuf.Data[i] != testData[i] {
			t.Errorf("Data at %d = %x, expected %x", i, readBuf.Data[i], testData[i])
		}
	}
}

func TestConn_ContextCancel(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	conn2 := NewConn(c2)

	ctx, cancel := context.WithCancel(context.Background())

	// Start reading
	errCh := make(chan error, 1)
	go func() {
		_, _, err := conn2.ReadNext(ctx)
		errCh <- err
	}()

	// Cancel context
	time.Sleep(100 * time.Millisecond)
	cancel()

	err := <-errCh
	if err == nil {
		t.Error("ReadNext should have returned an error on context cancellation")
	}
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}
