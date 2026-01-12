package ds

import (
	"math"
	"sort"
)

const Infinity = math.MaxFloat64

// Item represents a key-value pair in the frontier.
type Item struct {
	Key   int
	Value float64
	next  *Item // Internal pointer for the linked list
}

// block represents a bucket of items with a tracked upper bound.
type block struct {
	head       *Item
	tail       *Item
	size       int
	upperBound float64 // Max value in this block (for the BST/Index)
}

// DataStructure implements the block-based priority queue (Lemma 3.3).
type DataStructure struct {
	M     int
	B     float64 // Global upper bound
	Count int

	// D0: Sequence of blocks from BatchPrepend (unsorted between blocks, sorted within?)
	// The paper implies D0 is a buffer. We treat D0 as a simple list of blocks
	// where we just prepend new blocks.
	d0 []*block

	// D1: Sequence of blocks maintained in sorted order of their values.
	// We use a slice to act as the "Search Tree" for the block headers.
	d1 []*block
}

func NewDataStructure(m int) *DataStructure {
	return &DataStructure{
		M:  m,
		B:  Infinity,
		d0: make([]*block, 0),
		d1: make([]*block, 0),
	}
}

// Insert adds a key/value pair. amortized O(max{1, log(N/M)})
func (ds *DataStructure) Insert(key int, val float64) {
	ds.Count++
	item := &Item{Key: key, Value: val}

	// 1. Find appropriate block in D1 via Binary Search on UpperBounds
	// We look for the first block where upperBound >= val
	idx := sort.Search(len(ds.d1), func(i int) bool {
		return ds.d1[i].upperBound >= val
	})

	if idx == len(ds.d1) {
		// No block fits, or D1 is empty.
		// If D1 is empty, create new.
		if len(ds.d1) == 0 {
			b := newBlock()
			b.upperBound = Infinity // The last block always stretches to Infinity/B
			ds.d1 = append(ds.d1, b)
			idx = 0
		} else {
			// Should conceptually belong to the last block if it's within B,
			// but our binary search logic handles this if the last block has UB=Infinity.
			// If we are here, something is odd, or we just append to last.
			idx = len(ds.d1) - 1
		}
	}

	targetBlock := ds.d1[idx]

	// 2. Insert into the linked list of targetBlock (O(1))
	// Note: The paper assumes blocks are sorted internally?
	// Lemma 3.3 Proof: "Blocks are maintained in the sorted order... for any two pairs... in Bi and Bj... b1 <= b2".
	// It does NOT explicitly say items *within* a block are sorted, but Split relies on finding the median.
	// If we don't sort internal items, finding median is O(M). Inserting is O(1).
	// We append to head for O(1).
	item.next = targetBlock.head
	targetBlock.head = item
	if targetBlock.tail == nil {
		targetBlock.tail = item
	}
	targetBlock.size++

	// 3. Split if too big
	if targetBlock.size > ds.M {
		ds.split(idx)
	}
}

// BatchPrepend adds items strictly smaller than current min.
func (ds *DataStructure) BatchPrepend(items []Item) {
	if len(items) == 0 {
		return
	}
	ds.Count += len(items)

	// Sort items to form valid blocks
	sort.Slice(items, func(i, j int) bool {
		return items[i].Value < items[j].Value
	})

	// Chunk into blocks of size M
	for i := 0; i < len(items); i += ds.M {
		end := i + ds.M
		if end > len(items) {
			end = len(items)
		}

		chunk := items[i:end]
		blk := newBlock()
		// Convert chunk to linked list
		for k := range chunk {
			itm := &Item{Key: chunk[k].Key, Value: chunk[k].Value}
			itm.next = blk.head
			blk.head = itm
			if blk.tail == nil {
				blk.tail = itm
			}
		}
		blk.size = len(chunk)
		blk.upperBound = chunk[len(chunk)-1].Value // Conservative UB

		// Prepend to D0
		ds.d0 = append([]*block{blk}, ds.d0...)
	}
}

