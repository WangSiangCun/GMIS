package BPTree

import (
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
func (t *BPTree) GetData() map[int64]interface{} {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.getData(t.root)
}
func (t *BPTree) getData(node *BPNode) map[int64]interface{} {
	data := make(map[int64]interface{})
	for {
		if len(node.ChildNodes) > 0 {
			for i := 0; i < len(node.ChildNodes); i++ {
				data[node.ChildNodes[i].MaxKey] = t.getData(node.ChildNodes[i])
			}
			break
		} else {
			for i := 0; i < len(node.Items); i++ {
				data[node.Items[i].Key] = node.Items[i].Val
			}
			break
		}
	}
	return data
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
func (t *BPTree) Set(key int64, value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.InsertItem(nil, t.root, key, value)
}
func (t *BPTree) InsertItem(parent *BPNode, node *BPNode, key int64, value interface{}) {
	for i := 0; i < len(node.ChildNodes); i++ {
		if key <= node.ChildNodes[i].MaxKey || i == len(node.ChildNodes)-1 {
			t.InsertItem(node, node.ChildNodes[i], key, value)
			break
		}
	}

	//叶子结点，添加数据
	if len(node.ChildNodes) < 1 {
		node.InsertItem(key, value)
	}

	//结点分裂
	nodeNew := t.splitNode(node)
	if nodeNew != nil {
		//若父结点不存在，则创建一个父节点
		if parent == nil {
			parent = NewNode(t.width)
			parent.AddChildNode(node)
			t.root = parent
		}
		//添加结点到父亲结点
		parent.AddChildNode(nodeNew)
	}
}
func (t *BPTree) splitNode(node *BPNode) *BPNode {
	if len(node.ChildNodes) > t.width {
		//创建新结点
		halfW := t.width/2 + 1
		node2 := NewNode(t.width)
		node2.ChildNodes = append(node2.ChildNodes, node.ChildNodes[halfW:len(node.ChildNodes)]...)
		node2.MaxKey = node2.ChildNodes[len(node2.ChildNodes)-1].MaxKey

		//修改原结点数据
		node.ChildNodes = node.ChildNodes[0:halfW]
		node.MaxKey = node.ChildNodes[len(node.ChildNodes)-1].MaxKey

		return node2
	} else if len(node.Items) > t.width {
		//创建新结点
		halfW := t.width/2 + 1
		node2 := NewNode(t.width)
		node2.Items = append(node2.Items, node.Items[halfW:len(node.Items)]...)
		node2.MaxKey = node2.Items[len(node2.Items)-1].Key

		//修改原结点数据
		node.Next = node2
		node.Items = node.Items[0:halfW]
		node.MaxKey = node.Items[len(node.Items)-1].Key

		return node2
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
		node.Items = append([]BPItem{item}, node.Items...)
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
