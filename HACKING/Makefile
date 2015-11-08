# Copyright 2015 Luke Shumaker

srcdir := $(patsubst %/,%,$(dir $(lastword $(MAKEFILE_LIST))))
topdir ?= ..

build_secondary += $(srcdir)/*.html
build_secondary += $(srcdir)/*.pdf
build_secondary += $(srcdir)/*.png
build_secondary += $(srcdir)/*.svg
build_secondary := $(build_secondary)

ifeq (1,$(words $(MAKEFILE_LIST)))
  include $(topdir)/common.mk
endif

%.svg: %.dot
	dot -Tsvg $< > $@
%.png: %.dot
	dot -Tpng $< > $@
%.html: %.md
	pandoc -s -f markdown -t html5 $< > $@
%.pdf: %.md
	pandoc -s -f markdown -t pdf $< > $@
