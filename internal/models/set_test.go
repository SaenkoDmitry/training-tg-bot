package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_strikePlanned_int(t *testing.T) {
	assert.Equal(t, "<strike>15</strike> <b>0</b>", strikePlanned(15, 0, true))
	assert.Equal(t, "<strike>15</strike> <b>20</b>", strikePlanned(15, 20, true))
	assert.Equal(t, "20", strikePlanned(0, 20, true))

	assert.Equal(t, "15", strikePlanned(15, 0, false))
	assert.Equal(t, "15", strikePlanned(15, 20, false))
	assert.Equal(t, "0", strikePlanned(0, 20, false))
}

func Test_strikePlanned_float(t *testing.T) {
	assert.Equal(t, "<strike>15.5</strike> <b>0</b>", strikePlanned(15.5, 0.0, true))
	assert.Equal(t, "<strike>15.5</strike> <b>20.1</b>", strikePlanned(15.5, 20.1, true))
	assert.Equal(t, "20", strikePlanned(0, 20.0, true))

	assert.Equal(t, "15.5", strikePlanned(15.5, 0, false))
	assert.Equal(t, "15.5", strikePlanned(15.5, 20.1, false))
	assert.Equal(t, "0", strikePlanned(0.0, 20.1, false))
}
