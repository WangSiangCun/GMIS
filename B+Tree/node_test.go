package BPTree

import (
	"fmt"
	"testing"
)

func BenchmarkInsertItem(b *testing.B) {
	node := NewNode(10000000)
	for i := 0; i < b.N; i++ {
		node.InsertItem(int64(i), i)
	}
	fmt.Println(node)
}
func BenchmarkDeleteItem(b *testing.B) {
	node := NewNode(111111111)
	for i := 0; i < b.N; i++ {
		node.InsertItem(int64(i), i)
		node.DeleteItem(int64(i))
	}
	fmt.Println(node)
}
func TestBPNode_Split(t *testing.T) {
	n := BPNode{
		MaxKey: 5,
		ChildNodes: []*BPNode{&BPNode{
			MaxKey:     0,
			ChildNodes: nil,
			Items:      nil,
			Next:       nil,
			Wight:      0,
		}, &BPNode{
			MaxKey:     0,
			ChildNodes: nil,
			Items:      nil,
			Next:       nil,
			Wight:      0,
		}, &BPNode{
			MaxKey:     0,
			ChildNodes: nil,
			Items:      nil,
			Next:       nil,
			Wight:      0,
		}},
		Items: []*BPItem{&BPItem{1, 1}, &BPItem{3, 2}, &BPItem{5, 2}},
		Next:  nil,
		Wight: 3,
	}
	root := n.Split(nil)
	fmt.Println(root)
}
