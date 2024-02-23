package BPTree

import (
	"fmt"
	"testing"
)

func TestBPTree_Insert(t *testing.T) {
	bpt := NewBPTree(3)
	for i := 1; i <= 1000000; i++ {
		bpt.Insert(nil, bpt.root, int64(i), i)
		//fmt.Println(bpt.String())

	}
	fmt.Println(bpt.String())

}
