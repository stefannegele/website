#!/usr/bin/env python3
"""Fetch INNOQ RSS feed and write data/innoq.json for Hugo."""

import json
import re
import hashlib
import urllib.request
import xml.etree.ElementTree as ET
from datetime import datetime

FEED_URL = "https://www.innoq.com/en/written.atom"
AUTHOR_EMAIL = "stefan.negele@innoq.com"
AUTHOR_NAME = "Stefan Negele"
OUTPUT = "data/innoq.json"

NS = {"atom": "http://www.w3.org/2005/Atom"}


def strip_html(s):
    s = re.sub(r"<[^>]+>", "", s or "")
    return re.sub(r"\s+", " ", s).strip()


def fetch():
    req = urllib.request.Request(FEED_URL, headers={"User-Agent": "Mozilla/5.0"})
    with urllib.request.urlopen(req, timeout=15) as r:
        return r.read()


def parse(xml_bytes):
    root = ET.fromstring(xml_bytes)
    posts = []
    for entry in root.findall("atom:entry", NS):
        authors = entry.findall("atom:author", NS)
        is_stefan = any(
            (a.findtext("atom:email", namespaces=NS) or "").lower() == AUTHOR_EMAIL
            or (a.findtext("atom:name", namespaces=NS) or "").lower() == AUTHOR_NAME.lower()
            for a in authors
        )
        if not is_stefan:
            continue

        title = entry.findtext("atom:title", namespaces=NS) or ""
        url = ""
        for link in entry.findall("atom:link", NS):
            if link.get("rel") == "alternate":
                url = link.get("href", "")
                break

        pub = (
            entry.findtext("atom:published", namespaces=NS)
            or entry.findtext("atom:updated", namespaces=NS)
            or ""
        )
        try:
            dt = datetime.fromisoformat(pub.replace("Z", "+00:00"))
            date_str = dt.strftime("%Y-%m-%d")
        except Exception:
            date_str = pub[:10] if pub else ""

        summary_raw = entry.findtext("atom:summary", namespaces=NS) or ""
        summary = strip_html(summary_raw)[:300]
        if len(strip_html(summary_raw)) > 300:
            summary += "\u2026"

        slug = hashlib.md5(url.encode()).hexdigest()[:8]
        posts.append(
            {
                "title": title,
                "url": url,
                "date": date_str,
                "summary": summary,
                "source": "innoq",
                "slug": slug,
            }
        )

    return posts


if __name__ == "__main__":
    import os
    os.makedirs("data", exist_ok=True)
    xml_data = fetch()
    posts = parse(xml_data)
    with open(OUTPUT, "w") as f:
        json.dump(posts, f, indent=2, ensure_ascii=False)
    print(f"Wrote {len(posts)} posts to {OUTPUT}")
