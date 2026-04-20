# orkl

Go client for the [ORKL](https://orkl.eu) cyber threat intelligence library.

## CLI

### Install

```bash
go install github.com/cyberbuff/orkl/cmd/orkl@latest
```

### Usage

```
orkl [--timeout DURATION] <command> [flags]
```

| Global flag | Default | Description                        |
| ----------- | ------- | ---------------------------------- |
| `--timeout` | `60s`   | Request timeout (e.g. `10s`, `1m`) |

All commands print raw JSON to stdout.

### Commands

| Command                                                                                               | Description                                          |
| ----------------------------------------------------------------------------------------------------- | ---------------------------------------------------- |
| `library-info`                                                                                        | Library metadata (entry count, last update)          |
| `library-version`                                                                                     | Current library version                              |
| `library-versions [--limit N] [--offset N] [--order asc\|desc]`                                       | List library versions                                |
| `library-entries [--limit N] [--offset N] [--order-by field] [--order asc\|desc] [--origin pdf\|web]` | List library entries                                 |
| `library-entry <uuid>`                                                                                | Get a library entry by UUID                          |
| `library-entry-hash <sha1>`                                                                           | Get a library entry by SHA1 hash                     |
| `library-search <query> [--full] [--limit N] [--origin pdf\|web]`                                     | Search the library                                   |
| `source-entries`                                                                                      | List all sources                                     |
| `source-entry <uuid> [--full]`                                                                        | Get a source by UUID (`--full` includes report list) |
| `ta-entries`                                                                                          | List all threat actors                               |
| `ta-entry <uuid>`                                                                                     | Get a threat actor by UUID                           |

### Examples

```bash
orkl library-info
orkl library-search "APT29" --full --limit 10
orkl library-entry 3f4a1b2c-0000-0000-0000-000000000000
orkl ta-entries | jq '.[].main_name'
```

## Library

### Install

```bash
go get github.com/cyberbuff/orkl
```

### Usage

```go
import "github.com/cyberbuff/orkl"

client := orkl.NewClient("", 0) // uses defaults

body, err := client.Get(ctx, "/library/info", nil)

// With query params
params := url.Values{}
params.Set("query", "APT29")
params.Set("limit", "10")
body, err = client.Get(ctx, "/library/search", params)
```

`NewClient(baseURL string, timeout time.Duration)` — pass empty string and 0 to use the defaults (`https://orkl.eu/api/v1`, 30s timeout).

`Get(ctx, path, params)` returns the raw response body as `[]byte` or an error that includes the HTTP status and response body on non-2xx responses.
