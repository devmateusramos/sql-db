package btree

import (
	"bytes"
	"encoding/binary"
)

const (
	HEADER             = 4
	BTREE_PAGE_SIZE    = 4096
	BTREE_MAX_KEY_SIZE = 1000
	BTREE_MAX_VAL_SIZE = 3000
	// Decode the node format
	BNODE_NODE = 1 // internal nodes without values
	BNODE_LEAF = 2 // leaf nodes with values
)

func assert(cond bool) {
	if !cond {
		panic("invariant broken")
	}
}

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	assert(node1max <= BTREE_PAGE_SIZE)
}

type BNode []byte // can be dumped to the disk

/*
 B+Tree the database file is an array of pages(nodes) referenced
by page numbers (pointers)

* get reads a page from disk
* new allocates and writes a new page (copy-on-write)
* del deallocates a page

*/

type BTree struct {
	// pointer (a nonzero page number)
	root uint64
	// callbacks fro manging on-disk pages
	get func(uint64) []byte // dereference a pointer
	new func([]byte) uint64 // allocate a new page
	del func(uint64)        // deallocate a page
}

func (node BNode) bType() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

func (node BNode) nKeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

func (node BNode) setHeader(bType uint16, nKeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], bType)
	binary.LittleEndian.PutUint16(node[2:4], nKeys)
}

// pointers
func (node BNode) getPtr(idx uint16) uint64 {
	assert(idx < node.nKeys())
	pos := HEADER + (8 * idx)
	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) setPtr(idx uint16, val uint64) {
}

func offsetPos(node BNode, idx uint16) uint16 {
	assert(1 <= idx && idx <= node.nKeys())
	return HEADER + (8 * node.nKeys()) + (2 * (idx - 1))
}

func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}

	return binary.LittleEndian.Uint16(node[offsetPos(node, idx):])
}

func (node BNode) setOffSet(idx uint16, offset uint16) {}

// key-values
func (node BNode) kvPos(idx uint16) uint16 {
	assert(idx <= node.nKeys())
	return HEADER + (8 * node.nKeys()) + (2 * node.nKeys()) + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	assert(idx < node.nKeys())
	pos := node.kvPos(idx)
	kLen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:kLen]
}

// func (node BNode) getVal(idx uint16) []byte

// node size in bytes
func (node BNode) nBytes() uint16 {
	return node.kvPos(node.nKeys())
}

// Return the first kid node whose range intersects the key. (kid[i] <= key)
// TODO binary search
func nodeLookupLE(node BNode, key []byte) uint16 {
	nKeys := node.nKeys()
	found := uint16(0)

	// the first key is a copy from the parent node,
	// thus it's always less than equal to the key
	for i := uint16(1); i < nKeys; i++ {
		cmp := bytes.Compare(node.getKey(i), key)
		if cmp <= 0 {
			found = i
		}
		if cmp >= 0 {
			break
		}
	}
	return found
}

// add a new key to a leaf node
func leafInsert(
  new, old BNode, idx uint16,
  key, val []byte,
) {
  new.setHeader(BNODE_LEAF,old.nKeys()+1) // setup the header 
  nodeAppendRange(new, old, 0, 0, idx)
  nodeAppendKV(new, idx, 0, key, val)
  nodeAppendRange()
}

func nodeAppendKV(new BNode, idx uint16, ptr uint64, key, val []byte) {
  // ptrs
  new.setPtr(idx, ptr)
  // KVs
  pos := new.kvPos(idx)
  binary.LittleEndian.PutUint16(new[pos+0:], uint16(len(key)))
  binary.LittleEndian.PutUint16(new[pos+2:], uint16(len(val)))
  copy(new[pos+4:], key)
  copy(new[pos+4+uint16(len(key)):], val)
  // the offset of the next key
  new.setOffSet(idx + 1, new.getOffset(idx) + 4 + uint16(len(key) + len(val))
}

// copy multiple KVs into the position from the old node
func nodeAppendRange(new, old BnBNode, dstNew, srcOld,n uint16) {

}
