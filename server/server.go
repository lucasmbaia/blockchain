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
	"sync"
)

const (
	DEFAULT_PORT_SERVER = ":5689"
)

var (
	mutex	      = &sync.RWMutex{}
	index	      int32
	history	      = make(chan History)
	blockHistory  = make(chan *blockchain.Block)
	blocks	      []*blockchain.Block
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
}

type Operation struct {
	Done	chan struct{}
	Resume	chan struct{}
	Pause	chan struct{}
}

type History struct {
	Index int32
	Hash  utils.Hash
	Node  string
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
		operation:    &Operation{
			Done:   make(chan struct{}),
			Resume: make(chan struct{}),
			Pause:  make(chan struct{}),
		},
	}

	go client.listen()

	for _, node := range nodes {
		var c = dial(node)

		client.connection = append(client.connection, c)
	}

	go client.handleConnection()
	go client.operations()
	go client.getHistory()

	<-ctx.Done()
}

func (c *Client) getHistory() {
	var (
		wg	  sync.WaitGroup
		hash	  utils.Hash
		done	  = make(chan struct{})
		body	  []byte
		err	  error
		nodes	  = make(map[string]net.Conn)
		master  string
		b	  []*blockchain.Block
	)

	if body, err = encodeGossip(gossip{Option: "sizeof_chain"}); err != nil {
		panic(err)
	}

	wg.Add(len(c.connection))
	go func() {
		go func() {
			for {
				select {
				case h := <-history:
					if index < h.Index {
						index = h.Index
						hash = h.Hash
						master = h.Node
					}
					wg.Done()
				case <-done:
					return
				}
			}
		}()

		for _, conn := range c.connection {
			conn.conn.Write(body)
			nodes[conn.conn.RemoteAddr().(*net.TCPAddr).IP.String()] = conn.conn
		}
	}()
	wg.Wait()
	done <- struct{}{}

	wg.Add(int(index) + 1)
	go func() {
		go func() {
			for {
				select {
				case block := <-blockHistory:
					b = append(b, block)
					wg.Done()
				case <-done:
					return
				}
			}
		}()

		if body, err = encodeGossip(gossip{Option: "history_blocks"}); err != nil {
			panic(err)
		}

		nodes[master].Write(body)
	}()
	wg.Wait()
	done <- struct{}{}

	for i := len(b) - 1; i >= 0; i-- {
		blocks = append(blocks, b[i])
	}

	go c.mining(hash)
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
		case transaction := <-c.transaction:
			if len(c.connection) > 0 {
				var body []byte

				if body, err = encodeGossip(gossip{
					Option: "transaction",
					Body:	transaction.Serialize(),
				}); err == nil {
					for _, conn := range c.connection{
						conn.conn.Write(body)
					}
				} else {
					fmt.Printf("Error to send transaction gossip: %s\n", err.Error())
				}
			}
		}
	}
}

func (c *Client) mining(hash utils.Hash) {
	fmt.Println("MINERANDO")
	var (
		transactions	[]*blockchain.Transaction
		ctbx		*blockchain.Transaction
		block		*blockchain.Block
		operations	blockchain.Operations
		ctx		context.Context
		cancel		context.CancelFunc
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

	for _, tx := range blockchain.UnprocessedTransactions {
		if err := tx.ValidTransaction(getTransactionsInBlocks(tx)); err == nil {
			transactions = append(transactions, tx)
		}
	}
	blocks[len(blocks) - 1].CheckProcessedTransactions(transactions)

	block = blockchain.NewBlock(operations, index, transactions, []byte(""), hash)
	c.block <- block
	return
}

func (c *Client) handleConnection() {
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

			if err = json.Unmarshal(option[:len(option)-1], &gossip); err != nil {
				log.Printf("Error to deserialize gossip: %s\n", err.Error())
				continue
			}

			switch gossip.Option {
			case "block":
				c.operation.Pause <- struct{}{}

				var block = blockchain.Deserialize(gossip.Body)

				c.operation.Done <- struct{}{}
				go c.mining(block.Hash)
			case "sizeof_chain":
				var h History

				if err = json.Unmarshal(gossip.Body, &h); err == nil {
					h.Node = cli.conn.RemoteAddr().(*net.TCPAddr).IP.String()
					history <- h
				} else {
					log.Printf("Error to deserialize history: %s\n", err.Error())
				}
			case "history_blocks":
				var b blockchain.Block

				if err = json.Unmarshal(gossip.Body, &b); err == nil {
					blockHistory <- &b
				} else {
					log.Printf("Error to deserialize block: %s\n", err.Error())
				}
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

func encodeGossip(g gossip) ([]byte, error) {
	var (
		body  []byte
		err   error
	)

	if body, err = json.Marshal(g); err != nil {
		return body, err
	}

	body = append(body, byte('\n'))
	return body, nil
}

func getTransactionsInBlocks(tx *blockchain.Transaction) map[utils.Hash]*blockchain.Transaction {
	var transactions = make(map[utils.Hash]*blockchain.Transaction)

	for _, b := range blocks {
		for _, tx := range b.Transactions {
			for _, input := range tx.TXInput {
				if bytes.Compare(input.TXid[:], tx.ID[:]) == 0 {
					transactions[tx.ID] = tx
				}
			}
		}
	}

	return transactions
}
