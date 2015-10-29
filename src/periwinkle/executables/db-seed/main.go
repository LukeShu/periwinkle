// Copyright 2015 Luke Shumaker

package main

import (
	"periwinkle/cfg"
	"periwinkle/store"
)

func main() {
	store.DbSchema(cfg.DB)
	store.DbSeed(cfg.DB)
}
