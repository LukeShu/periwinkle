// Copyright 2015 Luke Shumaker

package main

import (
	"periwinkle/cfg"
	"periwinkle/store"
)

func main() {
	store.SchemaDrop(cfg.DB)
}
