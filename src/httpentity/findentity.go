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
func (router *Router) findEntity(upath string, request Request) (Entity, Response) {
	var group EntityGroup = router.Root
	for {
		nameChild, nameGrandchildren := pathSplit(upath)
		entity := group.Subentity(nameChild, request)

		if entity == nil {
			return nil, group.SubentityNotFound(nameChild, request)
		} else if nameGrandchildren == "" {
			return entity, Response{}
		} else {
			if newgroup, ok := entity.(EntityGroup); ok {
				group = newgroup
			} else {
				nameGrandchild, _ := pathSplit(upath)
				return nil, group.SubentityNotFound(nameChild+"/"+nameGrandchild, request)
			}
		}
	}
}
