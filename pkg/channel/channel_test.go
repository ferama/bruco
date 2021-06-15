package channel

import (
	"net"
	"testing"
)

func TestWrite(t *testing.T) {
	c, err := NewChannel()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	conn, err := net.Dial("unix", c.SocketPath)
	if err != nil {
		t.Fatal(err)
	}
	testData := []byte("test\n")

	conn.Write(testData)
	data, err := c.Read()
	if err != nil {
		t.Fatal(err)
	}
	if string(testData) != string(data) {
		t.Fatal("data read is not equal to data write")
	}
}
