---
name: orkl
description: "Query the ORKL cyber threat intelligence library using the orkl CLI. Use this skill whenever the user wants to research threat actors, look up CTI reports or intelligence sources, search for APT groups, find reports by SHA1 hash, investigate malware campaigns, or explore the ORKL library — even if they just say things like 'look up APT29', 'find threat reports on ransomware', 'what does ORKL say about...', or 'search threat intel for X'. Always reach for this skill for any threat intelligence research task."
---

# ORKL

Use the `orkl` CLI to query the ORKL Cyber Threat Intelligence Library — a comprehensive corpus of publicly released CTI reports, threat actors, and intelligence sources.

**Install if missing:**

```bash
go install github.com/cyberbuff/orkl/cmd/orkl@latest
```

## Global flags

| Flag                   | Default | Description                        |
| ---------------------- | ------- | ---------------------------------- |
| `--timeout <duration>` | `60s`   | Request timeout (e.g. `10s`, `1m`) |

---

## Common workflows

### Research a threat actor (e.g. APT29, Lazarus)

```bash
# 1. Search for reports
orkl library-search "APT29" --full --limit 10 | jq '.[].title'

# 2. Find the actor entry
orkl ta-entries | jq '.[] | select(.name | test("APT29"; "i"))'

# 3. Get full actor details by UUID
orkl ta-entry <uuid> | jq .
```

### Look up a suspicious file hash

```bash
orkl library-entry-hash <sha1> | jq .
```

### Explore recent reports

```bash
orkl library-entries --limit 20 --order-by created_at --order desc | jq '.[].title'
```

### Find reports from a specific source

```bash
# List sources, find UUID
orkl source-entries | jq '.[] | select(.name | test("Mandiant"; "i"))'

# Get source with all linked reports
orkl source-entry <uuid> --full | jq .
```

---

## Library commands

### Library metadata

```bash
orkl library-info          # Overall library stats
orkl library-version       # Current version
orkl library-versions [--limit N] [--offset N] [--order asc|desc]
```

### List entries

```bash
orkl library-entries [--limit N] [--offset N] [--order-by <field>] [--order asc|desc] [--origin pdf|web]
```

### Get a specific entry

```bash
orkl library-entry <uuid>           # By UUID
orkl library-entry-hash <sha1>      # By SHA1 hash
```

### Search the library

```bash
orkl library-search <query> [--full] [--limit N] [--origin pdf|web]
```

| Flag                | Description                                    |
| ------------------- | ---------------------------------------------- |
| `--full`            | Return full entry objects instead of summaries |
| `--limit N`         | Maximum number of results to return            |
| `--origin pdf\|web` | Filter by origin format                        |

---

## Source commands

```bash
orkl source-entries                  # List all intelligence sources
orkl source-entry <uuid> [--full]    # Get source details; --full includes linked reports
```

---

## Threat actor commands

```bash
orkl ta-entries          # List all threat actors
orkl ta-entry <uuid>     # Get full details for one actor
```

---

## Output & filtering

All commands output raw JSON. Pipe through `jq` to filter:

```bash
# Pretty-print any result
orkl library-info | jq .

# Extract just titles from a search
orkl library-search "ransomware" --limit 5 | jq '.[].title'

# Get name, aliases, and country for a threat actor
orkl ta-entry <uuid> | jq '{name: .name, aliases: .aliases, country: .country}'

# List sources by name
orkl source-entries | jq '.[].name'

# Find entries from a specific year
orkl library-entries --limit 100 | jq '[.[] | select(.created_at | startswith("2024"))]'
```

---

## Tips

- **Start broad, then drill down**: Use `library-search` to find relevant entries, grab a UUID, then use `library-entry` or `ta-entry` for full details.
- **Use `--full` on searches** when you need complete report content, not just titles/summaries.
- **SHA1 lookups** (`library-entry-hash`) are useful when investigating suspicious files — check if they appear in any known threat reports.
- **Combine with jq**: Always pipe through `jq` to extract what's relevant rather than dumping raw JSON.
