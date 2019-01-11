package blockchain

import (
    "net"
    "log"
    "bufio"
    "fmt"
    "io"
    //"strings"
)

const (
    DEFAULT_PORT_SERVER = ":5688"
)

var (
    p2pNodes  []net.Conn
)

type Client struct {
  write	*bufio.Writer
  read	*bufio.Reader
  bc	*Blockchain
}

func handleConnection(conn net.Conn) {
    defer func() {
	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
	    log.Printf("Connection Close if node IP %s", addr.IP.String())
	}
	conn.Close()
    }()

    var (
	err	error
	option	string
	b	*bufio.Reader
    )

    b = bufio.NewReader(conn)

    for {
	if option, err = b.ReadString('\n'); err != nil {
	    if err != io.EOF {
	      log.Println(err)
	    }
	    return
	}

	switch option {
	case "new_transaction":
	    fmt.Println("PORRA")
	case "valid_block":
	default:
	    fmt.Println(option)
	}
    }
}

func (c *Client) HandleNewTrancation(data []byte) {
  var (
    transaction	*Transaction
    err		error
  )

  transaction = DeserializeTransaction(data)

  if err = ValidTransaction(transaction, c.bc); err != nil {
    fmt.Println(err)
    return
  }

  AppendUnprocessedTransactions(transaction)
}

func StartFullNode(nodes []string) {
    var (
	err	error
	l	net.Listener
    )

    if l, err = net.Listen("tcp", DEFAULT_PORT_SERVER); err != nil {
	log.Fatal(err)
    }
    defer l.Close()

    for {
	var (
	  conn	  net.Conn
	  //client  Client
	)

	if conn, err = l.Accept(); err != nil {
	    log.Fatal(err)
	}

	/*client = Client{
	  write:  bufio.NewWriter(conn),
	  read:	  bufio.NewReader(conn),
	}*/

	go handleConnection(conn)
    }
}

func StartNode(nodes []string) {
    var (
	err error
	conn  net.Conn
	done = make(chan bool, 1)
    )

    if conn, err = net.Dial("tcp", "192.168.75.133:5688"); err != nil {
	log.Fatal(err)
    }
    defer conn.Close()

    conn.Write([]byte("ovo\n"))
    <-done
}
