package trie

import (
	"encoding/json"
	"math/big"

	"github.com/NethermindEth/juno/pkg/crypto/pedersen"
)

// encoding represents the enccoding of a node in a binary tree
// represented by the triplet (length, path, bottom).
type encoding struct {
	Length uint8    `json:"length"`
	Path   *big.Int `json:"path"`
	Bottom *big.Int `json:"bottom"`
}

// node represents a node in a binary tree.
type node struct {
	encoding
	Hash *big.Int `json:"hash"`
	Next []node   `json:"-"`
}

// newNode initialises a new node with two null links.
func newNode() node {
	return node{Next: make([]node, 2)}
}

// bytes returns a JSON byte representation of a node.
func (n *node) bytes() []byte {
	b, _ := json.Marshal(n)
	return b
}

// clear sets the links in the node n to null. This is done to conserve
// memory after a node has been committed to storage.
func (n *node) clear() {
	n.Next = nil
}

// isEmpty returns true if the in-memory representation of a node is
// the empty node i.e. encoded by the triplet (0, 0, 0).
func (n *node) isEmpty() bool {
	return n.Next == nil
}

// updateHash updates the node hash.
func (n *node) updateHash() {
	if n.Length == 0 {
		n.Hash = new(big.Int).Set(n.Bottom)
	} else {
		h, _ := pedersen.Digest(n.Bottom, n.Path)
		n.Hash = h.Add(h, new(big.Int).SetUint64(uint64(n.Length)))
	}
}