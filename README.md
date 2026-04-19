# logslice

A CLI tool to filter and slice structured log files by time range or field value.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build -o logslice .
```

## Usage

```bash
# Filter logs by time range
logslice --file app.log --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z"

# Filter by field value
logslice --file app.log --field level=error

# Combine filters
logslice --file app.log --from "2024-01-15T08:00:00Z" --field service=api

# Read from stdin
cat app.log | logslice --field level=warn
```

### Flags

| Flag | Description |
|------|-------------|
| `--file` | Path to the log file (defaults to stdin) |
| `--from` | Start of time range (RFC3339) |
| `--to` | End of time range (RFC3339) |
| `--field` | Filter by field in `key=value` format |
| `--time-key` | JSON key used for timestamps (default: `time`) |
| `--output` | Output format: `json` or `pretty` (default: `json`) |

## Supported Formats

- JSON / NDJSON (newline-delimited JSON)
- Logfmt

## License

MIT © 2024 [Your Name](https://github.com/yourusername)