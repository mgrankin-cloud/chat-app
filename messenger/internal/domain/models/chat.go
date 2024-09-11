package models

type Chat struct {
	ID       int64
	Name     string
	Photo    []byte
	ChatType int
	Status string
}
