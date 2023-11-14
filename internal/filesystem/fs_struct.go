package fs

import (
	"container/list"
)

type Type int8
type Id int8

const (
	Dir Type = iota
	File
)

var currId Id = 1

func NextId() Id {
	t := currId
	currId++
	return t
}

type Tree struct {
	root *Node
}

func NewTree() *Tree {
	root := NewRoot(Dir, "./")
	return &Tree{
		root: root,
	}
}

type WalkFn func(node *Node) error
type SearchFn func(in *Node) *Node

type Node struct {
	id       Id
	type_    Type
	item     IItem
	children *list.List
	parent   *Node
}

// NewNode with id incremented
func NewNode(_type Type, path string, root *Node) *Node {
	var item IItem
	if path != "" {
		item, _ = CreateItem(path)
	}
	return &Node{
		id:       NextId(),
		type_:    _type,
		item:     item,
		children: list.New(),
		parent:   root,
	}
}

func NewLeaf(_type Type, path string, root *Node) *Node {
	var item IItem
	item, _ = CreateItem(path)
	return &Node{
		id:     NextId(),
		type_:  _type,
		item:   item,
		parent: root,
	}
}

func NewRoot(_type Type, path string) *Node {
	var item IItem
	if path != "" {
		item, _ = CreateItem(path)
	}
	return &Node{
		id:       NextId(),
		type_:    _type,
		item:     item,
		children: list.New(),
	}
}

func (node *Node) Children() []*Node {
	var nodes []*Node
	for curr := node.children.Front(); curr != nil; curr = curr.Next() {
		nodes = append(nodes, curr.Value.(*Node))
	}
	return nodes
}

func (node *Node) Parent() *Node {
	if node.parent != nil {
		return node.parent
	}
	return nil
}

func (node *Node) Type() Type {
	return node.type_
}

func (node *Node) Id() Id {
	return node.id
}

func (node *Node) Item() IItem {
	return node.item
}

func (node *Node) Add(n *Node) {
	node.children.PushBack(n)
}

func (tree *Tree) Exists(node *Node) *Node {
	return tree.exists(tree.root, node)
}

func (tree *Tree) exists(in *Node, needle *Node) *Node {
	if in == needle {
		return in
	}
	if in.children.Len() > 0 {
		for curr := in.children.Front(); curr != nil; curr = curr.Next() {
			iter := curr.Value.(*Node)
			found := tree.exists(iter, needle)
			if found != nil {
				return found
			}
		}
	}
	return nil
}

func (tree *Tree) FindId(id Id) *Node {
	return tree.findBy(tree.Root(), func(next *Node) *Node {
		if next.Id() == id {
			return next
		}
		return nil
	})
}

func (tree *Tree) FindNode(node *Node) *Node {
	return tree.findBy(tree.Root(), func(next *Node) *Node {
		if next == node {
			return node
		}
		return nil
	})
}

func (tree *Tree) FindBy(fn SearchFn) *Node {
	return tree.findBy(tree.Root(), fn)
}

func (tree *Tree) findBy(in *Node, fn SearchFn) *Node {
	v := fn(in)
	if v != nil {
		return in
	}
	if in.children.Len() > 0 {
		for curr := in.children.Front(); curr != nil; curr = curr.Next() {
			iter := curr.Value.(*Node)
			found := tree.findBy(iter, fn)
			if found != nil {
				return found
			}
		}
	}
	return nil
}

// LWalk
// Deprecated
func (tree *Tree) LWalk(node *Node, fn WalkFn) error {
	if !node.IsLeaf() {
		for n := node.children.Front(); n != nil; n = n.Next() {
			iter := n.Value.(*Node)
			e := tree.LWalk(iter, fn)
			if e != nil {
				return e
			}
		}
	}
	return fn(node)
}

func (tree *Tree) LWalkLvl(node *Node, fn WalkFn, lvl int8) error {
	// including top level root
	if lvl >= 0 {
		lvl--
		if !node.IsLeaf() {
			for n := node.children.Front(); n != nil; n = n.Next() {
				iter := n.Value.(*Node)
				e := tree.LWalkLvl(iter, fn, lvl)
				if e != nil {
					return e
				}
			}
		}
	}
	return fn(node)
}

func (node *Node) IsRoot() bool {
	if node.parent == nil {
		return true
	}
	return false
}

func (type_ Type) String() string {
	if type_ == Dir {
		return "Dir"
	}
	return "File"
}

func (node *Node) IsLeaf() bool {
	if node.Type() == File || node.children.Len() == 0 {
		return true
	}
	return false
}

func (tree *Tree) Root() *Node {
	return tree.root
}

func (node *Node) PathTrace(traceTo *Node, path string) string {
	var temp string
	if node.parent != nil {
		temp = node.parent.item.Name() + "/" + node.item.Name()
		node.parent.PathTrace(traceTo, temp)
	}
	return path
}