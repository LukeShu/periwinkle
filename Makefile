# Set NET='' on the command line to not try to update things
NET ?= FORCE

# Set Q='' on the command line to enable verbose building
Q = @

# packages is the list of packages that we actually wrote and need built
packages = $(sort $(shell find src -type d -name '*.*' -prune -o -type f -name '*.go' -printf '%h\n'|cut -d/ -f2-))

# set deps to be a list of import strings of external packages we need to import
deps += bitbucket.org/ww/goautoneg
deps += github.com/dchest/captcha
deps += github.com/djherbis/times
deps += github.com/go-sql-driver/mysql
deps += github.com/jinzhu/gorm
deps += github.com/mattn/go-sqlite3
deps += golang.org/x/crypto/bcrypt

default: all
.PHONY: default

# RMRF is complicated because stupid NFS won't always allow deleting
# directory because of .nfsXXXXXXXXXXXX lock files
RMRF = { mkdir -p .rm && t=$$(mktemp -d .rm/XXXXXXXXXX) && mv -f -- $1 "$$t" 2>/dev/null && rm -rf -- "$$t"; }
DIRFAIL = { $(RMRF); false; }

# What directory is the Makefile in?
topdir := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))

# Configuration of the C compiler for C code called from Go
CFLAGS = -std=c99 -Wall -Wextra -Werror -Wno-old-style-declaration
CGO_CFLAGS = $(CFLAGS) -Wno-unused-parameter
CGO_ENABLED = 1
cgo_variables = CGO_ENABLED CGO_CFLAGS CGO_CPPFLAGS CGO_CXXFLAGS CGO_LDFLAGS CC CXX
export $(cgo_variables)

# A list of go source files; if any of these change, we need to rebuild
goext = go c s S cc cpp cxx h hh hpp hxx
gosrc_cmd = find -L src -name '.*' -prune -o \( -type f \( -false $(foreach e,$(goext),-o -name '*.$e') \) -o -type d \) -print
gosrc = $(shell $(gosrc_cmd))

# Iterate over external dependencies, and create a rule to download it
$(foreach d,$(deps),$(eval src/$d: $(NET); GOPATH='$(topdir)' go get -d -u $d || $(call DIRFAIL,$@)))

all: bin
.PHONY: all

#$(info $(gosrc) $(addprefix src/,$(deps)) $(addprefix .var.,$(cgo_variables)))
bin pkg: $(gosrc) $(addprefix src/,$(deps)) $(addprefix .var.,$(cgo_variables))
	GOPATH='$(topdir)' go install $(packages) || $(call DIRFAIL,bin pkg)
	$(Q)true $(foreach f,$^, && test $@ -nt $f ) || { \
		echo "# There's a discrepancy between Make and Go's dependency" && \
		echo "# tracking; nuking and starting over." && \
		PS4='' && set -x && \
		$(call RMRF,bin pkg) && \
		GOPATH='$(topdir)' go install $(packages) || $(call DIRFAIL,bin pkg); \
	}
	touch bin pkg

# Rule to nuke everything
clean:
	rm -rf -- .rm pkg bin src/*.*/ .var.*
.PHONY: clean

# Now, this is magic.  It stores the values of environment variables,
# so that if you change them in a way that would cause something to be
# rebuilt, then Make knows.
.var.%: FORCE
	$(Q)printf '%s' '$($*)' > .tmp$@ && { cmp -s .tmp$@ $@ && rm -f -- .tmp$@ || mv -Tf .tmp$@ $@; } || { rm -f -- .tmp$@; false; }

gofmt:
	gofmt -d $(addprefix src/,$(packages))
govet:
	GOPATH='$(topdir)' go vet $(packages)
.PHONY: gofmt govet

# Boilerplate
.SECONDARY:
.DELETE_ON_ERROR:
.PHONY: FORCE
