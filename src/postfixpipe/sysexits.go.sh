#!/usr/bin/env bash
# Copyright 2015 Luke Shumaker

sysexits_h=$1
{
	declare -i n=0
	IFS=
	< "$sysexits_h" \
	sed -r -e 's/#define EX_(\S+)\s/EX_\1 uint8 = /' -e '/#/d' |
	    while read -r line; do
		    echo "$line"
		    if [[ "$line" == " */" ]]; then
			    n+=1
		    fi
		    if [[ $n == 2 ]]; then
			    echo package postfixpipe
			    echo 'const ('
			    n+=1
		    fi
	    done
	echo ')'
} | gofmt
