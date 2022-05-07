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

func TestIntersectionWith(t *testing.T) {
	// Arrange
	testData := []struct {
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

	for _, td := range testData {
		// Act
		res := td.a.IntersectionWith(td.b)

		// Assert
		assert.Equal(t, td.expected, res)
	}
}

func TestIsSubsetOf(t *testing.T) {
	// Arrange
	testData := []struct {
		a        Interval
		b        Interval
		isSubset bool
	}{
		{NewInterval(2, 3), NewInterval(2, 3), true},
		{NewInterval(1, 3), NewInterval(2, 4), false},
		{NewInterval(3, 5), NewInterval(2, 4), false},
		{NewInterval(3, 5), NewInterval(1, 10), true},
		{NewInterval(1, 5), NewInterval(2, 4), false},
	}

	for _, td := range testData {
		// Act
		res := td.a.IsSubsetOf(td.b)

		// Assert
		assert.Equal(t, td.isSubset, res)
	}
}
