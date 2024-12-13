package btree

import "encoding/binary"

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
