# Copyright Luke Shumaker

d := $(patsubst %/,%,$(dir $(lastword $(MAKEFILE_LIST))))

generate += $d/sysexits.go

ifeq (1,$(words $(MAKEFILE_LIST)))

all: generate
.PHONY: all

generate: $(generate)
.PHONY: generate

clean:
.PHONY: clean
maintainer-clean: clean
	rm -f -- $(generate)
.PHONY: maintainer-clean

.DELETE_ON_ERROR:
.SECONDARY:

else

# fix these to the current value of `d`
generate := $(generate)

endif

%.go: %.go.sh
	./$^ > $@
$(d)/sysexits.go: /usr/include/sysexits.h
