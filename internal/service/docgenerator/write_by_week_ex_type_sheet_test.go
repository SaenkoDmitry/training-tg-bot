package docgenerator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_weeksSort(t *testing.T) {
	strs := []string{
		"05.01.26 – 11.01.26",
		"12.01.26 – 18.01.26",
		"19.01.26 – 25.01.26",
		"26.01.26 – 01.02.26",
		"29.12.25 – 04.01.26",
	}
	strs = weeksSort(strs)
	assert.Equal(t, strs, []string{
		"29.12.25 – 04.01.26",
		"05.01.26 – 11.01.26",
		"12.01.26 – 18.01.26",
		"19.01.26 – 25.01.26",
		"26.01.26 – 01.02.26",
	})
}
