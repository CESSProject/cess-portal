package hashtree

import (
	"crypto/sha256"

	"github.com/cbergoon/merkletree"
)

// HashTreeContent implements the Content interface provided by merkletree
// and represents the content stored in the tree.
type HashTreeContent struct {
	x string
}

// CalculateHash hashes the values of a HashTreeContent
func (t HashTreeContent) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.x)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Equals tests for equality of two Contents
func (t HashTreeContent) Equals(other merkletree.Content) (bool, error) {
	return t.x == other.(HashTreeContent).x, nil
}
