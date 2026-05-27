# env-diff

Compare two `.env` files and see exactly what changed.

```
$ env-diff .env.staging .env.production

  Comparing: .env.staging  →  .env.production
  ────────────────────────────────────────────────────────
  1 keys in .env.staging (removed)
  1 keys in .env.production (added)
  3 keys changed
  2 keys unchanged
  ────────────────────────────────────────────────────────

  REMOVED (only in first file)
    - DEBUG = true

  ADDED (only in second file)
    + CACHE_TTL = 3600

  CHANGED
    ~ DB_HOST
      - localhost
      + production-db.example.com
    ~ DB_NAME
      - myapp
      + myapp_prod
    ~ SECRET_KEY
      - abc123
      + xyz789

  UNCHANGED (2 keys)
     DB_PORT = 5432
     DB_USER = admin
```

## Install

```bash
go install github.com/TataneSan/env-diff@latest
```

Or build from source:

```bash
git clone https://github.com/TataneSan/env-diff.git
cd env-diff
go build -o env-diff .
```

## Usage

```bash
# Compare two env files
env-diff .env.staging .env.production

# Quiet mode — only show differences (hide unchanged)
env-diff -q .env.staging .env.production

# JSON output (pipe-friendly)
env-diff -f json .env.staging .env.production

# JSON + quiet
env-diff -f json -q .env.staging .env.production
```

## Flags

| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--format` | `-f` | `table` | Output format: `table` or `json` |
| `--quiet` | `-q` | `false` | Only show differences, hide unchanged keys |

## JSON Output

```bash
$ env-diff -f json .env.dev .env.prod
```

```json
{
  "changed": {
    "DB_HOST": {
      "old": "localhost",
      "new": "production-db.example.com"
    }
  },
  "first_file": ".env.dev",
  "only_in_first": {
    "DEBUG": "true"
  },
  "only_in_second": {
    "CACHE_TTL": "3600"
  },
  "second_file": ".env.prod",
  "summary": {
    "added": 1,
    "changed": 1,
    "removed": 1,
    "unchanged": 1
  },
  "unchanged": {
    "DB_PORT": "5432"
  }
}
```

## Features

- Colorized terminal output with symbols: `-` removed, `+` added, `~` changed
- JSON output for scripting and CI pipelines
- Handles quoted values (`"value"` and `'value'`)
- Skips comments (`#`) and blank lines
- Sorted output for deterministic comparison
- Quiet mode (`-q`) to focus only on differences
- Respects `NO_COLOR` environment variable

## Docker

```bash
docker run --rm -v $(pwd):/data ghcr.io/tatanesan/env-diff:latest /data/.env.dev /data/.env.prod
```

## License

MIT
