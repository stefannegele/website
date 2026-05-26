#!/usr/bin/env python3
"""
Fetches the INNOQ Atom feed and writes Stefan Negele's articles to data/innoq.json.
No external dependencies required.
"""

import json
import os
import re
import urllib.request
import xml.etree.ElementTree as ET
from datetime import datetime

FEED_URL = "https://www.innoq.com/en/written.atom"
AUTHOR_EMAIL = "stefan.negele@innoq.com"
AUTHOR_NAME = "Stefan Negele"
NS = {"atom": "http://www.w3.org/2005/Atom"}
OUTPUT = "data/innoq.json"


def fetch_feed(url):
    req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0 (stefannegele/website)"})
    with urllib.request.urlopen(req, timeout=15) as r:
        return r.read()


def is_stefan(entry):
    authors = entry.findall("atom:author", NS)
    for author in authors:
        email = (author.findtext("atom:email", namespaces=NS) or "").strip().lower()
        name = (author.findtext("atom:name", namespaces=NS) or "").strip()
        if email == AUTHOR_EMAIL or AUTHOR_NAME.lower() in name.lower():
            return True
    return False


def parse_date(raw):
    raw = raw.strip()
    for fmt in ("%Y-%m-%dT%H:%M:%S%z", "%Y-%m-%dT%H:%M:%SZ", "%Y-%m-%d"):
        try:
            dt = datetime.strptime(raw, fmt)
            return dt.strftime("%Y-%m-%d")
        except ValueError:
            continue
    return raw[:10]


def parse_feed(xml_bytes):
    root = ET.fromstring(xml_bytes)
    posts = []
    for entry in root.findall("atom:entry", NS):
        if not is_stefan(entry):
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

        raw_date = (
            entry.findtext("atom:published", namespaces=NS)
            or entry.findtext("atom:updated", namespaces=NS)
            or ""
        )
        date = parse_date(raw_date) if raw_date else ""

        summary_el = entry.find("atom:summary", NS)
        summary = (summary_el.text or "").strip() if summary_el is not None else ""
        if "<" in summary:
            summary = re.sub(r"<[^>]+>", "", summary).strip()

        posts.append({
            "title": title,
            "url": url,
            "date": date,
            "summary": summary,
            "source": "innoq",
        })

    posts.sort(key=lambda p: p.get("date", ""), reverse=True)
    return posts


def main():
    print("Fetching", FEED_URL, "...")
    xml_bytes = fetch_feed(FEED_URL)
    posts = parse_feed(xml_bytes)
    print("Found", len(posts), "article(s) by", AUTHOR_NAME)
    os.makedirs("data", exist_ok=True)
    with open(OUTPUT, "w", encoding="utf-8") as f:
        json.dump(posts, f, indent=2, ensure_ascii=False)
    print("Written to", OUTPUT)


if __name__ == "__main__":
    main()
