# Copyright 2015 Luke Shumaker

srcdir := $(patsubst %/,%,$(dir $(lastword $(MAKEFILE_LIST))))
topdir ?= ../../..

generate := $(generate) $(srcdir)/bits_locale.go
generate_secondary := $(generate_secondary) $(srcdir)/bits_locale.go.exe $(srcdir)/bits_locale.go.c

ifeq (1,$(words $(MAKEFILE_LIST)))
  include $(topdir)/common.mk
endif

$(srcdir)/bits_locale.go: $(srcdir)/bits_locale.go.exe
	rm -f -- $@.c
	$(call rel,$^) | gofmt > $@

$(srcdir)/bits_locale.go.exe: $(srcdir)/bits_locale.go.c .var.CC .var.CFLAGS .var.CPPFLAGS .var.LDFLAGS
	$(CC) $(CFLAGS) $(CPPFLAGS) $(LDFLAGS) $< -o $@

$(srcdir)/bits_locale.go.c: $(srcdir)/bits_locale.go.c.sh
	$(call rel,$^) > $@
