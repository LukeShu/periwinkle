# Copyright 2015 Luke Shumaker

_devtools_rest = $(wordlist 2,$(words $1),$1)
_devtools_merge = $(firstword $2)$(if $(call _devtools_rest,$2),$1$(call _devtools_merge,$1,$(call _devtools_rest,$2)))
_devtools_pathsearch = $(firstword $(wildcard $(addsuffix /$(1),$(subst :, ,$(PATH)))) $(topdir)/devtools/$(1)/bin/$(1))

GOIMPORTS = $(call _devtools_pathsearch,goimports)

$(topdir)/devtools/goimports/bin/goimports: $(NET)
	mkdir -p $(dir $(@D))
	cd $(dir $(@D)) && GOPATH=$$PWD go get -u golang.org/x/tools/cmd/goimports

GOLINT = $(call _devtools_pathsearch,golint)

$(topdir)/devtools/golint/bin/golint: $(NET)
	mkdir -p $(dir $(@D))
	cd $(dir $(@D)) && GOPATH=$$PWD go get -u github.com/golang/lint/golint

_devtools_path = $(filter $(topdir)/%,$(patsubst %/,%,$(dir $(GOIMPORTS) $(GOLINT))))
PATH := $(call _devtools_merge,:,$(_devtools_path) $(PATH))
