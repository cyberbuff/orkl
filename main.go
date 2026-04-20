package orkl

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"
)

func RunCLI(args []string, stdout, stderr io.Writer) int {
	return runWithBaseURL(context.Background(), args, stdout, stderr, defaultBaseURL)
}

func runWithBaseURL(ctx context.Context, args []string, stdout, stderr io.Writer, baseURL string) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	global := flag.NewFlagSet("orkl", flag.ContinueOnError)
	global.SetOutput(stderr)
	timeout := global.Duration("timeout", 60*time.Second, "request timeout")

	if err := global.Parse(args); err != nil {
		return 2
	}

	rest := global.Args()
	if len(rest) == 0 {
		printUsage(stderr)
		return 2
	}

	client := NewClient(baseURL, *timeout)
	switch rest[0] {
	case "library-info":
		body, err := client.Get(ctx, "/library/info", nil)
		return printRequest(stdout, stderr, body, err)
	case "library-version":
		body, err := client.Get(ctx, "/library/version", nil)
		return printRequest(stdout, stderr, body, err)
	case "library-versions":
		fs := flag.NewFlagSet("library-versions", flag.ContinueOnError)
		fs.SetOutput(stderr)
		limit := fs.Int("limit", 0, "maximum versions to return")
		offset := fs.Int("offset", 0, "number of versions to skip")
		order := fs.String("order", "", "sort order: asc or desc")
		if err := fs.Parse(rest[1:]); err != nil {
			return 2
		}
		params := queryValues(map[string]string{
			"limit":  optionalInt(*limit),
			"offset": optionalInt(*offset),
			"order":  optionalString(*order),
		})
		body, err := client.Get(ctx, "/library/version/entries", params)
		return printRequest(stdout, stderr, body, err)
	case "library-entries":
		fs := flag.NewFlagSet("library-entries", flag.ContinueOnError)
		fs.SetOutput(stderr)
		limit := fs.Int("limit", 0, "maximum entries to return")
		offset := fs.Int("offset", 0, "number of entries to skip")
		orderBy := fs.String("order-by", "", "sort field")
		order := fs.String("order", "", "sort order: asc or desc")
		origin := fs.String("origin", "", "entry origin: pdf or web")
		if err := fs.Parse(rest[1:]); err != nil {
			return 2
		}
		params := queryValues(map[string]string{
			"limit":    optionalInt(*limit),
			"offset":   optionalInt(*offset),
			"order_by": optionalString(*orderBy),
			"order":    optionalString(*order),
			"origin":   optionalString(*origin),
		})
		body, err := client.Get(ctx, "/library/entries", params)
		return printRequest(stdout, stderr, body, err)
	case "library-entry":
		if len(rest) < 2 {
			_, _ = fmt.Fprintln(stderr, "library-entry requires a uuid")
			return 2
		}
		body, err := client.Get(ctx, "/library/entry/"+url.PathEscape(rest[1]), nil)
		return printRequest(stdout, stderr, body, err)
	case "library-entry-hash":
		if len(rest) < 2 {
			_, _ = fmt.Fprintln(stderr, "library-entry-hash requires a sha1 hash")
			return 2
		}
		body, err := client.Get(ctx, "/library/entry/sha1/"+url.PathEscape(rest[1]), nil)
		return printRequest(stdout, stderr, body, err)
	case "library-search":
		if len(rest) < 2 {
			_, _ = fmt.Fprintln(stderr, "library-search requires a query")
			return 2
		}
		fs := flag.NewFlagSet("library-search", flag.ContinueOnError)
		fs.SetOutput(stderr)
		full := fs.Bool("full", false, "return full entries")
		limit := fs.Int("limit", 0, "maximum results to return")
		origin := fs.String("origin", "", "result origin: pdf or web")
		if err := fs.Parse(rest[2:]); err != nil {
			return 2
		}
		params := queryValues(map[string]string{
			"query":  rest[1],
			"full":   optionalBool(*full),
			"limit":  optionalInt(*limit),
			"origin": optionalString(*origin),
		})
		body, err := client.Get(ctx, "/library/search", params)
		return printRequest(stdout, stderr, body, err)
	case "source-entries":
		body, err := client.Get(ctx, "/source/entries", nil)
		return printRequest(stdout, stderr, body, err)
	case "source-entry":
		if len(rest) < 2 {
			_, _ = fmt.Fprintln(stderr, "source-entry requires a uuid")
			return 2
		}
		fs := flag.NewFlagSet("source-entry", flag.ContinueOnError)
		fs.SetOutput(stderr)
		full := fs.Bool("full", false, "include full report list")
		if err := fs.Parse(rest[2:]); err != nil {
			return 2
		}
		params := queryValues(map[string]string{
			"full": optionalBool(*full),
		})
		body, err := client.Get(ctx, "/source/entry/"+url.PathEscape(rest[1]), params)
		return printRequest(stdout, stderr, body, err)
	case "ta-entries":
		body, err := client.Get(ctx, "/ta/entries", nil)
		return printRequest(stdout, stderr, body, err)
	case "ta-entry":
		if len(rest) < 2 {
			_, _ = fmt.Fprintln(stderr, "ta-entry requires a uuid")
			return 2
		}
		body, err := client.Get(ctx, "/ta/entry/"+url.PathEscape(rest[1]), nil)
		return printRequest(stdout, stderr, body, err)
	default:
		_, _ = fmt.Fprintf(stderr, "unknown command: %s\n", rest[0])
		printUsage(stderr)
		return 2
	}
}

func printRequest(stdout, stderr io.Writer, body []byte, err error) int {
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 1
	}
	_, writeErr := stdout.Write(body)
	if writeErr != nil {
		_, _ = fmt.Fprintln(stderr, writeErr)
		return 1
	}
	return 0
}

func queryValues(items map[string]string) url.Values {
	values := url.Values{}
	for key, value := range items {
		if value != "" {
			values.Set(key, value)
		}
	}
	return values
}

func optionalString(v string) string {
	return v
}

func optionalInt(v int) string {
	if v == 0 {
		return ""
	}
	return strconv.Itoa(v)
}

func optionalBool(v bool) string {
	if !v {
		return ""
	}
	return "true"
}

func printUsage(stderr io.Writer) {
	_, _ = fmt.Fprintln(stderr, "Usage: orkl [--base-url URL] [--timeout DURATION] <command> [flags]")
	_, _ = fmt.Fprintln(stderr, "Commands: library-info, library-version, library-versions, library-entries, library-entry, library-entry-hash, library-search, source-entries, source-entry, ta-entries, ta-entry")
}
