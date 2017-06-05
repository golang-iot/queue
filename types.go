package queue

import (
	"time"
)

type Message struct {
	Message string
	Value   int
	Address int
	Created time.Time
}

type FileChunk struct{
	Id string
	Name string
	Total int64
	Current int64
	Content []byte
	ChunkSize int64
}