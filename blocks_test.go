package blocks

import (
	"bytes"
	"testing"

	u "github.com/ipfs/boxo/util"
	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

func TestBlocksBasic(t *testing.T) {
	// Test empty data
	empty := []byte{}
	NewBlock(empty)

	// Test nil case
	NewBlock(nil)

	// Test some data
	NewBlock([]byte("Hello world!"))
}

func TestData(t *testing.T) {
	data := []byte("some data")
	block := NewBlock(data)

	if !bytes.Equal(block.RawData(), data) {
		t.Error("data is wrong")
	}
}

func TestHash(t *testing.T) {
	data := []byte("some other data")
	block := NewBlock(data)

	hash, err := mh.Sum(data, mh.SHA2_256, -1)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(block.Multihash(), hash) {
		t.Error("wrong multihash")
	}
}

func TestCid(t *testing.T) {
	data := []byte("yet another data")
	block := NewBlock(data)
	c := block.Cid()

	if !bytes.Equal(block.Multihash(), c.Hash()) {
		t.Error("key contains wrong data")
	}
}

func TestManualHash(t *testing.T) {
	oldDebugState := u.Debug
	defer (func() {
		u.Debug = oldDebugState
	})()

	data := []byte("I can't figure out more names .. data")
	hash, err := mh.Sum(data, mh.SHA2_256, -1)
	if err != nil {
		t.Fatal(err)
	}

	c := cid.NewCidV0(hash)

	u.Debug = false
	block, err := NewBlockWithCid(data, c)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(block.Multihash(), hash) {
		t.Error("wrong multihash")
	}

	data[5] = byte((uint32(data[5]) + 5) % 256) // Transfrom hash to be different
	block, err = NewBlockWithCid(data, c)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(block.Multihash(), hash) {
		t.Error("wrong multihash")
	}

	u.Debug = true

	_, err = NewBlockWithCid(data, c)
	if err != ErrWrongHash {
		t.Fatal(err)
	}
}

func TestNewBlockWithPrefix(t *testing.T) {
	data := []byte("yet another test")
	hash, err := mh.Sum(data, mh.SHA2_256, -1)
	if err != nil {
		t.Fatal(err)
	}

	c := cid.NewCidV0(hash)
	prefix := c.Prefix()

	block, err := NewBlockWithPrefix(data, prefix)
	if err != nil {
		t.Fatal(err)
	}

	if !block.Cid().Equals(c) {
		t.Fatal("block CID is not equal to CID")
	}
	if !bytes.Equal(block.RawData(), data) {
		t.Fatal("block has wrong data")
	}

	prefix2 := block.Cid().Prefix()
	if prefix2 != prefix {
		t.Fatal("prefixes are not equal")
	}

	// Test with bad prefix.
	prefix = cid.Prefix{}
	_, err = NewBlockWithPrefix(data, prefix)
	if err == nil {
		t.Fatal("expected error with bad prefix")
	}
}
