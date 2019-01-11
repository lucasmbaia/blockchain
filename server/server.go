package server

import (
	"github.com/lucasmbaia/blockchain"
	"github.com/lucasmbaia/blockchain/utils"
	"net"
	"log"
	"bufio"
	"io"
	"encoding/gob"
)

const (
	DEFAULT_PORT_SERVER = ":5688"
)

var (
	index int64
)

type gossip struct {
	option	string
	body	[]byte
}

type Client struct {
	operation   blockchain.Operation
	connection  []*connection
	address	    []byte

	block	    chan blockchain.Block
	transaction chan *blockchain.Transaction
}

type connection struct {
	write	*bufio.Writer
	read	*bufio.Reader
	conn	*net.Conn
}

func (g *Gossip) Serialize() ([]byte, error) {
	var result bytes.Buffer
	var encoder *gob.Encoder = gob.NewEncoder(&result)

	if err := encoder.Encode(g); err != nil {
		return []byte{}, err
	}

	return result.Bytes(), nil
}

func DeserializeGossip(b []byte) (*Gossip, error) {
	var gossip Gossip
	var decoder *gob.Decoder = gob.NewDecoder(bytes.NewReader(b))

	if err := decoder.Decode(&gossip); err != nil {
		return &gossip, err
	}

	return &gossip, nil
}

func StartNode(ctx context.Context, address []byte, nodes []string) {
	var client = &Client{
		address:      address,
		block:	      make(chan blockchain.Block, 1),
		transaction:  make(chan blockchain.Transaction, 1),
	}

	for _, node := range nodes {
		var c = connection(node)

		client.connection = append(client.connection, c)
	}

	go client.handleConnection()
	go client.operations()

	<-ctx.Done()
}

func (c *Client) operations() {
	for {
		select {
		case <-c.block:
		case <-c.transaction:
		}
	}
}

func (c *Client) mining(hash utils.Hash) {
	var (
		transactions  []*blockchain.Transaction
		ctbx	      *Transaction
		block	      *blockchain.Block
	)

	go func() {
		select {
		case <-c.operation.Quit():
			return
		}
	}()

	ctbx = blockchain.NewCoinbase(c.address, "Coinbase Transaction", index)
	transactions = append(transactions, ctbx)

	block = blockchain.NewBlock(c.operation, transactions, []byte(""), hash)
	c.block <- block
	return
}

func (c *Client) handleConnection() {
	var (
		ctx	  context.Context
		cancel	  context.CancelFunc
	)

	c.operation = blockchain.Operation{
		Resume:	make(chan struct{}),
		Pause:	make(chan struct{}),
	}

	ctx, cancel = context.WithCancel(context.Background())
	c.Quit = ctx
	go c.mining()

	for _, cli := range c.connection {
		defer func() {
			if addr, ok := cli.conn.RemoteAddr().(*net.TCPAddr); ok {
				log.Printf("Connection Close if node IP %s\n", addr.IP.String())
			}
			cli.conn.Close()
		}()

		var (
			err	error
			option	[]byte
		)

		for {
			if option, err = cli.read.ReadBytes('\n'); err != nil {
				if err != io.EOF {
					log.Printf("Error to read bytes: %s\n", err.Error())
				}
				return
			}

			var gossip *Gossip

			if gossip, err = DeserializeGossip(option); err != nil {
				log.Printf("Error to deserialize gossip: %s\n", err.Error())
				continue
			}

			switch gossip.action {
			case "block":
				/*c.operation.Pause <- struct{}{}

				var block = blockchain.Deserialize(gossip.body)*/
				cancel()

				ctx, cancel = context.WithCancel(context.Background())
				c.Quit = ctx
				go c.mining(block.Hash)
			default:
				log.Printf("Invalid Option")
			}
		}
	}
}

func connection(addr string) *Client {
	var (
		client	*Client
		err	error
		conn	net.Conn
	)

	if conn, err = net.Dial("tcp", addr); err != nil {
		panic(err)
	}
	defer conn.Close()

	client = &Client{
		write:	bufio.NewWriter(conn),
		read:	bufio.NewReader(conn),
		conn:	&conn,
	}

	return client
}
