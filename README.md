# env-diff

Compare two .env files.

## Install

```bash
pip install .
```

## Usage

```
env-diff left.env right.env [--mask-values] [--shell-export] [--require-identical] [--no-conflicts] [--json]
```

## Example

```bash
env-diff .env.dev .env.prod
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0    | success |
| 1    | I/O or CLI error |
| 2    | gate condition not met |

## License

MIT
