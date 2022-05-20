package segmenttree

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindIntervalIndex(t *testing.T) {
	// Arrange
	testData := []struct {
		instant       uint32
		expectedIndex uint32
	}{
		{1, 0},
		{10, 1},
		{21, 2},
		{30, 3},
		{40, 3},
	}

	node := &Node{
		keys: []uint32{10, 20, 30},
	}

	for _, td := range testData {
		// Act
		index := node.findIntervalIndex(td.instant)

		// Assert
		assert.Equal(t, td.expectedIndex, index)
	}
}

func TestFindChildIndex(t *testing.T) {
	// Arrange
	parent := &Node{
		children: make([]*Node, 2),
	}

	child1 := &Node{
		parent: parent,
	}

	child2 := &Node{
		parent: parent,
	}

	parent.children[0] = child1
	parent.children[1] = child2

	// Act
	index1 := parent.findChildIndex(child1)
	index2 := parent.findChildIndex(child2)

	// Assert
	assert.Equal(t, uint32(0), index1)
	assert.Equal(t, uint32(1), index2)
}

func TestFindChildIndexWhenNotChild(t *testing.T) {
	// Assert
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()

	// Arrange
	parent := &Node{
		children: make([]*Node, 2),
	}

	child1 := &Node{
		parent: parent,
	}

	child2 := &Node{
		parent: parent,
	}

	child3 := &Node{
		parent: parent,
	}

	parent.children[0] = child1
	parent.children[1] = child2

	// Act
	_ = parent.findChildIndex(child3)
}

func TestGetIntervalEndWhenRoot(t *testing.T) {
	// Arrange
	testData := []struct {
		index       uint32
		expectedEnd uint32
	}{
		{0, 1},
		{1, 15},
		{2, 30},
		{3, math.MaxUint32},
	}

	node := &Node{
		keys: []uint32{1, 15, 30},
	}

	for _, td := range testData {
		// Act
		intervalEnd := node.getIntervalEnd(td.index)

		// Assert
		assert.Equal(t, td.expectedEnd, intervalEnd)
	}
}

func TestGetIntervalEndWhenNotRoot(t *testing.T) {
	// Arrange
	_, n1, n2, n3, n4 := SetupNodes()

	testData := []struct {
		node        *Node
		index       uint32
		expectedEnd uint32
	}{
		{n1, 2, 15},
		{n2, 1, 30},
		{n3, 2, 45},
		{n4, 1, math.MaxUint32},
	}

	for _, td := range testData {
		// Act
		intervalEnd := td.node.getIntervalEnd(td.index)

		// Assert
		assert.Equal(t, td.expectedEnd, intervalEnd)
	}
}

func TestGetIntervalStartWhenRoot(t *testing.T) {
	// Arrange
	testData := []struct {
		index       uint32
		expectedEnd uint32
	}{
		{0, 0},
		{1, 1},
		{2, 15},
		{3, 30},
	}

	node := &Node{
		keys: []uint32{1, 15, 30},
	}

	for _, td := range testData {
		// Act
		intervalStart := node.getIntervalStart(td.index)

		// Assert
		assert.Equal(t, td.expectedEnd, intervalStart)
	}
}

func TestGetIntervalStartWhenNotRoot(t *testing.T) {
	// Arrange
	_, n1, n2, n3, n4 := SetupNodes()

	testData := []struct {
		node        *Node
		index       uint32
		expectedEnd uint32
	}{
		{n1, 0, 0},
		{n2, 0, 15},
		{n3, 0, 30},
		{n4, 0, 45},
	}

	for _, td := range testData {
		// Act
		intervalStart := td.node.getIntervalStart(td.index)

		// Assert
		assert.Equal(t, td.expectedEnd, intervalStart)
	}
}

func TestGetIntervals(t *testing.T) {
	// Arrange
	n0, n1, n2, n3, n4 := SetupNodes()

	testData := []struct {
		node              *Node
		expectedIntervals []Interval
	}{
		{n0, []Interval{NewInterval(0, 15), NewInterval(15, 30), NewInterval(30, 45), NewInterval(45, math.MaxUint32)}},
		{n1, []Interval{NewInterval(0, 5), NewInterval(5, 10), NewInterval(10, 15)}},
		{n2, []Interval{NewInterval(15, 20), NewInterval(20, 30)}},
		{n3, []Interval{NewInterval(30, 35), NewInterval(35, 40), NewInterval(40, 45)}},
		{n4, []Interval{NewInterval(45, 50), NewInterval(50, math.MaxUint32)}},
	}

	for _, td := range testData {
		// Act
		intervals := td.node.getIntervals()

		// Assert
		assert.Equal(t, td.expectedIntervals, intervals)
	}
}

