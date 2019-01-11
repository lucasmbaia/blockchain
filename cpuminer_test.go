package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/lucasmbaia/blockchain/utils"
	"runtime"
	"testing"
	"time"
	"context"
)

const (
	BITS_TEST_MINIG = 24
)

func Test_Cpu_Miner(t *testing.T) {
	var (
		prevH   utils.Hash
		merkleH utils.Hash
		bh      BlockHeader
		td      time.Time
		prev    []byte
		merkle  []byte
		err     error
	)

	runtime.GOMAXPROCS(runtime.NumCPU())

	td = time.Date(2018, time.March, 31, 16, 37, 31, 0, time.UTC)
	if prev, err = hex.DecodeString(PREV_HASH_LT_BLOCK_125552); err != nil {
		t.Fatal(err)
	}

	if merkle, err = hex.DecodeString(MERKLE_HASH_LT_BLOCK_125552); err != nil {
		t.Fatal(err)
	}

	copy(prevH[:], prev)
	copy(merkleH[:], merkle)

	bh = BlockHeader{
		Version:    int32(VERSION_BLOCK_125552),
		//PrevBlock:  prevH,
		//MerkleRoot: merkleH,
		Bits:	    uint32(BITS_TEST_MINIG),
		Timestamp:  td,
	}

	fmt.Println(CpuMiner(&bh))
}

func Test_Cpu_Control_Miner(t *testing.T) {
	var (
		prevH	    utils.Hash
		merkleH	    utils.Hash
		bh	    BlockHeader
		td	    time.Time
		prev	    []byte
		merkle	    []byte
		err	    error
		operations  Operations
		done	    = make(chan struct{})
	)

	runtime.GOMAXPROCS(runtime.NumCPU())

	td = time.Date(2018, time.March, 31, 16, 37, 31, 0, time.UTC)
	if prev, err = hex.DecodeString(PREV_HASH_LT_BLOCK_125552); err != nil {
		t.Fatal(err)
	}

	if merkle, err = hex.DecodeString(MERKLE_HASH_LT_BLOCK_125552); err != nil {
		t.Fatal(err)
	}

	copy(prevH[:], prev)
	copy(merkleH[:], merkle)

	bh = BlockHeader{
		Version: int32(VERSION_BLOCK_125552),
		//PrevBlock:  prevH,
		//MerkleRoot: merkleH,
		Bits:      uint32(BITS_TEST_MINIG),
		Timestamp: td,
	}

	ctx, _ := context.WithCancel(context.Background())

	operations = Operations{
		Quit:	ctx,
		Resume:	make(chan struct{}),
		Pause:	make(chan struct{}),
	}

	go func() {
		CpuMinerControl(operations, &bh)
		done <- struct{}{}
	}()

	time.Sleep(2 * time.Second)
	operations.Pause <- struct{}{}

	time.Sleep(2 * time.Second)
	operations.Resume <- struct{}{}

	<-done
}
