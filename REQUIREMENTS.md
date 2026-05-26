# Website Requirements – Stefan Negele

## Ziel

Persönliche Website / Blog von Stefan Negele. Statische Site, gehostet auf GitHub Pages.

---

## Hosting & Deployment

- **Platform:** GitHub Pages (`https://stefannegele.github.io/website/`)
- **Deployment:** GitHub Actions, automatisch bei Push auf `main`
- **Zusätzlicher Trigger:** Täglicher Cron-Job (06:00 UTC) um externe Inhalte aktuell zu halten
- **Build-Tool:** Hugo (extended, aktuell v0.147.1)

---

## Inhaltstypen

### 1. Interne Artikel (`source: internal`)
- Markdown-Dateien unter `content/posts/`
- Felder: `title`, `date`, `summary`, `draft`
- Vollständiger Inhalt auf der Website

### 2. Externe Links (`source: external`)
- Markdown-Dateien unter `content/external/`
- Felder: `title`, `date`, `summary`, `external_url`
- Verlinken auf externe Seiten (z.B. INNOQ-Blog)

### 3. Statische Seiten
- Impressum (`content/impressum.md`)

---

## INNOQ RSS-Integration

- Quelle: `https://www.innoq.com/en/written.atom`
- Filter: Nur Artikel von `stefan.negele@innoq.com` / `Stefan Negele`
- Script (`scripts/fetch-innoq.py`) wird beim Build ausgeführt und schreibt `data/innoq.json`
- INNOQ-Artikel werden auf der Startseite zusammen mit internen Artikeln angezeigt

---

## Design & Sprache

- **Sprache:** Deutsch (`languageCode = "de"`)
- **Titel:** Stefan Negele
- **Beschreibung:** Artikel über Data Architecture, verteilte Systeme und Software Engineering
- **Kein externes Theme** – eigene Layouts
- **Code-Highlighting:** Dracula-Style
- Minimalistisches, schlichtes Design (kein Framework wie Bootstrap/Tailwind)

---

## Lokale Entwicklung

- `make dev` → Hugo Dev-Server auf Port 1313 mit Live Reload
- `make fetch` → INNOQ RSS fetchen → `data/innoq.json`
- `make build` → fetch + hugo --minify → `public/`
- `make post TITLE="..."` → neuen internen Post als Draft anlegen
- `make external TITLE="..."` → neuen externen Link-Post anlegen

---

## Nicht-Anforderungen

- Kein Server-Side-Rendering, kein Backend
- Kein Admin-Interface / CMS
- Keine Datenbank
- Kein JavaScript-Framework