func TestInsertEmptyNode(t *testing.T) {
	// Arrange
	node := &Node{
		keys:   []uint32{},
		values: []Addable{Float(0)},
		tree: &SegmentTreeImpl{
			aggregate:       Aggregate{Sum, Identity, Float(0)},
			branchingFactor: BRANCHING_FACTOR,
		},
	}
	intervalTuple := ValueIntervalTuple{value: Addable(Float(7)), interval: Interval{start: 17, end: 47}}

	// Act
	node.insert(0, intervalTuple)

	// Assert
	assert.Equal(t, intervalTuple.interval.start, node.keys[0])
	assert.Equal(t, intervalTuple.interval.end, node.keys[1])
	assert.Equal(t, Float(0), node.values[0])
	assert.Equal(t, intervalTuple.value, node.values[1])
	assert.Equal(t, Float(0), node.values[2])

}

func TestInsertSelfContained1(t *testing.T) {
	// Arrange
	tree := &SegmentTreeImpl{
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: BRANCHING_FACTOR,
	}
	node := &Node{
		keys:   []uint32{10, 20},
		values: []Addable{Float(1), Float(2), Float(3)},
		tree:   tree,
	}

	parent := &Node{
		keys:     []uint32{50},
		values:   []Addable{Float(0), Float(0)},
		children: []*Node{node, nil},
		tree:     tree,
	}

	node.parent = parent
	intervalTuple := ValueIntervalTuple{value: Addable(Float(7)), interval: Interval{start: 1, end: 3}}

	// Act
	node.insert(0, intervalTuple)

	// Assert
	assert.Equal(t, uint32(1), node.keys[0])
	assert.Equal(t, uint32(3), node.keys[1])
	assert.Equal(t, uint32(10), node.keys[2])
	assert.Equal(t, uint32(20), node.keys[3])
	assert.Equal(t, Float(1), node.values[0])
	assert.Equal(t, Float(8), node.values[1])
	assert.Equal(t, Float(1), node.values[2])
	assert.Equal(t, Float(2), node.values[3])
	assert.Equal(t, Float(3), node.values[4])
}

func TestInsertSelfContained2(t *testing.T) {
	// Arrange
	tree := &SegmentTreeImpl{
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: BRANCHING_FACTOR,
	}
	node := &Node{
		keys:   []uint32{10, 20},
		values: []Addable{Float(1), Float(2), Float(3)},
		tree:   tree,
	}

	parent := &Node{
		keys:     []uint32{50},
		values:   []Addable{Float(0), Float(0)},
		children: []*Node{node, nil},
		tree:     tree,
	}

	node.parent = parent
	intervalTuple := ValueIntervalTuple{value: Addable(Float(7)), interval: Interval{start: 11, end: 13}}

	// Act
	node.insert(1, intervalTuple)

	// Assert
	assert.Equal(t, uint32(10), node.keys[0])
	assert.Equal(t, uint32(11), node.keys[1])
	assert.Equal(t, uint32(13), node.keys[2])
	assert.Equal(t, uint32(20), node.keys[3])
	assert.Equal(t, Float(1), node.values[0])
	assert.Equal(t, Float(2), node.values[1])
	assert.Equal(t, Float(9), node.values[2])
	assert.Equal(t, Float(2), node.values[3])
	assert.Equal(t, Float(3), node.values[4])
}

func TestInsertSelfContained3(t *testing.T) {
	// Arrange
	tree := &SegmentTreeImpl{
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: BRANCHING_FACTOR,
	}
	node := &Node{
		keys:   []uint32{10, 20},
		values: []Addable{Float(1), Float(2), Float(3)},
		tree:   tree,
	}

	parent := &Node{
		keys:     []uint32{50},
		values:   []Addable{Float(0), Float(0)},
		children: []*Node{node, nil},
		tree:     tree,
	}

	node.parent = parent
	intervalTuple := ValueIntervalTuple{value: Addable(Float(7)), interval: Interval{start: 21, end: 23}}

	// Act
	node.insert(2, intervalTuple)

	// Assert
	assert.Equal(t, uint32(10), node.keys[0])
	assert.Equal(t, uint32(20), node.keys[1])
	assert.Equal(t, uint32(21), node.keys[2])
	assert.Equal(t, uint32(23), node.keys[3])
	assert.Equal(t, Float(1), node.values[0])
	assert.Equal(t, Float(2), node.values[1])
	assert.Equal(t, Float(3), node.values[2])
	assert.Equal(t, Float(10), node.values[3])
	assert.Equal(t, Float(3), node.values[4])
}

