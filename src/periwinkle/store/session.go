// Copyright 2015 Davis Webb
// Copyright 2015 Luke Shumaker

package store

type Session struct {
	id        int
	user_id   int
	last_used string // not 100% on this
}
