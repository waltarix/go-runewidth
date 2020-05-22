TABLE_FILE_VERSION := 13.0.0-r1
TABLE_FILE := https://github.com/waltarix/localedata/releases/download/$(TABLE_FILE_VERSION)/runewidth_table.go

CACHE_DIR := .cache
CACHE_FILE := $(CACHE_DIR)/$(TABLE_FILE_VERSION)

runewidth_table.go: $(CACHE_FILE)
	curl -sL $(TABLE_FILE) > $@

.PHONY: test
test: runewidth_table.go
	./go.test.sh

.PHONY: bench
bench: runewidth_table.go
	go test -bench .

$(CACHE_DIR):
	mkdir -p $@
$(CACHE_FILE): $(CACHE_DIR)
	touch $@

.PHONY: clean
clean:
	rm -rf $(CACHE_DIR) coverage.txt