func TestInsertNodeIntervalLeftLarger1(t *testing.T) {
	// Arrange
	tree := &SegmentTreeImpl{
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: BRANCHING_FACTOR,
	}
	node := &Node{
		keys:   []uint32{10, 20},
		values: []Addable{Float(1), Float(2), Float(3)},
		tree:   tree,
	}

	parent := &Node{
		keys:     []uint32{50},
		values:   []Addable{Float(0), Float(0)},
		children: []*Node{node, nil},
		tree:     tree,
	}

	node.parent = parent
	intervalTuple := ValueIntervalTuple{value: Addable(Float(7)), interval: Interval{start: 1, end: 100}}

	// Act
	node.insert(0, intervalTuple)

	// Assert
	assert.Equal(t, uint32(1), node.keys[0])
	assert.Equal(t, uint32(10), node.keys[1])
	assert.Equal(t, uint32(20), node.keys[2])
	assert.Equal(t, Float(1), node.values[0])
	assert.Equal(t, Float(8), node.values[1])
	assert.Equal(t, Float(2), node.values[2])
	assert.Equal(t, Float(3), node.values[3])
}

func TestInsertNodeIntervalLeftLarger2(t *testing.T) {
	// Arrange
	tree := &SegmentTreeImpl{
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: BRANCHING_FACTOR,
	}
	node := &Node{
		keys:   []uint32{10, 20},
		values: []Addable{Float(1), Float(2), Float(3)},
		tree:   tree,
	}

	parent := &Node{
		keys:     []uint32{50},
		values:   []Addable{Float(0), Float(0)},
		children: []*Node{node, nil},
		tree:     tree,
	}

	node.parent = parent
	intervalTuple := ValueIntervalTuple{value: Addable(Float(7)), interval: Interval{start: 11, end: 100}}

	// Act
	node.insert(1, intervalTuple)

	// Assert
	assert.Equal(t, uint32(10), node.keys[0])
	assert.Equal(t, uint32(11), node.keys[1])
	assert.Equal(t, uint32(20), node.keys[2])
	assert.Equal(t, Float(1), node.values[0])
	assert.Equal(t, Float(2), node.values[1])
	assert.Equal(t, Float(9), node.values[2])
	assert.Equal(t, Float(3), node.values[3])
}

func TestInsertNodeIntervalLeftLarger3(t *testing.T) {
	// Arrange
	tree := &SegmentTreeImpl{
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: BRANCHING_FACTOR,
	}
	node := &Node{
		keys:   []uint32{10, 20},
		values: []Addable{Float(1), Float(2), Float(3)},
		tree:   tree,
	}

	parent := &Node{
		keys:     []uint32{50},
		values:   []Addable{Float(0), Float(0)},
		children: []*Node{node, nil},
		tree:     tree,
	}

	node.parent = parent
	intervalTuple := ValueIntervalTuple{value: Addable(Float(7)), interval: Interval{start: 21, end: 100}}

	// Act
	node.insert(2, intervalTuple)

	// Assert
	assert.Equal(t, uint32(10), node.keys[0])
	assert.Equal(t, uint32(20), node.keys[1])
	assert.Equal(t, uint32(21), node.keys[2])
	assert.Equal(t, Float(1), node.values[0])
	assert.Equal(t, Float(2), node.values[1])
	assert.Equal(t, Float(3), node.values[2])
	assert.Equal(t, Float(10), node.values[3])
}

func TestInsertNodeIntervalRightLarger1(t *testing.T) {
	// Arrange
	tree := &SegmentTreeImpl{
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: BRANCHING_FACTOR,
	}
	node := &Node{
		keys:   []uint32{10, 20},
		values: []Addable{Float(1), Float(2), Float(3)},
		tree:   tree,
	}

	parent := &Node{
		keys:     []uint32{4},
		values:   []Addable{Float(0), Float(0)},
		children: []*Node{nil, node},
		tree:     tree,
	}

	node.parent = parent
	intervalTuple := ValueIntervalTuple{value: Addable(Float(7)), interval: Interval{start: 3, end: 6}}

	// Act
	node.insert(0, intervalTuple)

	// Assert
	assert.Equal(t, uint32(6), node.keys[0])
	assert.Equal(t, uint32(10), node.keys[1])
	assert.Equal(t, uint32(20), node.keys[2])
	assert.Equal(t, Float(8), node.values[0])
	assert.Equal(t, Float(1), node.values[1])
	assert.Equal(t, Float(2), node.values[2])
	assert.Equal(t, Float(3), node.values[3])
}

func TestInsertNodeIntervalRightLarger2(t *testing.T) {
	// Arrange
	tree := &SegmentTreeImpl{
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: BRANCHING_FACTOR,
	}
	node := &Node{
		keys:   []uint32{10, 20},
		values: []Addable{Float(1), Float(2), Float(3)},
		tree:   tree,
	}

	parent := &Node{
		keys:     []uint32{50},
		values:   []Addable{Float(0), Float(0)},
		children: []*Node{node, nil},
		tree:     tree,
	}

	node.parent = parent
	intervalTuple := ValueIntervalTuple{value: Addable(Float(7)), interval: Interval{start: 0, end: 11}}

	// Act
	node.insert(1, intervalTuple)

	// Assert
	assert.Equal(t, uint32(10), node.keys[0])
	assert.Equal(t, uint32(11), node.keys[1])
	assert.Equal(t, uint32(20), node.keys[2])
	assert.Equal(t, Float(1), node.values[0])
	assert.Equal(t, Float(9), node.values[1])
	assert.Equal(t, Float(2), node.values[2])
	assert.Equal(t, Float(3), node.values[3])
}

func TestInsertNodeIntervalRightLarger3(t *testing.T) {
	// Arrange
	tree := &SegmentTreeImpl{
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: BRANCHING_FACTOR,
	}
	node := &Node{
		keys:   []uint32{10, 20},
		values: []Addable{Float(1), Float(2), Float(3)},
		tree:   tree,
	}

	parent := &Node{
		keys:     []uint32{50},
		values:   []Addable{Float(0), Float(0)},
		children: []*Node{node, nil},
		tree:     tree,
	}

	node.parent = parent
	intervalTuple := ValueIntervalTuple{value: Addable(Float(7)), interval: Interval{start: 0, end: 21}}

	// Act
	node.insert(2, intervalTuple)

	// Assert
	assert.Equal(t, uint32(10), node.keys[0])
	assert.Equal(t, uint32(20), node.keys[1])
	assert.Equal(t, uint32(21), node.keys[2])
	assert.Equal(t, Float(1), node.values[0])
	assert.Equal(t, Float(2), node.values[1])
	assert.Equal(t, Float(10), node.values[2])
	assert.Equal(t, Float(3), node.values[3])
}

func TestInsertMatchingStartPoint1(t *testing.T) {
	// Arrange
	node := &Node{
		keys:   []uint32{10, 40},
		values: []Addable{Float(0), Float(2), Float(0)},
		tree: &SegmentTreeImpl{
			aggregate:       Aggregate{Sum, Identity, Float(0)},
			branchingFactor: BRANCHING_FACTOR,
		},
	}
	intervalTuple := ValueIntervalTuple{value: Float(3), interval: Interval{start: 0, end: 5}}

	// Act
	node.insert(0, intervalTuple)

	// Assert
	assert.Equal(t, uint32(5), node.keys[0])
	assert.Equal(t, uint32(10), node.keys[1])
	assert.Equal(t, uint32(40), node.keys[2])
	assert.Equal(t, Float(3), node.values[0])
	assert.Equal(t, Float(0), node.values[1])
	assert.Equal(t, Float(2), node.values[2])
	assert.Equal(t, Float(0), node.values[3])
}

func TestInsertMatchingStartPoint2(t *testing.T) {
	// Arrange
	node := &Node{
		keys:   []uint32{10, 40},
		values: []Addable{Float(0), Float(2), Float(0)},
		tree: &SegmentTreeImpl{
			aggregate:       Aggregate{Sum, Identity, Float(0)},
			branchingFactor: BRANCHING_FACTOR,
		},
	}
	intervalTuple := ValueIntervalTuple{value: Float(3), interval: Interval{start: 10, end: 30}}

	// Act
	node.insert(1, intervalTuple)

	// Assert
	assert.Equal(t, uint32(10), node.keys[0])
	assert.Equal(t, uint32(30), node.keys[1])
	assert.Equal(t, uint32(40), node.keys[2])
	assert.Equal(t, Float(0), node.values[0])
	assert.Equal(t, Float(5), node.values[1])
	assert.Equal(t, Float(2), node.values[2])
	assert.Equal(t, Float(0), node.values[3])
}

