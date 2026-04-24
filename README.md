# portwatch

Lightweight daemon that monitors open ports and alerts on unexpected changes with configurable rules.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a config file:

```bash
portwatch --config /etc/portwatch/config.yaml
```

Example `config.yaml`:

```yaml
interval: 30s
alert:
  method: log
  path: /var/log/portwatch.log
rules:
  - port: 22
    allow: true
  - port: 80
    allow: true
  - port: 443
    allow: true
  - port: "*"
    allow: false
    notify: true
```

portwatch will scan open ports at the defined interval and trigger an alert whenever a port outside your allowed rules is detected.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `./config.yaml` | Path to config file |
| `--interval` | `30s` | Override scan interval |
| `--verbose` | `false` | Enable verbose logging |

## License

MIT © [yourusername](https://github.com/yourusername)