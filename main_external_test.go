package merkleTree_test

import (
	"testing"

	"github.com/tclairet/merkle-tree"
)

func TestMerkleTree(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		tree, err := merkleTree.New([]string{"a", "b"})
		if err != nil {
			t.Fatal(err)
		}

		level0, err := tree.Level(0)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := level0[0], merkleTree.Sha256hash(merkleTree.Sha256hash("a"), merkleTree.Sha256hash("b")); got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
		if got, want := tree.Height(), 2; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}
