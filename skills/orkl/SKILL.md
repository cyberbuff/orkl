---
name: orkl
description: Query the ORKL cyber threat intelligence library using the orkl CLI. Use when looking up threat reports, library entries, threat actors, or sources from the ORKL API.
---

# ORKL

Use the `orkl` CLI to query the ORKL cyber threat intelligence library — threat reports, threat actors, and sources.

If not installed: `go install github.com/cyberbuff/orkl/cmd/orkl@latest`

## Global flags

| Flag                   | Default | Description                        |
| ---------------------- | ------- | ---------------------------------- |
| `--timeout <duration>` | `60s`   | Request timeout (e.g. `10s`, `1m`) |

## Library commands

### Get library info

```bash
orkl library-info
```

### Get current library version

```bash
orkl library-version
```

### List library versions

```bash
orkl library-versions [--limit N] [--offset N] [--order asc|desc]
```

### List library entries

```bash
orkl library-entries [--limit N] [--offset N] [--order-by <field>] [--order asc|desc] [--origin pdf|web]
```

### Get a library entry by UUID

```bash
orkl library-entry <uuid>
```

### Get a library entry by SHA1 hash

```bash
orkl library-entry-hash <sha1>
```

### Search the library

```bash
orkl library-search <query> [--full] [--limit N] [--origin pdf|web]
```

| Flag                | Description                                    |
| ------------------- | ---------------------------------------------- |
| `--full`            | Return full entry objects instead of summaries |
| `--limit N`         | Maximum number of results to return            |
| `--origin pdf\|web` | Filter results by origin                       |

Example — search for APT reports with full details:

```bash
orkl library-search "APT29" --full --limit 10
```

## Source commands

### List all sources

```bash
orkl source-entries
```

### Get a source entry by UUID

```bash
orkl source-entry <uuid> [--full]
```

`--full` includes the complete list of reports associated with the source.

## Threat actor commands

### List all threat actors

```bash
orkl ta-entries
```

### Get a threat actor by UUID

```bash
orkl ta-entry <uuid>
```

## Output

All commands print raw JSON to stdout. Pipe through `jq` to format or filter:

```bash
orkl library-info | jq .
orkl library-search "ransomware" --limit 5 | jq '.[].title'
orkl ta-entry <uuid> | jq '{name: .name, aliases: .aliases}'
```
