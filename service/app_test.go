package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp_OpenRedis(t *testing.T) {
	app := GetApp()
	app.OpenRedis()
}

func TestApp_OpenKafkaProducer(t *testing.T) {
	app := GetApp()
	app.OpenKafkaProducer()
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
