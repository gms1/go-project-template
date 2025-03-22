.PHONY: all prerequisites clean format lint build test upgrade help

define usage =
usage: make [COMMAND]

COMMAND:
	all           ... Runs: clean, format, lint, build and test

	clean         ... Clean dist and tmp folders
	format        ... Format
	lint          ... Run Linter
	build         ... Build this module
	test          ... Run tests

	upgrade       ... upgrade dependencies
	prerequisites ... install or upgrade prerequisites

	help          ... Display this usage information
endef

all:
	@./scripts/all.sh

prerequisites:
	@./scripts/prerequisites.sh

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

upgrade:
	@./scripts/upgrade.sh

help:; @ $(info $(usage)) :
