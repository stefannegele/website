.PHONY: dev fetch build post external

dev:
	hugo server --buildDrafts --bind 0.0.0.0

fetch:
	python3 scripts/fetch-innoq.py

build: fetch
	hugo --minify

post:
	@test -n "$(TITLE)" || (echo "Usage: make post TITLE=\"Mein Artikel\"" && exit 1)
	hugo new posts/$(shell echo "$(TITLE)" | tr '[:upper:]' '[:lower:]' | tr ' ' '-').md

external:
	@test -n "$(TITLE)" || (echo "Usage: make external TITLE=\"Artikel Titel\"" && exit 1)
	hugo new external/$(shell echo "$(TITLE)" | tr '[:upper:]' '[:lower:]' | tr ' ' '-').md
