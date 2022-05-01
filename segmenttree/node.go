package segmenttree

type Node struct {
	nodeId   uint32
	n        uint32
	keys     []uint32
	values   []Float
	children []*Node
	parent   *Node
	tree     *SegmentTreeImpl
	isLeaf   bool
}
