package processor

import (
	"bufio"
	"log"
	"net"
	"sync"
)

type channel struct {
	reader    *bufio.Reader
	writer    *bufio.Writer
	Port      int
	listener  net.Listener
	connected sync.WaitGroup

	rmu sync.Mutex
}

func newChannel() (*channel, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("dial error ", err)
		return nil, err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	// log.Printf("listening on port %d", port)

	c := &channel{
		listener: listener,
		Port:     port,
	}
	c.connected.Add(1)
	go c.listen()

	return c, nil
}

func (c *channel) listen() {
	conn, err := c.listener.Accept()
	if err != nil {
		log.Fatalln(err)
		return
	}
	// log.Printf("client connected: %s", conn.RemoteAddr())
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)
	c.connected.Done()
}

func (c *channel) Write(data []byte) error {
	c.connected.Wait()
	// prevents concurrent writes (short write error)
	c.rmu.Lock()
	defer c.rmu.Unlock()

	_, err := c.writer.Write(data)
	if err != nil {
		log.Printf("channel write error: %s", err)
		return err
	}
	c.writer.Flush()
	return nil
}

func (c *channel) Read() ([]byte, error) {
	c.connected.Wait()
	out, err := c.reader.ReadBytes('\n')
	if err != nil {
		log.Printf("channel read error: %s", err)
		return nil, err
	}
	return out, err
}

func (c *channel) Close() {
	c.listener.Close()
}
