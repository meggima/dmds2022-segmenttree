package segmenttree

type TreeFunction func(Addable, Addable) Addable

type SegmentTreeImpl struct {
	root            *Node
	operation       TreeFunction
	branchingFactor uint32
	nextNodeId      uint32
}

func NewSegmentTree(branchingFactor uint32, operation TreeFunction) *SegmentTreeImpl {
	tree := &SegmentTreeImpl{
		branchingFactor: branchingFactor,
		operation:       operation,
		nextNodeId:      1,
	}

	tree.root = tree.newNode()
	tree.root.values[0] = 0

	return tree
}

func (t *SegmentTreeImpl) newNode() *Node {
	node := &Node{
		nodeId:   t.nextNodeId,
		n:        0,
		keys:     make([]uint32, t.branchingFactor+1), // + 1 to account for an interval being split into three intervals
		values:   make([]Float, t.branchingFactor+2),  // + 2 to account for an interval being split into three intervals
		children: make([]*Node, t.branchingFactor+2),  // + 2 to account for an interval being split into three intervals
		isLeaf:   true,
		parent:   nil,
		tree:     t,
	}

	t.nextNodeId += 1

	return node
}

func (t *SegmentTreeImpl) GetAtInstant(instant uint32) float32 {
	return -1
}

func (t *SegmentTreeImpl) GetWithinInterval(interval Interval) []ValueIntervalTuple {
	return nil
}

func (t *SegmentTreeImpl) Insert(value ValueIntervalTuple) {

}

func (t *SegmentTreeImpl) Delete(value ValueIntervalTuple) {

}
