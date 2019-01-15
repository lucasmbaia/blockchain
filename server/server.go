package server

import (
	"github.com/lucasmbaia/blockchain"
	"github.com/lucasmbaia/blockchain/utils"
	"net"
	"log"
	"bufio"
	"io"
	"encoding/gob"
	"bytes"
	"context"
	"fmt"
	"encoding/json"
)

const (
	DEFAULT_PORT_SERVER = ":5689"
)

var (
	index int32
)

type gossip struct {
	Option	string
	Body	[]byte
}

type Client struct {
	operation   *Operation
	connection  []*connection
	address	    []byte

	block	    chan *blockchain.Block
	transaction chan *blockchain.Transaction
	mining	    chan struct{}
}

type Operation struct {
	Done	chan struct{}
	Resume	chan struct{}
	Pause	chan struct{}
}

type connection struct {
	write	*bufio.Writer
	read	*bufio.Reader
	conn	net.Conn
}

func (g *gossip) Serialize() ([]byte, error) {
	var result bytes.Buffer
	var encoder *gob.Encoder = gob.NewEncoder(&result)

	if err := encoder.Encode(g); err != nil {
		return []byte{}, err
	}

	return result.Bytes(), nil
}

func DeserializeGossip(b []byte) (*gossip, error) {
	var gossip gossip
	var decoder *gob.Decoder = gob.NewDecoder(bytes.NewReader(b))

	if err := decoder.Decode(&gossip); err != nil {
		return &gossip, err
	}

	return &gossip, nil
}

func StartNode(ctx context.Context, address []byte, nodes []string) {
	var client = &Client{
		address:      address,
		block:	      make(chan *blockchain.Block, 1),
		transaction:  make(chan *blockchain.Transaction, 1),
	}

	for _, node := range nodes {
		var c = dial(node)

		client.connection = append(client.connection, c)
	}

	go client.handleConnection()
	go client.operations()

	<-ctx.Done()
}

func (c *Client) getHistory() {

}

func (c *Client) operations() {
	var (
		err error
	)

	for {
		select {
		case block := <-c.block:
			if len(c.connection) > 0 {
				var body []byte
				var g = gossip{
					Option:	"block",
					Body:	block.Serialize(),
				}

				if body, err = json.Marshal(g); err == nil {
					body = append(body, byte('\n'))
					for _, conn := range c.connection{
						conn.conn.Write(body)
					}
				} else {
					log.Printf("Error to serializer: %s\n", err.Error())
				}
			}

			fmt.Printf("%x\n", block.Hash)
			go c.mining(block.Hash)
		case <-c.transaction:
		}
	}
}

func (c *Client) mining(hash utils.Hash) {
	fmt.Println("MINERANDO")
	var (
		transactions  []*blockchain.Transaction
		ctbx	      *blockchain.Transaction
		block	      *blockchain.Block
		operations     blockchain.Operations
		ctx	      context.Context
		cancel	      context.CancelFunc
	)

	ctx, cancel = context.WithCancel(context.Background())
	operations = blockchain.Operations{
		Quit:	ctx,
		Resume:	make(chan struct{}),
		Pause:	make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-c.operation.Done:
				cancel()
				return
			case <-c.operation.Resume:
				operations.Resume <- struct{}{}
			case <-c.operation.Pause:
				operations.Resume <- struct{}{}
			}
		}
	}()

	ctbx = blockchain.NewCoinbase(c.address, "Coinbase Transaction", int64(index))
	transactions = append(transactions, ctbx)

	block = blockchain.NewBlock(operations, index, transactions, []byte(""), hash)
	c.block <- block
	return
}

func (c *Client) handleConnection() {
	var hash utils.Hash

	c.operation = &Operation{
		Done:	make(chan struct{}),
		Resume:	make(chan struct{}),
		Pause:	make(chan struct{}),
	}

	go c.mining(hash)

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

			var gossip *gossip

			if gossip, err = DeserializeGossip(option); err != nil {
				log.Printf("Error to deserialize gossip: %s\n", err.Error())
				continue
			}

			switch gossip.Option {
			case "block":
				c.operation.Pause <- struct{}{}

				var block = blockchain.Deserialize(gossip.Body)

				c.operation.Done <- struct{}{}
				go c.mining(block.Hash)
			case "search_transaction":
				//var transaction = blockchain.DeserializeTransaction(gossip.body)

			default:
				log.Printf("Invalid Option")
			}
		}
	}
}

func dial(addr string) *connection {
	var (
		client	*connection
		err	error
		conn	net.Conn
	)

	if conn, err = net.Dial("tcp", addr); err != nil {
		panic(err)
	}

	client = &connection{
		write:	bufio.NewWriter(conn),
		read:	bufio.NewReader(conn),
		conn:	conn,
	}

	return client
}
