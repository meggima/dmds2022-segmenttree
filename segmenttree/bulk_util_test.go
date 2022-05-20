package segmenttree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertSortedWhenSortedAscending(t *testing.T) {
	// Arrange
	aggregate := Aggregate{Sum, InverseSum, Identity, Float(0)}
	list := make([]ValueTimeTuple, 0, 10)

	// Act
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(1), time: 1}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(2), time: 2}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(3), time: 3}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(4), time: 4}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(5), time: 5}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(6), time: 6}, list)

	// Assert
	assert.Len(t, list, 6)
	assert.Equal(t, Float(1), list[0].value)
	assert.Equal(t, Float(2), list[1].value)
	assert.Equal(t, Float(3), list[2].value)
	assert.Equal(t, Float(4), list[3].value)
	assert.Equal(t, Float(5), list[4].value)
	assert.Equal(t, Float(6), list[5].value)
}

func TestInsertSortedWhenSortedDescending(t *testing.T) {
	// Arrange
	aggregate := Aggregate{Sum, InverseSum, Identity, Float(0)}
	list := make([]ValueTimeTuple, 0, 10)

	// Act
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(6), time: 6}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(5), time: 5}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(4), time: 4}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(3), time: 3}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(2), time: 2}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(1), time: 1}, list)

	// Assert
	assert.Len(t, list, 6)
	assert.Equal(t, Float(1), list[0].value)
	assert.Equal(t, Float(2), list[1].value)
	assert.Equal(t, Float(3), list[2].value)
	assert.Equal(t, Float(4), list[3].value)
	assert.Equal(t, Float(5), list[4].value)
	assert.Equal(t, Float(6), list[5].value)
}

func TestInsertSortedWhenRandomOrder(t *testing.T) {
	// Arrange
	aggregate := Aggregate{Sum, InverseSum, Identity, Float(0)}
	list := make([]ValueTimeTuple, 0, 6)

	// Act
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(6), time: 6}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(1), time: 1}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(4), time: 4}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(5), time: 5}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(3), time: 3}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(2), time: 2}, list)

	// Assert
	assert.Len(t, list, 6)
	assert.Equal(t, Float(1), list[0].value)
	assert.Equal(t, Float(2), list[1].value)
	assert.Equal(t, Float(3), list[2].value)
	assert.Equal(t, Float(4), list[3].value)
	assert.Equal(t, Float(5), list[4].value)
	assert.Equal(t, Float(6), list[5].value)
}

func TestInsertSortedWhenMerge(t *testing.T) {
	// Arrange
	aggregate := Aggregate{Sum, InverseSum, Identity, Float(0)}
	list := make([]ValueTimeTuple, 0, 6)

	// Act
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(1), time: 1}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(2), time: 2}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(3), time: 3}, list) // will be merged
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(4), time: 3}, list) // will be merged
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(5), time: 5}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(6), time: 6}, list)

	// Assert
	assert.Len(t, list, 5)
	assert.Equal(t, Float(1), list[0].value)
	assert.Equal(t, Float(2), list[1].value)
	assert.Equal(t, Float(7), list[2].value)
	assert.Equal(t, Float(5), list[3].value)
	assert.Equal(t, Float(6), list[4].value)
}

func TestInsertSortedWhenInverse(t *testing.T) {
	// Arrange
	aggregate := Aggregate{Sum, InverseSum, Identity, Float(0)}
	list := make([]ValueTimeTuple, 0, 6)

	// Act
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(1), time: 1}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(2), time: 2}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(3), time: 3}, list)  // will be removed
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(-3), time: 3}, list) // will be removed
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(5), time: 5}, list)
	list = insertInOrder(aggregate, ValueTimeTuple{value: Float(6), time: 6}, list)

	// Assert
	assert.Len(t, list, 4)
	assert.Equal(t, Float(1), list[0].value)
	assert.Equal(t, Float(2), list[1].value)
	assert.Equal(t, Float(5), list[2].value)
	assert.Equal(t, Float(6), list[3].value)
}

func TestSplitAndSort(t *testing.T) {
	// Arrange
	aggregate := Aggregate{Sum, InverseSum, Identity, Float(0)}

	var testData []ValueIntervalTuple = []ValueIntervalTuple{
		{interval: NewInterval(10, 40), value: Float(2)},
		{interval: NewInterval(10, 30), value: Float(3)},
		{interval: NewInterval(20, 40), value: Float(1)},
		{interval: NewInterval(5, 15), value: Float(2)},
		{interval: NewInterval(35, 45), value: Float(4)},
		{interval: NewInterval(10, 50), value: Float(1)},
	}

	// Act
	result := createAndSortValueTimeTuples(aggregate, testData)

	// Assert
	assert.Len(t, result, 9)
	assert.Equal(t, ValueTimeTuple{time: 5, value: Float(2)}, result[0])
	assert.Equal(t, ValueTimeTuple{time: 10, value: Float(6)}, result[1])
	assert.Equal(t, ValueTimeTuple{time: 15, value: Float(-2)}, result[2])
	assert.Equal(t, ValueTimeTuple{time: 20, value: Float(1)}, result[3])
	assert.Equal(t, ValueTimeTuple{time: 30, value: Float(-3)}, result[4])
	assert.Equal(t, ValueTimeTuple{time: 35, value: Float(4)}, result[5])
	assert.Equal(t, ValueTimeTuple{time: 40, value: Float(-3)}, result[6])
	assert.Equal(t, ValueTimeTuple{time: 45, value: Float(-4)}, result[7])
	assert.Equal(t, ValueTimeTuple{time: 50, value: Float(-1)}, result[8])
}
