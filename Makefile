# Copyright 2015 Luke Shumaker

# Set Q='' to enable verbose building
Q ?= @

# Set NET=FORCE to update network-downloaded things
#NET ?= FORCE

# Configuration of the C compiler for C code called from Go
CFLAGS = -std=c99 -Wall -Wextra -Werror -Wno-old-style-declaration
CGO_CFLAGS = $(CFLAGS) -Wno-unused-parameter
CGO_ENABLED = 1

# Set deps to be a list of import strings of external packages we need to import
deps += bitbucket.org/ww/goautoneg
deps += github.com/dchest/captcha
deps += github.com/djherbis/times
deps += github.com/evanphx/json-patch
deps += github.com/go-sql-driver/mysql
deps += github.com/jinzhu/gorm
deps += github.com/mattbaird/jsonpatch
deps += github.com/mattn/go-sqlite3
deps += golang.org/x/crypto/bcrypt
deps += lukeshu.com/git/go/libsystemd.git

# List of our packages and executables in them
packages = $(sort $(shell find src -type d -name '*.*' -prune -o -type f -name '*.go' -printf '%h\n'|cut -d/ -f2-))
executables = $(patsubst periwinkle/executables/%,%,$(filter periwinkle/executables/%,$(packages)))

srcdir := $(abspath $(patsubst %/,%,$(dir $(lastword $(MAKEFILE_LIST)))))
topdir := $(srcdir)

subdirs += $(topdir)/src/postfixpipe $(topdir)/HACKING

generate += $(addprefix $(topdir)/src/,$(godeps))
generate_secondary += $(topdir)/src/*.*/
build += $(addprefix $(topdir)/bin/,$(executables))
build_secondary += $(topdir)/bin $(topdir)/pkg $(topdir)/*.sqlite

ifeq (1,$(words $(MAKEFILE_LIST)))
  include $(topdir)/common.mk
endif

include $(topdir)/golang.mk
$(call goget,$(topdir),$(deps))

$(addprefix %/bin/,$(executables)): $(generate) $(configure) %/src $(call gosrc,$(topdir))
	$(call goinstall,$*,$(addprefix periwinkle/executables/,$(executables)))

gofmt:
	gofmt -d $(addprefix $(topdir)/src/,$(packages))
govet:
	GOPATH='$(abspath $(topdir))' go vet $(packages)
.PHONY: gofmt govet
