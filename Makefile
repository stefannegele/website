.PHONY: dev build fetch clean setup-pages

# Local dev server with live reload
dev:
	hugo server --buildDrafts --bind 0.0.0.0 --port 1313

# Fetch INNOQ RSS into data/innoq.json
fetch:
	python3 scripts/fetch-innoq.py

# Build static site
build: fetch
	hugo --minify

# New internal blog post: make post TITLE="Mein Artikel"
post:
	hugo new content posts/$(shell echo "$(TITLE)" | tr '[:upper:]' '[:lower:]' | sed 's/ /-/g').md

# New external link: make external TITLE="Artikel Titel"
external:
	hugo new content external/$(shell echo "$(TITLE)" | tr '[:upper:]' '[:lower:]' | sed 's/ /-/g').md

clean:
	rm -rf public/ resources/ .hugo_build.lock
