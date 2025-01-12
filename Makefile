.PHONY: all clean format lint build test help

define usage =
usage: make [COMMAND]

COMMAND:
	all     ... Runs: clean, format, lint, build and test
	clean   ... Clean up
	format  ... Format
	lint    ... Run Linter
	build   ... Build this module
	test    ... Run tests
	help    ... Display this usage information
endef

all:
	@./scripts/all.sh

clean:
	@./scripts/clean.sh

format:
	@./scripts/format.sh

lint:
	@./scripts/lint.sh

build:
	@./scripts/build.sh

test:
	@./scripts/test.sh

help:; @ $(info $(usage)) :
