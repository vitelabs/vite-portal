package generics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterDuplicates(t *testing.T) {
	orig := []string {"test1", "test2", "test1", "test1.2"}
	actual := FilterDuplicates(orig...)
	assert.Equal(t, 3, len(actual))
}