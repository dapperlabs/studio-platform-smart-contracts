.PHONY: test
test:
	$(MAKE) test -C ./lib/go

.PHONY: ci
ci:
	$(MAKE) ci -C ./lib/go

.PHONY: generate
generate:
	$(MAKE) generate -C ./lib/go/contracts
	$(MAKE) generate -C ./lib/go/templates
