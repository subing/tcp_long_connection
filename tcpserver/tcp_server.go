package tcpserver

import (
	"bufio"
	"net"
	"tcp_long_connection/logger"
)

// Client holds info about connection
type Client struct {
	conn     net.Conn
	Server   *server
	incoming chan string // Channel for incoming data from client
}

// TCP server
type server struct {
	clients                  []*Client
	address                  string        // Address to open connection: localhost:9999
	joins                    chan net.Conn // Channel for new connections
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, p Package, message string)
}

// Read client data from channel string
func (c *Client) listen() {
	var p Package
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}
		c.Server.onNewMessage(c, p, message)
	}
}

//Read client data from channel byte
func (c *Client) listenByte() {
	reader := bufio.NewReader(c.conn)
	for {
		p, message, err := Decode(reader)
		if err != nil {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}
		c.Server.onNewMessage(c, p, message)
	}
}
func (c *Client) Send(p Package, message string) error {
	sndBuf, err := p.Encode(message)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	_, err = c.conn.Write(sndBuf)
	return err
}

func (c *Client) GetClientAddr() string {
	return c.conn.RemoteAddr().String()
}

// Called right after server starts listening new client
func (s *server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *server) OnNewMessage(callback func(c *Client, p Package, message string)) {
	s.onNewMessage = callback
}

// Creates new Client instance and starts listening
func (s *server) newClient(conn net.Conn) {
	client := &Client{
		conn:   conn,
		Server: s,
	}
	//go client.listen()
	go client.listenByte()
	s.onNewClientCallback(client)
}

// Listens new connections channel and creating new client
func (s *server) listenChannels() {
	for {
		select {
		case conn := <-s.joins:
			s.newClient(conn)
		}
	}
}

// Start network server
func (s *server) Listen() {
	go s.listenChannels()

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		logger.Error("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		s.joins <- conn
	}
}

// Creates new tcp server instance
func New(address string) *server {
	logger.Info("Creating server with address", address)
	server := &server{
		address: address,
		joins:   make(chan net.Conn),
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, p Package, message string) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}
