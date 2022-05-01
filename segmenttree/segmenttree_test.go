package segmenttree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	BRANCHING_FACTOR uint32 = 4
)

// Yang et. al 2003, 3.1
// Lookup 19
func TestGetAtInstant(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	tree := setupTree()

	// Act
	res := tree.GetAtInstant(19)

	// Assert
	assert.Equal(float32(6), res)
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
	tree := NewSegmentTree(BRANCHING_FACTOR, Sum)

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
	n4 := n0.children[3]

	assert.Equal(Float(1), n0.values[2])

	assert.Equal(uint32(2), n2.n)
	assert.Equal(uint32(17), n2.keys[0])
	assert.Equal(uint32(20), n2.keys[1])
	assert.Equal(Float(5), n2.values[0])
	assert.Equal(Float(6), n2.values[1])
	assert.Equal(Float(7), n2.values[2])

	assert.Equal(uint32(2), n4.n)
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

	assert.Equal(uint32(2), n2.n)

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

	assert.Equal(uint32(3), n2.n)

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

	assert.Equal(uint32(1), n0.n)
	assert.Equal(uint32(30), n0.keys[0])
	assert.Equal(Float(0), n0.values[0])
	assert.Equal(Float(0), n0.values[1])

	assert.Equal(uint32(2), n01.n)
	assert.Equal(uint32(10), n01.keys[0])
	assert.Equal(uint32(15), n01.keys[1])
	assert.Equal(Float(0), n01.values[0])
	assert.Equal(Float(0), n01.values[1])
	assert.Equal(Float(1), n01.values[2])

	assert.Equal(uint32(1), n02.n)
	assert.Equal(uint32(45), n02.keys[0])
	assert.Equal(Float(0), n02.values[0])
	assert.Equal(Float(0), n02.values[1])

	assert.Equal(uint32(2), n11.n)
	assert.Equal(uint32(5), n11.keys[0])
	assert.Equal(uint32(7), n11.keys[1])
	assert.Equal(Float(0), n11.values[0])
	assert.Equal(Float(2), n11.values[1])
	assert.Equal(Float(3), n11.values[2])

	assert.Equal(uint32(1), n12.n)
	assert.Equal(uint32(12), n12.keys[0])
	assert.Equal(Float(9), n12.values[0])
	assert.Equal(Float(8), n12.values[1])

	assert.Equal(uint32(1), n2.n)
	assert.Equal(uint32(20), n2.keys[0])
	assert.Equal(Float(5), n2.values[0])
	assert.Equal(Float(6), n2.values[1])

	assert.Equal(uint32(2), n3.n)
	assert.Equal(uint32(35), n3.keys[0])
	assert.Equal(uint32(40), n3.keys[1])
	assert.Equal(Float(4), n3.values[0])
	assert.Equal(Float(8), n3.values[1])
	assert.Equal(Float(5), n3.values[2])

	assert.Equal(uint32(1), n4.n)
	assert.Equal(uint32(50), n4.keys[0])
	assert.Equal(Float(1), n4.values[0])
	assert.Equal(Float(0), n4.values[1])
}

// Yang et. al 2003, 3.4 & 3.6
// Delete 1, [17, 47) & imerge
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

	assert.Equal(uint32(2), n1.n)
	assert.Equal(uint32(5), n1.keys[0])
	assert.Equal(uint32(10), n1.keys[1])
	assert.Equal(Float(0), n1.values[0])
	assert.Equal(Float(2), n1.values[1])
	assert.Equal(Float(8), n1.values[2])

	assert.Equal(uint32(1), n2.n)
	assert.Equal(uint32(20), n2.keys[0])
	assert.Equal(Float(5), n2.values[0])
	assert.Equal(Float(6), n2.values[1])

	assert.Equal(uint32(2), n3.n)
	assert.Equal(uint32(35), n3.keys[0])
	assert.Equal(uint32(40), n3.keys[1])
	assert.Equal(Float(4), n3.values[0])
	assert.Equal(Float(8), n3.values[1])
	assert.Equal(Float(5), n3.values[2])

	assert.Equal(uint32(1), n4.n)
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

	assert.Equal(uint32(1), n0.n)
	assert.Equal(uint32(30), n0.keys[0])
	assert.Equal(Float(0), n0.values[0])
	assert.Equal(Float(0), n0.values[1])

	assert.Equal(uint32(1), n01.n)
	assert.Equal(uint32(10), n01.keys[0])
	assert.Equal(Float(0), n01.values[0])
	assert.Equal(Float(0), n01.values[1])

	assert.Equal(uint32(1), n02.n)
	assert.Equal(uint32(45), n02.keys[0])
	assert.Equal(Float(0), n02.values[0])
	assert.Equal(Float(0), n02.values[1])

	assert.Equal(uint32(1), n11.n)
	assert.Equal(uint32(5), n11.keys[0])
	assert.Equal(Float(0), n11.values[0])
	assert.Equal(Float(2), n11.values[1])

	assert.Equal(uint32(2), n2.n)
	assert.Equal(uint32(15), n2.keys[0])
	assert.Equal(uint32(20), n2.keys[1])
	assert.Equal(Float(8), n2.values[0])
	assert.Equal(Float(6), n2.values[1])
	assert.Equal(Float(7), n2.values[2])

	assert.Equal(uint32(2), n3.n)
	assert.Equal(uint32(35), n3.keys[0])
	assert.Equal(uint32(40), n3.keys[1])
	assert.Equal(Float(4), n3.values[0])
	assert.Equal(Float(8), n3.values[1])
	assert.Equal(Float(5), n3.values[2])

	assert.Equal(uint32(1), n4.n)
	assert.Equal(uint32(50), n4.keys[0])
	assert.Equal(Float(1), n4.values[0])
	assert.Equal(Float(0), n4.values[1])
}

// Yang et. al 2003, Fig 4
func setupTree() *SegmentTreeImpl {
	n1 := &Node{
		nodeId:   1,
		n:        2,
		keys:     []uint32{5, 10, 0},
		values:   []Float{Float(0), Float(2), Float(8), Float(-1)},
		children: []*Node{nil, nil, nil, nil},
	}

	n2 := &Node{
		nodeId:   2,
		n:        1,
		keys:     []uint32{20, 0, 0},
		values:   []Float{Float(5), Float(6), Float(-1), Float(-1)},
		children: []*Node{nil, nil, nil, nil},
	}

	n3 := &Node{
		nodeId:   3,
		n:        2,
		keys:     []uint32{35, 40, 0},
		values:   []Float{Float(4), Float(8), Float(5), Float(-1)},
		children: []*Node{nil, nil, nil, nil},
	}

	n4 := &Node{
		nodeId:   4,
		n:        1,
		keys:     []uint32{50, 0, 0},
		values:   []Float{Float(1), Float(0), Float(-1), Float(-1)},
		children: []*Node{nil, nil, nil, nil},
	}

	n0 := &Node{
		nodeId:   0,
		n:        3,
		keys:     []uint32{15, 30, 45},
		values:   []Float{Float(0), Float(1), Float(0), Float(0)},
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
		operation:       Sum,
		branchingFactor: BRANCHING_FACTOR,
	}

	n0.tree = tree
	n1.tree = tree
	n2.tree = tree
	n3.tree = tree
	n4.tree = tree

	return tree
}
