package timeutil_test

import (
	"testing"

	timeutil "go-clean-architecture/utils/time"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeNow(t *testing.T) {
	value := timeutil.GetTimeNow()
	assert.Equal(t, value, value)
}
