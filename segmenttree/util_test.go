package segmenttree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDataMinUint32 = []struct {
	a        uint32
	b        uint32
	expected uint32
}{
	{10, 20, 10},
	{20, 10, 10},
	{10, 10, 10},
}

var testDataMaxUint32 = []struct {
	a        uint32
	b        uint32
	expected uint32
}{
	{10, 20, 20},
	{20, 10, 20},
	{10, 10, 10},
}

func TestMinUint32(t *testing.T) {
	for _, testData := range testDataMinUint32 {
		// Arrange
		assert := assert.New(t)

		// Act
		res := MinUint32(testData.a, testData.b)

		// Assert
		assert.Equal(testData.expected, res)
	}
}

func TestMaxUint32(t *testing.T) {
	for _, testData := range testDataMaxUint32 {
		// Arrange
		assert := assert.New(t)

		// Act
		res := MaxUint32(testData.a, testData.b)

		// Assert
		assert.Equal(testData.expected, res)
	}
}
