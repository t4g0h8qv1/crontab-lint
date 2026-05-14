# crontab-lint

Static analyzer and validator for crontab expressions with human-readable explanations.

## Installation

```bash
go install github.com/yourusername/crontab-lint@latest
```

## Usage

Validate a crontab expression directly from the command line:

```bash
$ crontab-lint "*/5 * * * *"
✔ Valid expression
  → Runs every 5 minutes
```

```bash
$ crontab-lint "0 25 * * *"
✗ Invalid expression
  → Field 'hour' value 25 is out of range (0-23)
```

You can also lint an entire crontab file:

```bash
$ crontab-lint --file /etc/crontab
Line 4: ✔ Valid   → Runs at midnight every day
Line 7: ✗ Invalid → Field 'minute' value 61 is out of range (0-59)
```

### Flags

| Flag | Description |
|------|-------------|
| `--file` | Path to a crontab file to validate |
| `--explain` | Print a human-readable explanation for each expression |
| `--json` | Output results in JSON format |

### Example with `--explain`

```bash
$ crontab-lint --explain "0 9 * * 1-5"
✔ Valid expression
  → Runs at 09:00 AM, Monday through Friday
```

## Contributing

Contributions are welcome. Please open an issue or submit a pull request.

## License

MIT © yourusername