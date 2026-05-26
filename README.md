# stefannegele.de

Personal website and blog. Built with [Hugo](https://gohugo.io/), hosted on GitHub Pages.

**Live:** https://stefannegele.github.io/website/

---

## Content

| Type | Location | Description |
|---|---|---|
| Articles | `content/posts/` | Internal blog posts (Markdown) |
| External links | `content/external/` | Links to external publications |
| INNOQ articles | `data/innoq.json` | Fetched from INNOQ feed, append-only |
| Static pages | `content/` | Imprint etc. |

## Development

```bash
make fetch   # fetch new INNOQ articles → data/innoq.json
make dev     # local server at localhost:1313
make build   # production build → public/
```

New post:
```bash
make post TITLE="My Article"
make external TITLE="Article Title"
```

## Deployment

GitHub Actions deploys automatically on push to `main` and via daily cron (06:00 UTC).  
Requires **Pages → Source → GitHub Actions** in repository settings.
