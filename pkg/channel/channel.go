package channel

import (
	"bufio"
	"encoding/binary"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
)

// Channel is a unix socket server wrapper object. It provides a buffered
// abstraction on top of a socket connection.
type Channel struct {
	SocketPath string

	reader    *bufio.Reader
	writer    *bufio.Writer
	listener  net.Listener
	connected sync.WaitGroup

	wmu sync.Mutex
}

// NewChannel builds a Channel object
func NewChannel() (*Channel, error) {
	tmpFile, err := ioutil.TempFile("", "bruco-channel-socket-")
	if err != nil {
		log.Fatal("tmp file error ", err)
		return nil, err
	}
	socketPath := tmpFile.Name()
	tmpFile.Close()
	os.Remove(tmpFile.Name())

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal("dial error ", err)
		return nil, err
	}

	c := &Channel{
		listener:   listener,
		SocketPath: socketPath,
	}
	c.connected.Add(1)
	go c.listen()

	return c, nil
}

func (c *Channel) listen() {
	conn, err := c.listener.Accept()
	if err != nil {
		log.Fatalln(err)
		return
	}
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)
	c.connected.Done()
}

// Write writes data to the channel
func (c *Channel) Write(data []byte) error {
	c.connected.Wait()
	// prevents concurrent writes (short write error)
	c.wmu.Lock()
	defer c.wmu.Unlock()

	// send the msg len as the first 4 bytes
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(len(data)))
	_, err := c.writer.Write(bs)
	if err != nil {
		log.Printf("channel write error: %s", err)
		return err
	}

	l, err := c.writer.Write(data)
	if err != nil {
		log.Printf("channel write error: %s", err)
		return err
	}
	if l != len(data) {
		log.Printf("channel write error: %s", err)
		return err
	}

	c.writer.Flush()
	return nil
}

// Read reads data from the channel
func (c *Channel) Read() ([]byte, error) {
	c.connected.Wait()
	out, err := c.reader.ReadBytes('\n')
	if err != nil {
		// log.Printf("channel read error: %s", err)
		return nil, err
	}
	return out, err
}

// Close close the channel
func (c *Channel) Close() {
	c.listener.Close()
}
