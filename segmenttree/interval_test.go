package segmenttree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInterval(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	// Act
	interval := NewInterval(1, 2)

	// Assert
	assert.Equal(Interval{1, 2}, interval)
}

func TestNewIntervalInvalid(t *testing.T) {
	// Assert
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()

	// Act
	NewInterval(2, 1)
}

var testDataIntersectionWith = []struct {
	a        Interval
	b        Interval
	expected Interval
}{
	{NewInterval(1, 2), NewInterval(3, 4), EmptyInterval},
	{NewInterval(3, 4), NewInterval(1, 2), EmptyInterval},
	{NewInterval(1, 3), NewInterval(2, 4), NewInterval(2, 3)},
	{NewInterval(2, 4), NewInterval(1, 3), NewInterval(2, 3)},
	{NewInterval(1, 4), NewInterval(1, 4), NewInterval(1, 4)},
	{NewInterval(1, 4), NewInterval(2, 3), NewInterval(2, 3)},
	{NewInterval(1, 2), NewInterval(2, 3), NewInterval(2, 2)},
}

func TestIntersectionWith(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	for _, testData := range testDataIntersectionWith {
		// Act
		res := testData.a.IntersectionWith(testData.b)

		// Assert
		assert.Equal(testData.expected, res)
	}
}
