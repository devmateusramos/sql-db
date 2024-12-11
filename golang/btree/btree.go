package btree

const (
	HEADER             = 4
	BTREE_PAGE_SIZE    = 4096
	BTREE_MAX_KEY_SIZE = 1000
	BTREE_MAX_VAL_SIZE = 3000
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
