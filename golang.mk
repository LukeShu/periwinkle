# Copyright 2015 Luke Shumaker

_golang_cgo_variables = CGO_ENABLED CGO_CFLAGS CGO_CPPFLAGS CGO_CXXFLAGS CGO_LDFLAGS CC CXX
export $(_golang_cgo_variables)
_golang_src_cmd = find -L $1/src -name '.*' -prune -o \( -type f \( -false $(foreach e,go c s S cc cpp cxx h hh hpp hxx,-o -name '*.$e') \) -o -type d \) -print

# Iterate over external dependencies, and create a rule to download it
goget = $(foreach d,$2,$(eval $1/src/$d: $(NET); GOPATH='$(abspath $1)' go get -d -u $d || { rm -rf -- $$@; false; }))

gosrc = $(shell $(_golang_src_cmd)) $(addprefix .var.,$(_golang_cgo_variables))
define goinstall
	$(Q)for target in $(addprefix $1/bin/,$(notdir $2)); do              \
		if test -e $$target; then                                    \
			for dep in $(filter .var.%,$^); do                   \
				if test $$dep -nt $$target; then             \
					rm -rf -- $1/bin $1/pkg || exit $$?; \
					exit 0;                              \
				fi                                           \
			done                                                 \
		fi                                                           \
	done
	GOPATH='$(abspath $1)' go install -x $2
	$(Q)true $(foreach e,$(notdir $2), && test -f $1/bin/$e -a -x $1/bin/$e && touch $1/bin/$e)
endef
