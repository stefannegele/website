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

# Run in Docker dev container
docker-dev:
	docker build -f Dockerfile.dev -t website-dev .
	docker run --rm -it -v $(PWD):/site -p 1313:1313 website-dev make dev

clean:
	rm -rf public/ resources/ .hugo_build.lock
