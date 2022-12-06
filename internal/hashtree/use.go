package hashtree

import (
	"errors"
	"io"
	"os"

	"github.com/cbergoon/merkletree"
)

// NewHashTree build file to build hash tree
func NewHashTree(chunkPath []string) (*merkletree.MerkleTree, error) {
	if len(chunkPath) == 0 {
		return nil, errors.New("Empty data")
	}
	var list = make([]merkletree.Content, 0)
	for i := 0; i < len(chunkPath); i++ {
		f, err := os.Open(chunkPath[i])
		if err != nil {
			return nil, err
		}
		temp, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		f.Close()
		list = append(list, HashTreeContent{x: string(temp)})
	}

	//Create a new Merkle Tree from the list of Content
	return merkletree.NewTree(list)
}
