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
		n:    3,
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
		n:    3,
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
		n:    3,
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

func SetupNodes() (*Node, *Node, *Node, *Node, *Node) {
	n0 := &Node{
		n:        3,
		keys:     []uint32{15, 30, 45},
		children: make([]*Node, 5),
	}

	n1 := &Node{
		n:      2,
		keys:   []uint32{5, 10},
		parent: n0,
	}

	n2 := &Node{
		n:      1,
		keys:   []uint32{20},
		parent: n0,
	}

	n3 := &Node{
		n:      2,
		keys:   []uint32{35, 40},
		parent: n0,
	}

	n4 := &Node{
		n:      1,
		keys:   []uint32{50},
		parent: n0,
	}

	n0.children = []*Node{n1, n2, n3, n4}
	return n0, n1, n2, n3, n4
}
