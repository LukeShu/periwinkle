# Set NET='' on the command line to not try to update things
NET ?= NET

# Set Q='' on the command line to enable verbose building
Q = @

# packages is the list of packages that we actually wrote and need built
packages = $(sort $(shell cd src && find periwinkle/ -name '*.go' -printf '%h\n'))

# set deps to be a list of import strings of external packages we need to import
deps += bitbucket.org/ww/goautoneg
deps += github.com/dchest/captcha
deps += github.com/go-sql-driver/mysql
deps += github.com/jinzhu/gorm
deps += golang.org/x/crypto/bcrypt
deps += github.com/mattn/go-sqlite3

default: all
.PHONY: default

# What directory is the Makefile in?
topdir := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))

# Configuration of the C compiler for C code called from Go
CFLAGS = -std=c99 -Wall -Wextra -Werror -pedantic
CGO_CFLAGS = $(CFLAGS) -Wno-unused-parameter
CGO_ENABLED = 1
cgo_variables = CGO_ENABLED CGO_CFLAGS CGO_CPPFLAGS CGO_CXXFLAGS CGO_LDFLAGS CC CXX
export $(cgo_variables)

# A list of go source files; if any of these change, we need to rebuild
goext = go c s S cc cpp cxx h hh hpp hxx
gosrc_cmd = find -L src -name '.*' -prune -o \( -type f \( -false $(foreach e,$(goext),-o -name '*.$e') \) -o -type d \) -print
gosrc = $(shell $(gosrc_cmd))

# Iterate over external dependencies, and create a rule to download it
$(foreach d,$(deps),$(eval src/$d: $(NET); GOPATH='$(topdir)' go get -d -u $d || { rm -rf -- $$@; false; }))

all: bin
.PHONY: all

# The rule to build the Go code.  The first line nukes the built files
# if there is a discrepancy between Make and Go's internal
# dependency tracker.
bin pkg: $(gosrc) $(addprefix src/,$(deps)) $(addprefix .var.,$(cgo_variables))
	$(Q)true $(foreach f,$(filter-out .var.%,$^), && test $@ -nt $f ) || rm -rf -- bin pkg
	GOPATH='$(topdir)' go install $(packages) || { rm -rf -- bin; false; }

# Rule to nuke everything
clean:
	rm -rf -- pkg bin src/*.*/ .var.*
.PHONY: clean

# Now, this is magic.  It stores the values of environment variables,
# so that if you change them in a way that would cause something to be
# rebuilt, then Make knows.
.var.%: FORCE
	$(Q)printf '%s' '$($*)' > .tmp$@ && { cmp -s .tmp$@ $@ && rm -f -- .tmp$@ || mv -Tf .tmp$@ $@; } || { rm -f -- .tmp$@; false; }

# Boilerplate
.SECONDARY:
.DELETE_ON_ERROR:
.PHONY: FORCE NET
