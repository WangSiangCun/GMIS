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
