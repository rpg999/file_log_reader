include src/*/Makefile

format: format-fix-fast

# common typo
fromat: format-fix-fast

format-fix-fast:
	gofmt -w -s src/