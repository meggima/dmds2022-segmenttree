package segmenttree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	BRANCHING_FACTOR uint32 = 4
)

var testDataGetAtInstant = []struct {
	instant       uint32
	expectedValue Float
}{
	{19, Float(6)},
	{49, Float(1)},
	{28, Float(7)},
	{0, Float(0)},
	{7, Float(2)},
}

// Yang et. al 2003, 3.1
// Lookup 19
func TestGetAtInstant(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	tree := setupTree()

	for _, testData := range testDataGetAtInstant {
		// Act
		res := tree.GetAtInstant(testData.instant)

		// Assert
		assert.Equal(testData.expectedValue, res)
	}
}

// Yang et. al 2003, 3.2
// Range [14, 28)
func TestGetWithinInterval(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	tree := setupTree()

	// Act
	res := tree.GetWithinInterval(Interval{start: 14, end: 28})

	// Assert
	assert.Len(res, 3)
	assert.Contains(res, ValueIntervalTuple{value: Float(8), interval: Interval{start: 14, end: 15}})
	assert.Contains(res, ValueIntervalTuple{value: Float(6), interval: Interval{start: 15, end: 20}})
	assert.Contains(res, ValueIntervalTuple{value: Float(7), interval: Interval{start: 20, end: 28}})
}

func TestNewTree(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	// Act
	tree := NewSegmentTree(BRANCHING_FACTOR, Aggregate{Sum, Identity, Float(0)})

	// Assert
	n0 := tree.root

	assert.Equal(BRANCHING_FACTOR, tree.branchingFactor)
	assert.Equal(Float(0), n0.values[0])
}

// Yang et. al 2003, 3.3
// Insert 1, [17, 47)
func TestInsert1(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	tree := setupTree()

	// Act
	tree.Insert(ValueIntervalTuple{value: Float(1), interval: Interval{start: 17, end: 47}})

	// Assert
	n0 := tree.root
	n2 := n0.children[1]
	n3 := n0.children[2]
	n4 := n0.children[3]

	assert.Equal(Float(0), n0.values[0])
	assert.Equal(Float(1), n0.values[1])
	assert.Equal(Float(1), n0.values[2])
	assert.Equal(Float(0), n0.values[3])

	assert.Equal(uint32(2), n2.size())
	assert.Equal(uint32(17), n2.keys[0])
	assert.Equal(uint32(20), n2.keys[1])
	assert.Equal(Float(5), n2.values[0])
	assert.Equal(Float(6), n2.values[1])
	assert.Equal(Float(7), n2.values[2])

	assert.Equal(uint32(2), n3.size())
	assert.Equal(uint32(35), n3.keys[0])
	assert.Equal(uint32(40), n3.keys[1])
	assert.Equal(Float(4), n3.values[0])
	assert.Equal(Float(8), n3.values[1])
	assert.Equal(Float(5), n3.values[2])

	assert.Equal(uint32(2), n4.size())
	assert.Equal(uint32(47), n4.keys[0])
	assert.Equal(uint32(50), n4.keys[1])
	assert.Equal(Float(2), n4.values[0])
	assert.Equal(Float(1), n4.values[1])
	assert.Equal(Float(0), n4.values[2])
}

// Yang et. al 2003, 3.3
// Insert 1, [24, 30)
func TestInsert2(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	tree := setupTree()

	// Act
	tree.Insert(ValueIntervalTuple{value: Float(1), interval: Interval{start: 24, end: 30}})

	// Assert
	n0 := tree.root
	n2 := n0.children[1]

	assert.Equal(uint32(2), n2.size())

	assert.Equal(uint32(20), n2.keys[0])
	assert.Equal(uint32(24), n2.keys[1])

	assert.Equal(Float(5), n2.values[0])
	assert.Equal(Float(6), n2.values[1])
	assert.Equal(Float(7), n2.values[2])
}

// Yang et. al 2003, 3.3
// Insert 1, [24, 28)
func TestInsert3(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	tree := setupTree()

	// Act
	tree.Insert(ValueIntervalTuple{value: Float(1), interval: Interval{start: 24, end: 28}})

	// Assert
	n0 := tree.root
	n2 := n0.children[1]

	assert.Equal(uint32(3), n2.size())

	assert.Equal(uint32(20), n2.keys[0])
	assert.Equal(uint32(24), n2.keys[1])
	assert.Equal(uint32(28), n2.keys[2])

	assert.Equal(Float(5), n2.values[0])
	assert.Equal(Float(6), n2.values[1])
	assert.Equal(Float(7), n2.values[2])
	assert.Equal(Float(6), n2.values[3])
}

// Yang et. al 2003, fig. 9
// Insert 1, [7, 12) & split
func TestInsert4(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	tree := setupTree()

	// Act
	tree.Insert(ValueIntervalTuple{value: Float(1), interval: Interval{start: 7, end: 12}})

	// Assert
	n0 := tree.root
	n01 := n0.children[0]
	n02 := n0.children[1]
	n11 := n01.children[0]
	n12 := n01.children[1]
	n2 := n01.children[2]
	n3 := n02.children[0]
	n4 := n02.children[1]

	assert.Equal(uint32(1), n0.size())
	assert.Equal(uint32(30), n0.keys[0])
	assert.Equal(Float(0), n0.values[0])
	assert.Equal(Float(0), n0.values[1])

	assert.Equal(uint32(2), n01.size())
	assert.Equal(uint32(10), n01.keys[0])
	assert.Equal(uint32(15), n01.keys[1])
	assert.Equal(Float(0), n01.values[0])
	assert.Equal(Float(0), n01.values[1])
	assert.Equal(Float(1), n01.values[2])

	assert.Equal(uint32(1), n02.size())
	assert.Equal(uint32(45), n02.keys[0])
	assert.Equal(Float(0), n02.values[0])
	assert.Equal(Float(0), n02.values[1])

	assert.Equal(uint32(2), n11.size())
	assert.Equal(uint32(5), n11.keys[0])
	assert.Equal(uint32(7), n11.keys[1])
	assert.Equal(Float(0), n11.values[0])
	assert.Equal(Float(2), n11.values[1])
	assert.Equal(Float(3), n11.values[2])

	assert.Equal(uint32(1), n12.size())
	assert.Equal(uint32(12), n12.keys[0])
	assert.Equal(Float(9), n12.values[0])
	assert.Equal(Float(8), n12.values[1])

	assert.Equal(uint32(1), n2.size())
	assert.Equal(uint32(20), n2.keys[0])
	assert.Equal(Float(5), n2.values[0])
	assert.Equal(Float(6), n2.values[1])

	assert.Equal(uint32(2), n3.size())
	assert.Equal(uint32(35), n3.keys[0])
	assert.Equal(uint32(40), n3.keys[1])
	assert.Equal(Float(4), n3.values[0])
	assert.Equal(Float(8), n3.values[1])
	assert.Equal(Float(5), n3.values[2])

	assert.Equal(uint32(1), n4.size())
	assert.Equal(uint32(50), n4.keys[0])
	assert.Equal(Float(1), n4.values[0])
	assert.Equal(Float(0), n4.values[1])
}

func TestInsertTwiceSameRangeSimpleElement(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	n0 := &Node{
		nodeId:   0,
		keys:     []uint32{},
		values:   []Addable{Float(0)},
		children: []*Node{},
		isLeaf:   true,
	}
	n0.parent = nil
	tree := &SegmentTreeImpl{
		root:            n0,
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: BRANCHING_FACTOR,
	}
	n0.tree = tree

	tree.Insert(ValueIntervalTuple{value: Float(2), interval: Interval{start: 10, end: 40}})
	tree.Insert(ValueIntervalTuple{value: Float(3), interval: Interval{start: 10, end: 40}})

	// Assert
	assert.Equal(2, int(n0.size()))
}

// Yang et. al 2003, 3.4 & 3.6
// Delete 1, [17, 47) & interval merge
func TestDelete1(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	tree := setupTree()
	tree.Insert(ValueIntervalTuple{value: Float(1), interval: Interval{start: 17, end: 47}})

	// Act
	tree.Delete(ValueIntervalTuple{value: Float(1), interval: Interval{start: 17, end: 47}})

	// Assert
	n0 := tree.root
	n1 := n0.children[0]
	n2 := n0.children[1]
	n3 := n0.children[2]
	n4 := n0.children[3]

	assert.Equal(uint32(2), n1.size())
	assert.Equal(uint32(5), n1.keys[0])
	assert.Equal(uint32(10), n1.keys[1])
	assert.Equal(Float(0), n1.values[0])
	assert.Equal(Float(2), n1.values[1])
	assert.Equal(Float(8), n1.values[2])

	assert.Equal(uint32(1), n2.size())
	assert.Equal(uint32(20), n2.keys[0])
	assert.Equal(Float(5), n2.values[0])
	assert.Equal(Float(6), n2.values[1])

	assert.Equal(uint32(2), n3.size())
	assert.Equal(uint32(35), n3.keys[0])
	assert.Equal(uint32(40), n3.keys[1])
	assert.Equal(Float(4), n3.values[0])
	assert.Equal(Float(8), n3.values[1])
	assert.Equal(Float(5), n3.values[2])

	assert.Equal(uint32(1), n4.size())
	assert.Equal(uint32(50), n4.keys[0])
	assert.Equal(Float(1), n4.values[0])
	assert.Equal(Float(0), n4.values[1])
}

func TestDelete2(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	tree := setupTree()
	tree.Insert(ValueIntervalTuple{value: Float(1), interval: Interval{start: 7, end: 12}})

	// Act
	tree.Delete(ValueIntervalTuple{value: Float(1), interval: Interval{start: 7, end: 12}})

	// Assert
	n0 := tree.root
	n01 := n0.children[0]
	n02 := n0.children[1]
	n11 := n01.children[0]
	n2 := n01.children[1]
	n3 := n02.children[0]
	n4 := n02.children[1]

	assert.Equal(uint32(1), n0.size())
	assert.Equal(uint32(30), n0.keys[0])
	assert.Equal(Float(0), n0.values[0])
	assert.Equal(Float(0), n0.values[1])

	assert.Equal(uint32(1), n01.size())
	assert.Equal(uint32(10), n01.keys[0])
	assert.Equal(Float(0), n01.values[0])
	assert.Equal(Float(0), n01.values[1])

	assert.Equal(uint32(1), n02.size())
	assert.Equal(uint32(45), n02.keys[0])
	assert.Equal(Float(0), n02.values[0])
	assert.Equal(Float(0), n02.values[1])

	assert.Equal(uint32(1), n11.size())
	assert.Equal(uint32(5), n11.keys[0])
	assert.Equal(Float(0), n11.values[0])
	assert.Equal(Float(2), n11.values[1])

	assert.Equal(uint32(2), n2.size())
	assert.Equal(uint32(15), n2.keys[0])
	assert.Equal(uint32(20), n2.keys[1])
	assert.Equal(Float(8), n2.values[0])
	assert.Equal(Float(6), n2.values[1])
	assert.Equal(Float(7), n2.values[2])

	assert.Equal(uint32(2), n3.size())
	assert.Equal(uint32(35), n3.keys[0])
	assert.Equal(uint32(40), n3.keys[1])
	assert.Equal(Float(4), n3.values[0])
	assert.Equal(Float(8), n3.values[1])
	assert.Equal(Float(5), n3.values[2])

	assert.Equal(uint32(1), n4.size())
	assert.Equal(uint32(50), n4.keys[0])
	assert.Equal(Float(1), n4.values[0])
	assert.Equal(Float(0), n4.values[1])
}

// Yang et. al 2003, Fig 4
func setupTree() *SegmentTreeImpl {
	n1 := &Node{
		nodeId:   1,
		keys:     []uint32{5, 10},
		values:   []Addable{Float(0), Float(2), Float(8)},
		children: []*Node{},
		isLeaf:   true,
	}

	n2 := &Node{
		nodeId:   2,
		keys:     []uint32{20},
		values:   []Addable{Float(5), Float(6)},
		children: []*Node{},
		isLeaf:   true,
	}

	n3 := &Node{
		nodeId:   3,
		keys:     []uint32{35, 40},
		values:   []Addable{Float(4), Float(8), Float(5)},
		children: []*Node{},
		isLeaf:   true,
	}

	n4 := &Node{
		nodeId:   4,
		keys:     []uint32{50},
		values:   []Addable{Float(1), Float(0)},
		children: []*Node{},
		isLeaf:   true,
	}

	n0 := &Node{
		nodeId:   0,
		keys:     []uint32{15, 30, 45},
		values:   []Addable{Float(0), Float(1), Float(0), Float(0)},
		children: []*Node{n1, n2, n3, n4},
		isLeaf:   false,
	}

	n0.parent = nil
	n1.parent = n0
	n2.parent = n0
	n3.parent = n0
	n4.parent = n0

	tree := &SegmentTreeImpl{
		root:            n0,
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: BRANCHING_FACTOR,
	}

	n0.tree = tree
	n1.tree = tree
	n2.tree = tree
	n3.tree = tree
	n4.tree = tree

	return tree
}

