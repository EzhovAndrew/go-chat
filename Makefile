.PHONY: help
.DEFAULT_GOAL := help

# Include modular makefiles
include scripts/make/tools.mk
include scripts/make/proto.mk
include scripts/make/dev.mk

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

