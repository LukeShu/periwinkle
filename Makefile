# Copyright 2015 Luke Shumaker

# Set Q='' to enable verbose building
Q ?= @

# Set NET=FORCE to update network-downloaded things
#NET ?= FORCE

# What mode should gofmt run in? (another good option is `-w`)
GOFMT_MODE ?= -d

#POSTBUILD = systemctl --user restart listen-http.service listen-twilio.service || true

# Configuration of the C compiler for C code called from Go
CFLAGS = -std=c99 -Wall -Wextra -Werror -Wno-old-style-declaration
CGO_CFLAGS = $(CFLAGS) -Wno-unused-parameter
CGO_ENABLED = 1

# Set deps to be a list of import strings of external packages we need to import
deps += github.com/LukeShu/go-docopt
deps += github.com/dchest/captcha
deps += github.com/djherbis/times
deps += github.com/evanphx/json-patch
deps += github.com/go-sql-driver/mysql
deps += github.com/jinzhu/gorm
deps += github.com/mattbaird/jsonpatch
deps += github.com/mattn/go-sqlite3
deps += golang.org/x/crypto/bcrypt
deps += gopkg.in/yaml.v2
deps += lukeshu.com/git/go/libsystemd.git

# List of our packages and executables in them
packages = $(sort $(shell find src -type d -name '*.*' -not -name lukeshu.com -not -name '*.git' -prune -o -type f -name '*.go' -printf '%h\n'|cut -d/ -f2-))
toppackages = $(sort $(shell find src -type d -name '*.*' -not -name lukeshu.com -not -name '*.git' -prune -o -type f -name '*.go' -printf '%h\n'|cut -d/ -f2))
cmds = $(filter periwinkle/cmd/% locale/cmd/%,$(packages))

# What to ignore from golint
golint-filter = | grep -vE "/(sysexits|env|exit-status)\.go:[0-9]+:[0-9]+: don't use ALL_CAPS in Go names; use CamelCase"

srcdir := $(abspath $(patsubst %/,%,$(dir $(lastword $(MAKEFILE_LIST)))))
topdir := $(srcdir)

subdirs += $(topdir)/src/postfixpipe $(topdir)/src/locale/gettext $(topdir)/HACKING

generate += $(addprefix $(topdir)/src/,$(deps))
generate_secondary += $(topdir)/src/*.*/
build += $(addprefix $(topdir)/bin/,$(notdir $(cmds)))
build_secondary += $(topdir)/bin $(topdir)/pkg $(topdir)/*.sqlite

ifeq (1,$(words $(MAKEFILE_LIST)))
  include $(topdir)/common.mk
endif

include $(topdir)/devtools.mk

include $(topdir)/golang.mk
$(call goget,$(topdir),$(deps))

# Build all executables in one shot, because otherwise multiple
# instances of `go install` will not play nice with eachother in
# `pkg/`
$(addprefix %/bin/,$(notdir $(cmds))): $(generate) $(configure) %/src $(call gosrc,$(topdir))
	$(call goinstall,$*,$(cmds))
	$(POSTBUILD)

check: gofmt goimports govet gotest
.PHONY: check

# directory-oriented
gofmt: generate
	{ gofmt -s $(GOFMT_MODE) $(addprefix $(topdir)/src/,$(toppackages)) 2>&1 | tee /dev/stderr | test -z "$$(cat)"; } 2>&1
goimports: generate $(GOIMPORTS)
	{ goimports -d $(addprefix $(topdir)/src/,$(toppackages)) 2>&1 | tee /dev/stderr | test -z "$$(cat)"; } 2>&1
govet: generate
	GOPATH='$(abspath $(topdir))' go tool vet -composites=false $(addprefix $(topdir)/src/,$(toppackages))
.PHONY: gofmt goimports govet

# package-oriented
gotest: build
	GOPATH='$(abspath $(topdir))' go test -cover -v $(packages)
.PHONY: gotest

golint: generate $(GOLINT)
	export GOPATH='$(abspath $(topdir))'; { { $(foreach p,$(packages),golint $p; )} $(golint-filter) | tee /dev/stderr | test -z "$$(cat)"; } 2>&1
.PHONY: golint