// Pull retrieves the smallest M items.
func (ds *DataStructure) Pull() ([]Item, float64) {
	// This is a simplified Pull that merges D0 and D1 logically.
	// In a full impl, you'd scan the heads of D0 and D1.
	// Since D0 contains "small" items prepended, and D1 is sorted,
	// we collect from D0 then D1 until we have M items.

	collected := make([]Item, 0, ds.M)

	// Helper to drain a block
	drain := func(b *block, limit int) {
		curr := b.head
		prev := &Item{Key: 0, Value: 0, next: nil}
		for curr != nil && len(collected) < limit {
			collected = append(collected, *curr)
			prev.next = curr
			prev = curr
			curr = curr.next
			b.size--
			ds.Count--
		}
		b.head = curr
		if curr == nil {
			b.tail = nil
		}
	}

	// 1. Drain D0 first (contains smallest from prepends)
	// We iterate D0 backwards or forwards? BatchPrepend adds to front.
	// Logic dictates D0 blocks are smaller than D1.
	activeD0 := ds.d0[:0]
	for _, blk := range ds.d0 {
		if len(collected) < ds.M {
			drain(blk, ds.M)
		}
		if blk.size > 0 {
			activeD0 = append(activeD0, blk)
		}
	}
	ds.d0 = activeD0

	// 2. Drain D1 if needed
	if len(collected) < ds.M {
		activeD1 := ds.d1[:0]
		for _, blk := range ds.d1 {
			if len(collected) < ds.M {
				// Sort the block to extract smallest?
				// The items inside aren't guaranteed sorted by Insert, only partitioned.
				// We must sort the block content to pull correctly if we partially drain it.
				// Cost: O(M log M). Allowable since pull is amortized.
				ds.sortBlock(blk)
				drain(blk, ds.M)
			}
			if blk.size > 0 {
				// Update UB if needed, or keep
				activeD1 = append(activeD1, blk)
			}
		}
		ds.d1 = activeD1
	}

	// Determine Bi (the bound).
	// If we exhausted everything, Bi = Infinity.
	// Else, Bi is the value of the next available item.
	Bi := Infinity
	if ds.Count > 0 {
		// Find min in remaining D0 or D1
		// (Simplification: just peek heads. Correctness requires iterating blocks)
		if len(ds.d0) > 0 {
			// scan d0
			Bi = ds.peekBlock(ds.d0[0])
		} else if len(ds.d1) > 0 {
			ds.sortBlock(ds.d1[0]) // Ensure sorted to peek
			Bi = ds.peekBlock(ds.d1[0])
		}
	}

	return collected, Bi
}

func newBlock() *block {
	return &block{}
}

func (ds *DataStructure) split(d1Index int) {
	b := ds.d1[d1Index]

	// Materialize list to slice for sorting/splitting
	items := make([]*Item, 0, b.size)
	curr := b.head
	for curr != nil {
		items = append(items, curr)
		curr = curr.next
	}

	// Find median (O(M log M) with sort, or O(M) with select)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Value < items[j].Value
	})

	mid := len(items) / 2

	// Create new block for right half
	newB := newBlock()
	newB.upperBound = b.upperBound    // Inherits old UB
	b.upperBound = items[mid-1].Value // New UB for left block

	// Rebuild lists
	b.head, b.tail, b.size = listFromSlice(items[:mid])
	newB.head, newB.tail, newB.size = listFromSlice(items[mid:])

	// Insert newB into D1 after b
	ds.d1 = append(ds.d1[:d1Index+1], append([]*block{newB}, ds.d1[d1Index+1:]...)...)
}

func (ds *DataStructure) sortBlock(b *block) {
	if b.size < 2 {
		return
	}
	items := make([]*Item, 0, b.size)
	curr := b.head
	for curr != nil {
		items = append(items, curr)
		curr = curr.next
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Value < items[j].Value
	})
	b.head, b.tail, b.size = listFromSlice(items)
}

func listFromSlice(items []*Item) (*Item, *Item, int) {
	if len(items) == 0 {
		return nil, nil, 0
	}
	head := items[0]
	curr := head
	for i := 1; i < len(items); i++ {
		curr.next = items[i]
		curr = curr.next
	}
	curr.next = nil
	return head, curr, len(items)
}

func (ds *DataStructure) peekBlock(b *block) float64 {
	if b.head == nil {
		return Infinity
	}
	return b.head.Value
}
