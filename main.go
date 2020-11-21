package merkleTree

import (
	"crypto/sha256"
	"fmt"
)

var _ MerkleTree = (*merkleTree)(nil)

type HashFn func(string, ...string) string

type MerkleTree interface {
	Height() int
	Level(int) ([]string, error)
}

type merkleTree struct {
	nodes  map[string]Node
	root   string
	height int

	hash HashFn
}

func New(inputs []string) (*merkleTree, error) {
	tree := &merkleTree{
		hash: Sha256hash,
	}
	if err := tree.from(inputs); err != nil {
		return nil, err
	}
	return tree, nil
}

func (tree *merkleTree) Level(index int) ([]string, error) {
	if tree.height == 0 {
		return nil, fmt.Errorf("merkle tree not initialized")
	}

	if index >= tree.height {
		return nil, fmt.Errorf("cannot retrieve level '%d' height is '%d'", index, tree.height)
	}

	parents := []string{tree.root}
	if index == 0 {
		return parents, nil
	}

	var nodes []string
	for i := 0; i < index; i++ {
		nodes = []string{}
		for _, child := range parents {
			nodes = append(nodes, tree.childOf(child)...)
		}
		parents = nodes
	}

	return nodes, nil
}

func (tree *merkleTree) Height() int {
	return tree.height
}

func (tree *merkleTree) from(inputs []string) error {
	if inputs == nil || len(inputs) == 0 {
		return fmt.Errorf("invalid inputs")
	}

	tree.height = 1
	tree.nodes = make(map[string]Node)

	var hashes []string
	for _, data := range inputs {
		hashes = append(hashes, tree.hash(data))
	}

	for _, hash := range hashes {
		tree.nodes[hash] = Node{
			hash: hash,
		}
	}

	for len(hashes) != 1 {
		hashes = tree.buildBranch(hashes)
	}
	tree.root = hashes[0]

	return nil
}

func (tree *merkleTree) buildBranch(hashes []string) []string {
	var newNodes []string

	for i := 0; i < (len(hashes) / 2); i++ {
		hash := tree.hash(hashes[i*2], hashes[i*2+1])
		tree.nodes[hash] = Node{
			hash:       hash,
			leftChild:  hashes[i*2],
			rightChild: hashes[i*2+1],
		}
		newNodes = append(newNodes, hash)
	}

	if isOdd(len(hashes)) {
		newNodes = append(newNodes, hashes[len(hashes)-1])
	}

	tree.height++
	return newNodes
}

func (tree *merkleTree) childOf(hash string) []string {
	node := tree.nodes[hash]
	if node.leftChild == "" && node.rightChild == "" {
		return []string{hash}
	}
	return []string{node.leftChild, node.rightChild}
}

type Node struct {
	hash       string
	leftChild  string
	rightChild string
}

func Sha256hash(left string, right ...string) string {
	if len(right) == 0 {
		return fmt.Sprintf("%x", sha256.Sum256([]byte(left)))
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(left+right[0])))
}

func isOdd(number int) bool {
	return number%2 == 1
}
