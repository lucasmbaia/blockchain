package blockchain

import (
	"github.com/lucasmbaia/blockchain/utils"
	"net"
	"log"
	"bufio"
	"io"
	"encoding/json"
	//"encoding/gob"
	//"bytes"
	"math/big"
	"fmt"
	"sync"
	"context"
	"errors"
)

const (
	DEFAULT_PORT_SERVER = ":5688"
	MULTIPLIER	    = 100000000
)

var (
	p2pNodes  []net.Conn
)

type server struct {
	sync.RWMutex

	operation   *operation
	connections map[string]*connection
	wcAddress   []byte
}

type operation struct {
	done	chan struct{}
	resume	chan struct{}
	pause	chan struct{}
}

type connection struct {
	write *bufio.Writer
	read  *bufio.Reader
	conn  net.Conn
}

type gossip struct {
	Option  string
	Body    []byte
	Error	error
}

type Infos struct {
	Private	string
	From	string
	To	string
	Value	float64
}

func (i *Infos) Serialize() ([]byte, error) {
	return json.Marshal(i)
	/*var result bytes.Buffer
	var encoder *gob.Encoder = gob.NewEncoder(&result)

	if err := encoder.Encode(i); err != nil {
		log.Printf("Error to serialize infos: %s\n", err)
	}

	return result.Bytes()*/
}

func DeserializeInfos(b []byte) (*Infos, error) {
	var (
		infos Infos
		err   error
	)

	err = json.Unmarshal(b, &infos)
	return &infos, err
	/*var decoder *gob.Decoder = gob.NewDecoder(bytes.NewReader(b))

	if err := decoder.Decode(infos); err != nil {
		log.Printf("Error to deserialize: %s\n", err)
	}

	return &infos*/
}

func (s *server) removeNode(node string) {
	if node != "" {
		s.Lock()
		if _, ok := s.connections[node]; ok {
			delete(s.connections, node)
		}
		s.Unlock()
	}
}

func (s *server) replyToNodes(g gossip, exception string) error {
	var (
		body  []byte
		err   error
	)

	if body, err = encodeGossip(g); err != nil {
		return err
	}

	s.Lock()
	for node, connection := range s.connections {
		if node != exception {
			if _, err = connection.conn.Write(body); err != nil {
				return err
			}
		}
	}
	s.Unlock()

	return nil
}

func (s *server) handleConnection(c *connection, bc *Blockchain) {
	defer func() {
		if addr, ok := c.conn.RemoteAddr().(*net.TCPAddr); ok {
			log.Printf("Connection Close if node IP %s", addr.IP.String())
			s.removeNode(addr.IP.String())
		}
		c.conn.Close()
	}()

	var (
		err	error
		option	[]byte
	)

	for {
		if option, err = c.read.ReadBytes('\n'); err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}

		var g gossip

		err = json.Unmarshal(option[:len(option)-1], &g)

		switch g.Option {
		case "sizeof_chain":
			var (
				block	*Block
				body	[]byte
				encoded	[]byte
			)

			if block, err = getHead(bc); err != nil {
				break
			}

			history := struct{
				Index int32
				Hash  utils.Hash
			}{block.Index, block.Hash}

			if body, err = json.Marshal(history); err != nil {
				break
			}

			if encoded, err = encodeGossip(gossip{Option: "sizeof_chain", Body: body}); err != nil {
				break
			}

			c.conn.Write(encoded)
		case "history_blocks":
			var blocks = getBlocks(bc)

			for _, b := range blocks {
				body, _ := json.Marshal(b)
				encoded, _ := encodeGossip(gossip{Option: "history_blocks", Body: body})
				c.conn.Write(encoded)
			}
		case "transaction":
		case "local_transaction":
			var infos *Infos

			if infos, err = DeserializeInfos(g.Body); err != nil {
				fmt.Println(err)
			}

			if err = s.validTransaction(infos, bc); err != nil {
				fmt.Println(err)
			}

			body, _ := encodeGossip(gossip{Option: "local_transaction"})
			c.conn.Write(body)
			fmt.Println("TRANSACTION")
		}
	}
}

func (s *server) validTransaction(i *Infos, bc *Blockchain) error {
	var (
		err     error
		amount  uint64
		w       *Wallet
		valid   bool
		tx      *Transaction
		txs     []Transaction
	)

	amount = uint64(i.Value * MULTIPLIER)

	if valid, w, err = UnlockWallet(i.Private, i.From); err != nil {
		return errors.New(fmt.Sprintf("Error to unlock waller: %s", err.Error()))
	}

	if !valid {
		return errors.New(fmt.Sprintf("The private key is not allowed to unlock the wallet"))
	}

	if tx, err = bc.NewTransaction([]byte(i.From), []byte(i.To), amount); err != nil {
		return errors.New(fmt.Sprintf("Error to generate transaction: %s", err.Error()))
	}

	for _, input := range tx.TXInput {
		_, transaction, _ := bc.FindTransaction(input.TXid)

		txs = append(txs, transaction)
	}

	if err = tx.SignTransaction(w, txs); err != nil {
		return errors.New(fmt.Sprintf("Error to sign transaction: %s", err.Error()))
	}

	return nil
}

func getBlocks(bc *Blockchain) ([]*Block) {
	var (
		blocks	[]*Block
		bci	= bc.Iterator()
		stop	= big.NewInt(0)
	)

	for {
		block, _ := bci.Next()
		blocks = append(blocks, block)

		if HashToBig(&block.Header.PrevBlock).Cmp(stop) == 0 {
			break
		}
	}

	return blocks
}

func getHead(bc *Blockchain) (*Block, error) {
	var bci = bc.Iterator()

	return bci.Next()
}

func StartFullNode(ctx context.Context, address []byte, nodes []string) {
	var (
		err error
		l   net.Listener
		bc  = NewBlockchain(address)
		s   *server
	)

	if l, err = net.Listen("tcp", DEFAULT_PORT_SERVER); err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	s = &server{
		wcAddress:    address,
		connections:  make(map[string]*connection),
		operation:    &operation{
			done:	make(chan struct{}),
			resume: make(chan struct{}),
			pause:	make(chan struct{}),
		},
	}

	go func() {
		for {
			var (
				client	= &connection{}
			)

			if client.conn, err = l.Accept(); err != nil {
				log.Fatalf(fmt.Sprintf("Error to open port: %s\n", err.Error()))
			}

			client.write = bufio.NewWriter(client.conn)
			client.read = bufio.NewReader(client.conn)

			s.Lock()
			s.connections[client.conn.RemoteAddr().(*net.TCPAddr).IP.String()] = client
			s.Unlock()

			go s.handleConnection(client, bc)
		}
	}()

	<-ctx.Done()
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
