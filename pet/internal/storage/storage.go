package storage

import "errors"

var (
	ErrGoodsExists   = errors.New("goods already exists")
	ErrGoodsNotFound = errors.New("goods not found")
	Errmove          = errors.New("move is not success")
)
