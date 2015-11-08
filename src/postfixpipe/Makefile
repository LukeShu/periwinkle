# Copyright 2015 Luke Shumaker

srcdir := $(patsubst %/,%,$(dir $(lastword $(MAKEFILE_LIST))))
topdir ?= ../..

generate := $(generate) $(srcdir)/sysexits.go

ifeq (1,$(words $(MAKEFILE_LIST)))
  include $(topdir)/common.mk
endif

%.go: %.go.sh
	$(call rel,$^) > $@

$(srcdir)/sysexits.go: /usr/include/sysexits.h
