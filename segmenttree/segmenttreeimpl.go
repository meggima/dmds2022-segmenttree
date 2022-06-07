package segmenttree

type Aggregate struct {
	operation       func(Addable, Addable) Addable
	additionElement func(Addable) Addable
	neutralElement  Addable
}

type SegmentTreeImpl struct {
	root            *Node
	aggregate       Aggregate
	branchingFactor uint32
}

func NewSegmentTree(branchingFactor uint32, aggregate Aggregate) *SegmentTreeImpl {
	tree := &SegmentTreeImpl{
		branchingFactor: branchingFactor,
		aggregate:       aggregate,
	}

	tree.root = tree.newNode()
	tree.root.values[0] = Float(0)

	return tree
}

func (t *SegmentTreeImpl) newNode() *Node {
	node := &Node{
		keys:     make([]uint32, t.branchingFactor+1),  // + 1 to account for an interval being split into three intervals
		values:   make([]Addable, t.branchingFactor+2), // + 2 to account for an interval being split into three intervals
		children: make([]*Node, t.branchingFactor+2),   // + 2 to account for an interval being split into three intervals
		isLeaf:   true,
		parent:   nil,
		tree:     t,
	}

	return node
}

func (tree *SegmentTreeImpl) GetAtInstant(instant uint32) Addable {
	return tree.lookup(tree.root, instant)
}

func (tree *SegmentTreeImpl) GetWithinInterval(interval Interval) []ValueIntervalTuple {
	return tree.rangeQuery(tree.root, interval, tree.aggregate.neutralElement)
}

func (tree *SegmentTreeImpl) Insert(value ValueIntervalTuple) {
	valueToInsert := tree.aggregate.additionElement(value.value)

	tree.insert(tree.root, ValueIntervalTuple{value: valueToInsert, interval: value.interval})
}

func (tree *SegmentTreeImpl) Delete(value ValueIntervalTuple) {
	valueToInsert := tree.aggregate.additionElement(value.value)

	valueToInsert = valueToInsert.Inverse()

	tree.insert(tree.root, ValueIntervalTuple{value: valueToInsert, interval: value.interval})

}

func (tree *SegmentTreeImpl) lookup(node *Node, instant uint32) Addable {
	var intervalIndex = node.findIntervalIndex(instant)

	if node.isLeaf {
		return node.values[intervalIndex]
	}

	return tree.aggregate.operation(node.values[intervalIndex], tree.lookup(node.children[intervalIndex], instant))
}

func (tree *SegmentTreeImpl) rangeQuery(node *Node, interval Interval, value Addable) []ValueIntervalTuple {
	var result []ValueIntervalTuple = make([]ValueIntervalTuple, 0)

	for index, nodeInterval := range node.getIntervals() {
		intersection := interval.IntersectionWith(nodeInterval)

		if intersection.GetLength() == 0 {
			continue
		}

		if node.isLeaf {
			newTuple := ValueIntervalTuple{
				value:    tree.aggregate.operation(node.values[index], value),
				interval: intersection,
			}
			result = append(result, newTuple)
		} else {
			childResult := tree.rangeQuery(node.children[index], interval, tree.aggregate.operation(node.values[index], value))
			result = append(result, childResult...)
		}
	}

	return result
}

func (tree *SegmentTreeImpl) insert(node *Node, tupleToInsert ValueIntervalTuple) {
	has_next_interval := true
	index := 0
	intervals := node.getIntervals()
	for has_next_interval {
		nodeInterval := intervals[index]

		intersection := nodeInterval.IntersectionWith(tupleToInsert.interval)

		if intersection.GetLength() == 0 {
			// Do nothing
		} else if index < len(node.values) && node.values[index] == tree.aggregate.operation(node.values[index], tupleToInsert.value) {
			// Do nothing
		} else if nodeInterval.IsSubsetOf(tupleToInsert.interval) {
			node.values[index] = tree.aggregate.operation(node.values[index], tupleToInsert.value)
		} else {
			if !node.isLeaf {
				tree.insert(node.children[index], tupleToInsert)
			} else {
				index += node.insert(index, tupleToInsert)
				intervals = node.getIntervals() // recalculate as they might have changed
			}
		}
		if index+1 >= len(intervals) {
			has_next_interval = false
		}
		index++
	}
	if node.size()+1 > tree.branchingFactor {
		node.split()
	}

}
