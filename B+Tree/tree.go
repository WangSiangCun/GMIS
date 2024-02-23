package BPTree

import (
	"fmt"
	"sync"
)

// BPTree
// root：指向B+树的根结点
// width：用于表示B+树的阶M
// halfW：用于[M/2]=ceil(M/2)
type BPTree struct {
	mutex sync.RWMutex
	root  *BPNode
	width int
	halfW int
}

func NewBPTree(width int) *BPTree {
	if width < 3 {
		width = 3
	}
	var bt = &BPTree{}
	bt.root = NewNode(width)
	bt.width = width
	bt.halfW = (bt.width + 1) / 2
	return bt
}
func (t *BPTree) String() string {
	return t.getTreeString(t.root)
}
func (t *BPTree) getTreeString(node *BPNode) (res string) {

	now := node
	for now != nil {
		for i := 0; i < len(now.Items); i++ {
			res += fmt.Sprintf(" %v ", now.Items[i].Key)
		}
		res += fmt.Sprintf("|")
		now = now.Next
	}
	res += "\n"
	if len(node.ChildNodes) != 0 {
		res += t.getTreeString(node.ChildNodes[0])
	}

	return res

}

// Insert 递归插入
func (t *BPTree) Insert(lastNode *BPNode, nowNode *BPNode, key int64, value interface{}) {

	//找到叶子节点插入item
	if nowNode.IsLeafNode() {
		nowNode.InsertItem(key, value)
	} else {
		for i := 0; i < len(nowNode.ChildNodes); i++ {
			//从子节点中找到合适的
			if key <= nowNode.ChildNodes[i].MaxKey || i == len(nowNode.ChildNodes)-1 {
				t.Insert(nowNode, nowNode.ChildNodes[i], key, value)
				break
			}
		}

	}
	//找到之后的处理
	if nowNode.IsNeedSplitNode() {
		if res := nowNode.Split(lastNode); res != nil {
			t.root = res
		}
	}

}
func (t *BPTree) Get(key int64) interface{} {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	node := t.root
	for i := 0; i < len(node.ChildNodes); i++ {
		if key <= node.ChildNodes[i].MaxKey {
			node = node.ChildNodes[i]
			i = 0
		}
	}

	//没有到达叶子结点
	if len(node.ChildNodes) > 0 {
		return nil
	}

	for i := 0; i < len(node.Items); i++ {
		if node.Items[i].Key == key {
			return node.Items[i].Val
		}
	}
	return nil
}
func (t *BPTree) Remove(key int64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.DeleteItem(nil, t.root, key)
}
func (t *BPTree) DeleteItem(parent *BPNode, node *BPNode, key int64) {
	for i := 0; i < len(node.ChildNodes); i++ {
		if key <= node.ChildNodes[i].MaxKey {
			t.DeleteItem(node, node.ChildNodes[i], key)
			break
		}
	}

	if len(node.ChildNodes) < 1 {
		//删除记录后若结点的子项<m/2，则从兄弟结点移动记录，或者合并结点
		node.DeleteItem(key)
		if len(node.Items) < t.halfW {
			t.itemMoveOrMerge(parent, node)
		}
	} else {
		//若结点的子项<m/2，则从兄弟结点移动记录，或者合并结点
		node.MaxKey = node.ChildNodes[len(node.ChildNodes)-1].MaxKey
		if len(node.ChildNodes) < t.halfW {
			t.childMoveOrMerge(parent, node)
		}
	}
}
func (t *BPTree) childMoveOrMerge(parent *BPNode, node *BPNode) {
	if parent == nil {
		return
	}

	//获取兄弟结点
	var node1 *BPNode = nil
	var node2 *BPNode = nil
	for i := 0; i < len(parent.ChildNodes); i++ {
		if parent.ChildNodes[i] == node {
			if i < len(parent.ChildNodes)-1 {
				node2 = parent.ChildNodes[i+1]
			} else if i > 0 {
				node1 = parent.ChildNodes[i-1]
			}
			break
		}
	}

	//将左侧结点的子结点移动到删除结点
	if node1 != nil && len(node1.ChildNodes) > t.halfW {
		item := node1.ChildNodes[len(node1.ChildNodes)-1]
		node1.ChildNodes = node1.ChildNodes[0 : len(node1.ChildNodes)-1]
		node.ChildNodes = append([]*BPNode{item}, node.ChildNodes...)
		return
	}

	//将右侧结点的子结点移动到删除结点
	if node2 != nil && len(node2.ChildNodes) > t.halfW {
		item := node2.ChildNodes[0]
		node2.ChildNodes = node1.ChildNodes[1:]
		node.ChildNodes = append(node.ChildNodes, item)
		return
	}

	if node1 != nil && len(node1.ChildNodes)+len(node.ChildNodes) <= t.width {
		node1.ChildNodes = append(node1.ChildNodes, node.ChildNodes...)
		parent.DeleteChildNode(node)
		return
	}

	if node2 != nil && len(node2.ChildNodes)+len(node.ChildNodes) <= t.width {
		node.ChildNodes = append(node.ChildNodes, node2.ChildNodes...)
		parent.DeleteChildNode(node2)
		return
	}
}
func (t *BPTree) itemMoveOrMerge(parent *BPNode, node *BPNode) {
	//获取兄弟结点
	var node1 *BPNode = nil
	var node2 *BPNode = nil
	for i := 0; i < len(parent.ChildNodes); i++ {
		if parent.ChildNodes[i] == node {
			if i < len(parent.ChildNodes)-1 {
				node2 = parent.ChildNodes[i+1]
			} else if i > 0 {
				node1 = parent.ChildNodes[i-1]
			}
			break
		}
	}

	//将左侧结点的记录移动到删除结点
	if node1 != nil && len(node1.Items) > t.halfW {
		item := node1.Items[len(node1.Items)-1]
		node1.Items = node1.Items[0 : len(node1.Items)-1]
		node1.MaxKey = node1.Items[len(node1.Items)-1].Key
		node.Items = append([]*BPItem{item}, node.Items...)
		return
	}

	//将右侧结点的记录移动到删除结点
	if node2 != nil && len(node2.Items) > t.halfW {
		item := node2.Items[0]
		node2.Items = node1.Items[1:]
		node.Items = append(node.Items, item)
		node.MaxKey = node.Items[len(node.Items)-1].Key
		return
	}

	//与左侧结点进行合并
	if node1 != nil && len(node1.Items)+len(node.Items) <= t.width {
		node1.Items = append(node1.Items, node.Items...)
		node1.Next = node.Next
		node1.MaxKey = node1.Items[len(node1.Items)-1].Key
		parent.DeleteChildNode(node)
		return
	}

	//与右侧结点进行合并
	if node2 != nil && len(node2.Items)+len(node.Items) <= t.width {
		node.Items = append(node.Items, node2.Items...)
		node.Next = node2.Next
		node.MaxKey = node.Items[len(node.Items)-1].Key
		parent.DeleteChildNode(node2)
		return
	}
}
