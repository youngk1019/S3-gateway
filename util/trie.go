package util

type Trie[T comparable] struct {
	Node   map[T]*Trie[T]
	Value  T
	IsLeaf bool
}

func NewTrie[T comparable](val T) *Trie[T] {
	q := &Trie[T]{
		Node:   make(map[T]*Trie[T]),
		Value:  val,
		IsLeaf: false,
	}
	return q
}

func (t *Trie[T]) Insert(isLeaf bool, items ...T) {
	p := t
	for _, u := range items {
		if p.Node[u] == nil {
			p.Node[u] = NewTrie[T](u)
			p = p.Node[u]
		} else {
			p = p.Node[u]
		}
	}
	p.IsLeaf = isLeaf
}
