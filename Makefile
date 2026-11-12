# This Makefile is a proxy for the Justfile.
# It allows you to run `make <target>` which will be forwarded to `just <target>`.
# Variables given on the command line are passed on as well, e.g. `make build VERSION=1.2.3`.

.PHONY: default
default:
	@just --list

.PHONY: %
%:
	@just $(MAKEOVERRIDES) $@
