%.svg: %.dot
	dot -Tsvg $< > $@
%.png: %.dot
	dot -Tpng $< > $@
%.html: %.md
	pandoc -s -f markdown -t html5 $< > $@
%.pdf: %.md
	pandoc -s -f markdown -t pdf $< > $@

.DELETE_ON_ERROR:
.SECONDARY:
