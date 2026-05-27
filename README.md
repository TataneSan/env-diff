# env-diff

Compare environment variables between two contexts. Supports comparing .env files, current shell environment, and command output.

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

```
env-diff file <fileA> <fileB>     Compare two .env files
env-diff current <file>           Compare current env with a .env file
env-diff cmd <cmdA> <cmdB>        Compare output of two commands
```

## Examples

Compare two .env files:
```bash
env-diff file .env .env.production
```

Compare current shell environment with a .env file:
```bash
env-diff current .env
```

Compare environment from two commands:
```bash
env-diff cmd 'printenv APP_*' 'cat .env | grep APP_'
```

## Output

The tool outputs a color-coded diff:
- `+` green: variable added in the second context
- `-` red: variable removed (only in first context)
- `~` yellow: variable modified (different values)

## Requirements

- Linux or macOS
- Go 1.21+

## License

MIT
