package service

import "testing"

func TestApp_OpenDB(t *testing.T) {
	app := GetApp()
	app.OpenDB()
}

func TestApp_OpenRedis(t *testing.T) {
	app := GetApp()
	app.OpenRedis()
}
