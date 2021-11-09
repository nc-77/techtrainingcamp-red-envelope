package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp_OpenDB(t *testing.T) {
	app := GetApp()
	app.OpenDB()
}

func TestApp_OpenRedis(t *testing.T) {
	app := GetApp()
	app.OpenRedis()
}

func TestCheckSnatchedPr(t *testing.T) {
	value, ok := CheckSnatchedPr("80")
	assert.Equal(t, 80, value)
	assert.Equal(t, true, ok)

	_, ok = CheckSnatchedPr("101")
	assert.Equal(t, false, ok)

	_, ok = CheckSnatchedPr("-1")
	assert.Equal(t, false, ok)

	_, ok = CheckSnatchedPr("hello")
	assert.Equal(t, false, ok)
}
