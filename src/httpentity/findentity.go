// Copyright 2015 Luke Shumaker

package httpentity

import (
	"fmt"
	"strings"
)

// Takes the normalized path without the leading slash
func findEntity(entity Entity, req Request, upath string) Entity {
	if entity == nil {
		return nil // 404
	} else if upath == "" {
		return entity
	} else {
		parts := strings.SplitN(upath, "/", 2)
		if len(parts) != 2 {
			panic(fmt.Sprintf(s("Path parser logic failure: %#v"), upath))
		}
		child := parts[0]
		grandchildren := parts[1]
		return findEntity(entity.Subentity(child, req), req, grandchildren)
	}
}
