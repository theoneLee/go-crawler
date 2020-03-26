package eumn

import "sync"

var (
	CurCount = 0
	MaxCount = 2
	RwMutex  = sync.RWMutex{}
)
