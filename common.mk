# Copyright 2015 Luke Shumaker

rel = $(patsubst $(abspath .)/%,./%,$(abspath $1))

all: build
.PHONY: all

include $(addsuffix /Makefile,$(subdirs))

generate: $(generate)
.PHONY: generate

configure: generate $(configure)
.PHONY: configure

build: configure $(build)
.PHONY: build

# un-build
clean:
	rm -rf -- $(build) $(build_secondary)
.PHONY: clean

# un-configure
distclean: clean
	rm -rf -- $(configure) $(configure_secondary)
.PHONY: distclean

# un-generate
maintainer-clean: distclean
	rm -rf -- $(generate) $(generate_secondary)
.PHONY: maintainer-clean


# Now, this is magic.  It stores the values of environment variables,
# so that if you change them in a way that would cause something to be
# rebuilt, then Make knows.
.var.%: FORCE
	$(Q)printf '%s' '$($*)' > .tmp$@ && { cmp -s .tmp$@ $@ && rm -f -- .tmp$@ || mv -Tf .tmp$@ $@; } || { rm -f -- .tmp$@; false; }

.DELETE_ON_ERROR:
.SECONDARY:
.PHONY: FORCE
