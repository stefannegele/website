#!/usr/bin/env python3
"""
Fetches INNOQ articles for Stefan Negele and merges them into data/innoq.json.
Append-only, deduplication by URL. No external dependencies.
"""

import json
import os
import re
import urllib.request
import xml.etree.ElementTree as ET
from datetime import datetime

ARTICLES_FEED = "https://www.innoq.com/en/written.atom"
AUTHOR_NAME="Stefan Negele"
NS = {"atom": "http://www.w3.org/2005/Atom"}
OUTPUT = "data/innoq.json"


def fetch(url):
    req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0 (stefannegele/website)"})
    with urllib.request.urlopen(req, timeout=15) as r:
        return r.read()


def parse_date(raw):
    raw = raw.strip()
    for fmt in ("%Y-%m-%dT%H:%M:%S%z", "%Y-%m-%dT%H:%M:%SZ", "%Y-%m-%d"):
        try:
            return datetime.strptime(raw, fmt).strftime("%Y-%m-%d")
        except ValueError:
            continue
    return raw[:10]


def strip_html(text):
    return re.sub(r"<[^>]+>", "", text).strip()


def fetch_articles():
    print("Fetching", ARTICLES_FEED, "...")
    root = ET.fromstring(fetch(ARTICLES_FEED))
    posts = []
    for entry in root.findall("atom:entry", NS):
        authors = entry.findall("atom:author", NS)
        is_stefan = any(
            AUTHOR_NAME.lower() in (a.findtext("atom:name", namespaces=NS) or "").lower()
            for a in authors
        )
        if not is_stefan:
            continue

        title = (entry.findtext("atom:title", namespaces=NS) or "").strip()

        url = ""
        for link in entry.findall("atom:link", NS):
            if link.get("rel", "alternate") == "alternate":
                url = link.get("href", "")
                break
        if not url:
            link_el = entry.find("atom:link", NS)
            if link_el is not None:
                url = link_el.get("href", "")

        raw_date = (entry.findtext("atom:published", namespaces=NS)
                    or entry.findtext("atom:updated", namespaces=NS) or "")
        date = parse_date(raw_date) if raw_date else ""

        summary_el = entry.find("atom:summary", NS)
        summary = (summary_el.text or "").strip() if summary_el is not None else ""
        if "<" in summary:
            summary = strip_html(summary)

        posts.append({"title": title, "url": url, "date": date, "summary": summary})

    return posts


def load_existing():
    if not os.path.exists(OUTPUT):
        return []
    with open(OUTPUT, encoding="utf-8") as f:
        return json.load(f)


def merge(existing, new):
    known_urls = {item["url"] for item in existing}
    added = 0
    for item in new:
        if item["url"] not in known_urls:
            existing.append(item)
            known_urls.add(item["url"])
            added += 1
    existing.sort(key=lambda p: p.get("date", ""), reverse=True)
    return existing, added


def main():
    fetched = fetch_articles()
    print("Found", len(fetched), "article(s) in feed")

    existing = load_existing()
    merged, added = merge(existing, fetched)

    os.makedirs("data", exist_ok=True)
    with open(OUTPUT, "w", encoding="utf-8") as f:
        json.dump(merged, f, indent=2, ensure_ascii=False)
    print("Added", added, "new article(s). Total:", len(merged), "in", OUTPUT)


if __name__ == "__main__":
    main()