func TestInsertMatchingStartPoint3(t *testing.T) {
	// Arrange
	node := &Node{
		keys:   []uint32{10, 40},
		values: []Addable{Float(0), Float(2), Float(0)},
		tree: &SegmentTreeImpl{
			aggregate:       Aggregate{Sum, Identity, Float(0)},
			branchingFactor: BRANCHING_FACTOR,
		},
	}
	intervalTuple := ValueIntervalTuple{value: Float(3), interval: Interval{start: 40, end: 100}}

	// Act
	node.insert(2, intervalTuple)

	// Assert
	assert.Equal(t, uint32(10), node.keys[0])
	assert.Equal(t, uint32(40), node.keys[1])
	assert.Equal(t, uint32(100), node.keys[2])
	assert.Equal(t, Float(0), node.values[0])
	assert.Equal(t, Float(2), node.values[1])
	assert.Equal(t, Float(3), node.values[2])
	assert.Equal(t, Float(0), node.values[3])
}

func TestInsertMatchingEndPoint1(t *testing.T) {
	// Arrange
	node := &Node{
		keys:   []uint32{10, 40},
		values: []Addable{Float(0), Float(2), Float(0)},
		tree: &SegmentTreeImpl{
			aggregate:       Aggregate{Sum, Identity, Float(0)},
			branchingFactor: BRANCHING_FACTOR,
		},
	}
	intervalTuple := ValueIntervalTuple{value: Float(3), interval: Interval{start: 3, end: 10}}

	// Act
	node.insert(0, intervalTuple)

	// Assert
	assert.Equal(t, uint32(3), node.keys[0])
	assert.Equal(t, uint32(10), node.keys[1])
	assert.Equal(t, uint32(40), node.keys[2])
	assert.Equal(t, Float(0), node.values[0])
	assert.Equal(t, Float(3), node.values[1])
	assert.Equal(t, Float(2), node.values[2])
	assert.Equal(t, Float(0), node.values[3])
}

func TestInsertMatchingEndPoint2(t *testing.T) {
	// Arrange
	node := &Node{
		keys:   []uint32{10, 40},
		values: []Addable{Float(0), Float(2), Float(0)},
		tree: &SegmentTreeImpl{
			aggregate:       Aggregate{Sum, Identity, Float(0)},
			branchingFactor: BRANCHING_FACTOR,
		},
	}
	intervalTuple := ValueIntervalTuple{value: Float(3), interval: Interval{start: 30, end: 40}}

	// Act
	node.insert(1, intervalTuple)

	// Assert
	assert.Equal(t, uint32(10), node.keys[0])
	assert.Equal(t, uint32(30), node.keys[1])
	assert.Equal(t, uint32(40), node.keys[2])
	assert.Equal(t, Float(0), node.values[0])
	assert.Equal(t, Float(2), node.values[1])
	assert.Equal(t, Float(5), node.values[2])
	assert.Equal(t, Float(0), node.values[3])
}

func TestInsertMatchingEndPoint3(t *testing.T) {
	// Arrange
	node := &Node{
		keys:   []uint32{10, 40},
		values: []Addable{Float(0), Float(2), Float(0)},
		tree: &SegmentTreeImpl{
			aggregate:       Aggregate{Sum, Identity, Float(0)},
			branchingFactor: BRANCHING_FACTOR,
		},
	}
	intervalTuple := ValueIntervalTuple{value: Float(3), interval: Interval{start: 50, end: math.MaxUint32}}

	// Act
	node.insert(2, intervalTuple)

	// Assert
	assert.Equal(t, uint32(10), node.keys[0])
	assert.Equal(t, uint32(40), node.keys[1])
	assert.Equal(t, uint32(50), node.keys[2])
	assert.Equal(t, Float(0), node.values[0])
	assert.Equal(t, Float(2), node.values[1])
	assert.Equal(t, Float(0), node.values[2])
	assert.Equal(t, Float(3), node.values[3])
}

func SetupNodes() (*Node, *Node, *Node, *Node, *Node) {
	n0 := &Node{
		keys:     []uint32{15, 30, 45},
		children: make([]*Node, 5),
	}

	n1 := &Node{
		keys:   []uint32{5, 10},
		parent: n0,
	}

	n2 := &Node{
		keys:   []uint32{20},
		parent: n0,
	}

	n3 := &Node{
		keys:   []uint32{35, 40},
		parent: n0,
	}

	n4 := &Node{
		keys:   []uint32{50},
		parent: n0,
	}

	n0.children = []*Node{n1, n2, n3, n4}
	return n0, n1, n2, n3, n4
}
