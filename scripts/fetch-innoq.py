#!/usr/bin/env python3
"""
Fetches INNOQ articles and podcast episodes for Stefan Negele.
No external dependencies required.
"""

import json
import os
import re
import urllib.request
import xml.etree.ElementTree as ET
from datetime import datetime

ARTICLES_FEED = "https://www.innoq.com/en/written.atom"
PODCAST_FEED = "https://innoq.podigee.io/feed/mp3"
AUTHOR_NAME = "Stefan Negele"
NS_ATOM = {"atom": "http://www.w3.org/2005/Atom"}
NS_RSS = {"itunes": "http://www.itunes.com/dtds/podcast-1.0.dtd"}
OUTPUT = "data/innoq.json"


def fetch(url):
    req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0 (stefannegele/website)"})
    with urllib.request.urlopen(req, timeout=15) as r:
        return r.read()


def parse_date(raw):
    raw = raw.strip()
    for fmt in ("%Y-%m-%dT%H:%M:%S%z", "%Y-%m-%dT%H:%M:%SZ", "%Y-%m-%d",
                "%a, %d %b %Y %H:%M:%S %z", "%a, %d %b %Y %H:%M:%S GMT"):
        try:
            dt = datetime.strptime(raw, fmt)
            return dt.strftime("%Y-%m-%d")
        except ValueError:
            continue
    return raw[:10]


def strip_html(text):
    return re.sub(r"<[^>]+>", "", text).strip()


def fetch_articles():
    print("Fetching articles from", ARTICLES_FEED, "...")
    root = ET.fromstring(fetch(ARTICLES_FEED))
    posts = []
    for entry in root.findall("atom:entry", NS_ATOM):
        authors = entry.findall("atom:author", NS_ATOM)
        is_stefan = any(
            AUTHOR_NAME.lower() in (a.findtext("atom:name", namespaces=NS_ATOM) or "").lower()
            for a in authors
        )
        if not is_stefan:
            continue

        title = (entry.findtext("atom:title", namespaces=NS_ATOM) or "").strip()
        url = ""
        for link in entry.findall("atom:link", NS_ATOM):
            if link.get("rel", "alternate") == "alternate":
                url = link.get("href", "")
                break
        if not url:
            link_el = entry.find("atom:link", NS_ATOM)
            if link_el is not None:
                url = link_el.get("href", "")

        raw_date = (entry.findtext("atom:published", namespaces=NS_ATOM)
                    or entry.findtext("atom:updated", namespaces=NS_ATOM) or "")
        date = parse_date(raw_date) if raw_date else ""

        summary_el = entry.find("atom:summary", NS_ATOM)
        summary = (summary_el.text or "").strip() if summary_el is not None else ""
        if "<" in summary:
            summary = strip_html(summary)

        posts.append({"title": title, "url": url, "date": date, "summary": summary, "type": "article"})

    print("Found", len(posts), "article(s)")
    return posts


def fetch_podcasts():
    print("Fetching podcasts from", PODCAST_FEED, "...")
    root = ET.fromstring(fetch(PODCAST_FEED))
    channel = root.find("channel")
    if channel is None:
        return []

    posts = []
    for item in channel.findall("item"):
        # Check title or description for Stefan's name
        title = (item.findtext("title") or "").strip()
        desc = strip_html(item.findtext("description") or "")

        if AUTHOR_NAME.lower() not in title.lower() and AUTHOR_NAME.lower() not in desc.lower():
            continue

        url = (item.findtext("link") or "").strip()
        raw_date = (item.findtext("pubDate") or "").strip()
        date = parse_date(raw_date) if raw_date else ""

        # Short summary: first sentence of description
        summary = desc.split(".")[0].strip() if desc else ""
        if len(summary) > 160:
            summary = summary[:157] + "..."

        posts.append({"title": title, "url": url, "date": date, "summary": summary, "type": "podcast"})

    print("Found", len(posts), "podcast episode(s)")
    return posts


def main():
    items = fetch_articles() + fetch_podcasts()
    items.sort(key=lambda p: p.get("date", ""), reverse=True)
    os.makedirs("data", exist_ok=True)
    with open(OUTPUT, "w", encoding="utf-8") as f:
        json.dump(items, f, indent=2, ensure_ascii=False)
    print("Written", len(items), "items to", OUTPUT)


if __name__ == "__main__":
    main()