// Yang et. al 2003, Fig 19
func TestSumDosageScenarioInsert(t *testing.T) {

	// Arrange
	assert := assert.New(t)

	n0 := &Node{
		nodeId:   0,
		keys:     []uint32{},
		values:   []Addable{Float(0)},
		children: []*Node{},
		isLeaf:   true,
	}

	n0.parent = nil
	tree := &SegmentTreeImpl{
		root:            n0,
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: 4,
	}

	n0.tree = tree

	// Act
	tree.Insert(ValueIntervalTuple{value: Float(2), interval: Interval{start: 10, end: 40}})
	tree.Insert(ValueIntervalTuple{value: Float(3), interval: Interval{start: 10, end: 30}})
	tree.Insert(ValueIntervalTuple{value: Float(1), interval: Interval{start: 20, end: 40}})
	// split nodes
	tree.Insert(ValueIntervalTuple{value: Float(2), interval: Interval{start: 5, end: 15}})
	// split nodes
	tree.Insert(ValueIntervalTuple{value: Float(4), interval: Interval{start: 35, end: 45}})
	tree.Insert(ValueIntervalTuple{value: Float(1), interval: Interval{start: 10, end: 50}})
	// split nodes

	// Assert
	n0 = tree.root
	assert.Len(n0.children, 4)
	n00 := n0.children[0]
	n01 := n0.children[1]
	n02 := n0.children[2]
	n03 := n0.children[3]

	assert.Equal(uint32(15), n0.keys[0])
	assert.Equal(uint32(30), n0.keys[1])
	assert.Equal(uint32(45), n0.keys[2])
	assert.Len(n0.keys, 3)
	assert.Equal(Float(0), n0.values[0])
	assert.Equal(Float(1), n0.values[1])
	assert.Equal(Float(0), n0.values[2])
	assert.Equal(Float(0), n0.values[3])
	assert.Len(n0.values, 4)

	assert.Equal(uint32(5), n00.keys[0])
	assert.Equal(uint32(10), n00.keys[1])
	assert.Len(n00.keys, 2)
	assert.Equal(Float(0), n00.values[0])
	assert.Equal(Float(2), n00.values[1])
	assert.Equal(Float(8), n00.values[2])
	assert.Len(n00.values, 3)
	assert.Len(n00.children, 0)

	assert.Equal(uint32(20), n01.keys[0])
	assert.Len(n01.keys, 1)
	assert.Equal(Float(5), n01.values[0])
	assert.Equal(Float(6), n01.values[1])
	assert.Len(n01.values, 2)
	assert.Len(n01.children, 0)

	assert.Equal(uint32(35), n02.keys[0])
	assert.Equal(uint32(40), n02.keys[1])
	assert.Len(n02.keys, 2)
	assert.Equal(Float(4), n02.values[0])
	assert.Equal(Float(8), n02.values[1])
	assert.Equal(Float(5), n02.values[2])
	assert.Len(n02.values, 3)
	assert.Len(n02.children, 0)

	assert.Equal(uint32(50), n03.keys[0])
	assert.Len(n03.keys, 1)
	assert.Equal(Float(1), n03.values[0])
	assert.Equal(Float(0), n03.values[1])
	assert.Len(n03.values, 2)
	assert.Len(n03.children, 0)
}

// Yang et. al 2003, Fig 19
func TestSumDosageScenarioDelete(t *testing.T) {

	// Arrange
	assert := assert.New(t)

	n0 := &Node{
		nodeId:   0,
		keys:     []uint32{},
		values:   []Addable{},
		children: []*Node{},
		isLeaf:   false,
	}
	n0.parent = nil
	tree := &SegmentTreeImpl{
		root:            n0,
		aggregate:       Aggregate{Sum, Identity, Float(0)},
		branchingFactor: BRANCHING_FACTOR,
	}
	n0.tree = tree

	tree.Insert(ValueIntervalTuple{value: Float(2), interval: Interval{start: 10, end: 40}})
	tree.Insert(ValueIntervalTuple{value: Float(3), interval: Interval{start: 10, end: 30}})
	tree.Insert(ValueIntervalTuple{value: Float(1), interval: Interval{start: 20, end: 40}})
	tree.Insert(ValueIntervalTuple{value: Float(2), interval: Interval{start: 5, end: 15}})
	tree.Insert(ValueIntervalTuple{value: Float(4), interval: Interval{start: 35, end: 45}})
	tree.Insert(ValueIntervalTuple{value: Float(1), interval: Interval{start: 10, end: 50}})

	// Act
	tree.Delete(ValueIntervalTuple{value: Float(1), interval: Interval{start: 10, end: 50}})
	tree.Delete(ValueIntervalTuple{value: Float(4), interval: Interval{start: 35, end: 45}})
	// merge and remove node
	tree.Delete(ValueIntervalTuple{value: Float(2), interval: Interval{start: 5, end: 15}})
	tree.Delete(ValueIntervalTuple{value: Float(1), interval: Interval{start: 20, end: 40}})
	// merge and remove node
	tree.Delete(ValueIntervalTuple{value: Float(3), interval: Interval{start: 10, end: 30}})
	// merge and remove node
	tree.Delete(ValueIntervalTuple{value: Float(2), interval: Interval{start: 10, end: 40}})
	// empty tree

	// Assert
	n0 = tree.root

	assert.Len(n0.keys, 0)
	assert.Equal(Float(0), n0.values[0])
	assert.Len(n0.values, 0)
	assert.Len(n0.children, 0)
}
