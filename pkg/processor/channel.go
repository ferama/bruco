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
	go c.listen(true)

	return c, nil
}

func (c *channel) listen(onNewChannel bool) {
	if !onNewChannel {
		c.connected.Add(1)
	}
	conn, err := c.listener.Accept()
	if err != nil {
		// log.Println(err)
		return
	}
	// log.Printf("client connected: %s", conn.RemoteAddr())
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)
	c.connected.Done()
}

func (c *channel) Write(data []byte) error {
	c.connected.Wait()
	_, err := c.writer.Write(data)
	if err != nil {
		// log.Printf("write err: %s", err)
		c.listen(false)
		return err
	}
	c.writer.Flush()
	return nil
}

func (c *channel) Read() ([]byte, error) {
	c.connected.Wait()
	out, err := c.reader.ReadBytes('\n')
	if err != nil {
		// log.Printf("read err: %s", err)
		c.listen(false)
		return nil, err
	}
	return out, err
}

func (c *channel) Close() {
	// log.Println("ch closed")
	c.listener.Close()
}
