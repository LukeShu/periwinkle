# Copyright © 2015 Luke Shumaker
# This work is free. You can redistribute it and/or modify it under the
# terms of the Do What The Fuck You Want To Public License, Version 2,
# as published by Sam Hocevar. See http://www.wtfpl.net/ for more details.

rel = $(patsubst $(abspath .)/%,./%,$(abspath $1))

all: build
.PHONY: all

-include $(addsuffix /Makefile,$(subdirs))

generate: $(generate)
.PHONY: generate

configure: generate $(configure)
.PHONY: configure

build: configure $(build)
.PHONY: build

install: build $(install)
.PHONY: install

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

# un-install
uninstall:
	rm -f -- $(install)
	rmdir -p -- $(sort $(dir $(install))) 2>/dev/null || true
.PHONY: uninstall


# Now, this is magic.  It stores the values of environment variables,
# so that if you change them in a way that would cause something to be
# rebuilt, then Make knows.
.var.%: FORCE
	$(Q)printf '%s' '$($*)' > .tmp$@ && { cmp -s .tmp$@ $@ && rm -f -- .tmp$@ || mv -Tf .tmp$@ $@; } || { rm -f -- .tmp$@; false; }

.DELETE_ON_ERROR:
.SECONDARY:
.PHONY: FORCE
