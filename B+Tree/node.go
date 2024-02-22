package BPTree

import "fmt"

// BPItem 用于数据记录。
type BPItem struct {
	Key int64
	Val interface{}
}

func (item *BPItem) String() string {
	return fmt.Sprintf(" Key:%v-value:%v  ", item.Key, item.Val)
}

// BPNode
// MaxKey：用于存储子树的最大关键字
// ChildNodes：结点的子树（叶子结点的ChildNodes=nil）
// Items：叶子结点的数据记录（索引结点的Items=nil）
// Next：叶子结点中指向下一个叶子结点，用于实现叶子结点链表
type BPNode struct {
	MaxKey     int64
	ChildNodes []*BPNode
	Items      []BPItem
	Next       *BPNode
}

func NewNode(width int) *BPNode {
	var node = &BPNode{}
	node.Items = make([]BPItem, width+1)
	node.Items = node.Items[0:0]
	return node
}

// InsertItem 插入Node的Item，并维护Items有序
func (node *BPNode) InsertItem(key int64, value interface{}) {
	item := BPItem{key, value}
	num := len(node.Items)
	if num < 1 {
		node.Items = append(node.Items, item)
		node.MaxKey = item.Key
		return
	} else if key < node.Items[0].Key {
		node.Items = append([]BPItem{item}, node.Items...)
		return
	} else if key > node.Items[num-1].Key {
		node.Items = append(node.Items, item)
		node.MaxKey = item.Key
		return
	}

	for i := 0; i < num; i++ {
		if node.Items[i].Key > key {
			node.Items = append(node.Items, BPItem{})
			copy(node.Items[i+1:], node.Items[i:])
			node.Items[i] = item
			return
		} else if node.Items[i].Key == key {
			node.Items[i] = item //直接赋值，等于修改了
			return
		}
	}
}

// DeleteItem 删除对应Key的Item，不存在则返回false，删除成功即true
func (node *BPNode) DeleteItem(key int64) bool {
	num := len(node.Items)
	for i := 0; i < num; i++ {
		//大于要查找的key 说明后面的key也不会存在相等了 因为是递增的
		if node.Items[i].Key > key {
			return false
		} else if node.Items[i].Key == key {
			copy(node.Items[i:], node.Items[i+1:])
			if len(node.Items) != 1 {
				node.Items = node.Items[0 : len(node.Items)-1]
				node.MaxKey = node.Items[len(node.Items)-1].Key
			} else {
				node.Items = nil
				node.MaxKey = 0
			}

			return true
		}
	}
	return false
}
func (node *BPNode) String() string {
	res := ""
	for i := 0; i < len(node.Items); i++ {
		res += fmt.Sprintf("%v-%v  ", i, node.Items[i].String())
	}
	res += fmt.Sprintf("MaxKey:%v ChildNodes:%v Next:%v \n", node.MaxKey, node.ChildNodes, node.Next)
	return res
}

// AddChildNode 插入Node的孩子节点，并维护孩子节点们按照MaxKey有序
func (node *BPNode) AddChildNode(child *BPNode) {
	num := len(node.ChildNodes)
	if num < 1 {
		node.ChildNodes = append(node.ChildNodes, child)
		node.MaxKey = child.MaxKey
		return
	} else if child.MaxKey < node.ChildNodes[0].MaxKey {
		node.ChildNodes = append([]*BPNode{child}, node.ChildNodes...)
		return
	} else if child.MaxKey > node.ChildNodes[num-1].MaxKey {
		node.ChildNodes = append(node.ChildNodes, child)
		node.MaxKey = child.MaxKey
		return
	}

	for i := 0; i < num; i++ {
		if node.ChildNodes[i].MaxKey > child.MaxKey {
			node.ChildNodes = append(node.ChildNodes, nil)
			copy(node.ChildNodes[i+1:], node.ChildNodes[i:])
			node.ChildNodes[i] = child
			return
		}
	}
}

// DeleteChildNode 删除指定的Node的孩子节点，并维护孩子节点们按照MaxKey有序
func (node *BPNode) DeleteChildNode(child *BPNode) bool {
	num := len(node.ChildNodes)
	for i := 0; i < num; i++ {
		if node.ChildNodes[i] == child {
			copy(node.ChildNodes[i:], node.ChildNodes[i+1:])
			node.ChildNodes = node.ChildNodes[0 : len(node.ChildNodes)-1]
			node.MaxKey = node.ChildNodes[len(node.ChildNodes)-1].MaxKey
			return true
		}
	}
	return false
}
