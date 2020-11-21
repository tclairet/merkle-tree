package merkleTree

import (
	"reflect"
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	if got, want := Sha256hash("a", "b"), Sha256hash("ab"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestMerkleTree(t *testing.T) {
	fakeHash := func(left string, right ...string) string {
		if len(right) == 0 {
			return left
		}
		return left + right[0]
	}

	t.Run("from", func(t *testing.T) {
		tree := merkleTree{hash: fakeHash}

		cases := []struct {
			inputs         []string
			expectedRoot   string
			expectedHeight int
		}{
			{[]string{"a"}, "a", 1},
			{[]string{"a", "b"}, "ab", 2},
			{[]string{"a", "b", "c"}, "abc", 3},
			{[]string{"a", "b", "c", "d"}, "abcd", 3},
			{[]string{"a", "b", "c", "d", "e"}, "abcde", 4},
			{[]string{"a", "b", "c", "d", "e", "f"}, "abcdef", 4},
			{[]string{"a", "b", "c", "d", "e", "f"}, "abcdef", 4},
			{[]string{"a", "b", "c", "d", "e", "f", "g", "h"}, "abcdefgh", 4},
			{[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}, "abcdefghi", 5},
		}

		for _, c := range cases {
			tree.from(c.inputs)

			if got, want := tree.root, c.expectedRoot; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
			if got, want := tree.Height(), c.expectedHeight; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
		}
	})

	t.Run("level", func(t *testing.T) {
		tree := merkleTree{hash: fakeHash}

		cases := []struct {
			inputs   []string
			level    int
			expected []string
		}{
			{[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}, 0, []string{"abcdefghi"}},
			{[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}, 1, []string{"abcdefgh", "i"}},
			{[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}, 2, []string{"abcd", "efgh", "i"}},
			{[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}, 3, []string{"ab", "cd", "ef", "gh", "i"}},
			{[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}, 4, []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}},
		}

		for _, c := range cases {
			tree.from(c.inputs)

			leaves, err := tree.Level(c.level)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := leaves, c.expected; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		}
	})

	t.Run("invalid level", func(t *testing.T) {
		tree, _ := New([]string{"a", "b"})
		cases := []struct {
			tree     *merkleTree
			level    int
			expected string
		}{
			{&merkleTree{}, 0, "merkle tree not initialized"},
			{tree, 2, "cannot retrieve level"},
		}

		for _, c := range cases {
			_, err := c.tree.Level(c.level)
			if got, want := err.Error(), c.expected; !strings.Contains(got, want) {
				t.Fatalf("'%v', does not contains %v", got, want)
			}
		}
	})

	t.Run("new", func(t *testing.T) {
		tree, err := New([]string{"a", "b"})
		if err != nil {
			t.Fatal(err)
		}

		if got, want := tree.root, Sha256hash(Sha256hash("a"), Sha256hash("b")); got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
		if got, want := tree.height, 2; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("new invalid inputs", func(t *testing.T) {
		_, err := New(nil)
		if got, want := err.Error(), "invalid inputs"; !strings.Contains(got, want) {
			t.Fatalf("'%v', does not contains %v", got, want)
		}
	})
}
