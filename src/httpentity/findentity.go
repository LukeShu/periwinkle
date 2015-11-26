// Copyright 2015 Luke Shumaker

package httpentity

import "strings"

func pathSplit(upath string) (string, string) {
	parts := strings.SplitN(upath, "/", 2)
	switch len(parts) {
	case 1:
		return parts[0], ""
	case 2:
		return parts[0], parts[1]
	}
	panic("not reached")
}

// Takes the normalized path without the leading slash, and with or
// without a trailing slash.
func (router *Router) findEntity(upath string, request Request) (Entity, *Response) {
	var entity Entity = router.Root
	var handle404 func(string, Request) Response
	for {
		if g, ok := entity.(EntityGroup); ok {
			handle404 = g.SubentityNotFound
		}
		name_child, name_grandchildren := pathSplit(upath)
		entity := entity.Subentity(name_child, request)

		if entity == nil {
			response := handle404(name_child, request)
			return nil, &response
		} else if name_grandchildren == "" {
			return entity, nil
		} else {
			upath = name_grandchildren
		}
	}
}
