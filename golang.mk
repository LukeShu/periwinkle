# Copyright 2015 Luke Shumaker

_golang_cgo_variables = CGO_ENABLED CGO_CFLAGS CGO_CPPFLAGS CGO_CXXFLAGS CGO_LDFLAGS CC CXX
export $(_golang_cgo_variables)
_golang_src_cmd = find -L $1/src -name '.*' -prune -o \( -type f \( -false $(foreach e,go c s S cc cpp cxx h hh hpp hxx,-o -name '*.$e') \) -o -type d \) -print

# Iterate over external dependencies, and create a rule to download it
goget = $(foreach d,$2,$(eval $1/src/$d: $(NET); GOPATH='$(abspath $1)' go get -d -u $d || { rm -rf -- $$@; false; }))

gosrc = $(shell $(_golang_src_cmd)) $(addprefix .var.,$(_golang_cgo_variables))
define goinstall
	$(Q)true $(foreach f,$(filter .var.%,$^), && test $@ -nt $f ) || rm -rf -- $1/bin $1/pkg
	GOPATH='$(abspath $1)' go install $2
	$(Q)true $(foreach e,$(notdir $2), && test -f $1/bin/$e && test -x $1/bin/$e && touch $1/bin/$e)
endef
