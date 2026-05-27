# env-diff

Compare two .env files and show the differences between them.

## Features

- Compare any two .env files side by side
- Shows added, removed, modified, and unchanged keys
- Handles quoted values (single and double quotes)
- Ignores comments and empty lines
- Sorted output for easy reading
- Single binary, no dependencies

## Install

```bash
go install github.com/TataneSan/env-diff@latest
```

Or build from source:

```bash
git clone https://github.com/TataneSan/env-diff.git
cd env-diff
go build -o env-diff .
cp env-diff /usr/local/bin/
```

## Usage

```
env-diff <file1> <file2>
```

### Examples

Compare production and staging:
```bash
env-diff .env.production .env.staging
```

Check what's missing compared to the example:
```bash
env-diff .env.example .env
```

## Output

```
  Comparing: .env.production vs .env.staging
  Total keys: 8 | Added: 1 | Removed: 0 | Modified: 2 | Unchanged: 5

  ~ API_KEY
    .env.production: prod-secret-key
    .env.staging: staging-secret-key
  ~ DB_HOST
    .env.production: db.prod.example.com
    .env.staging: db.staging.example.com
  + DEBUG_MODE
    .env.staging: true
```

### Legend

- `+` Key added (present in file2, not in file1)
- `-` Key removed (present in file1, not in file2)
- `~` Key modified (different values)

## Requirements

- Go 1.21+

## License

MIT
