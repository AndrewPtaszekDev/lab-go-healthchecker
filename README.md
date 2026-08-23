# Lab Health

A small SSH-served TUI that checks lab service-health status live.

Wish runs an SSH server, and every connection gets its own dashboard served as a TUI via Bubble Tea

## Run

```sh
go run .
```

The SSH server listens on `0.0.0.0:23234`. Connect with:

```sh
ssh localhost -p 23234
```

Or with Docker:

```sh
docker build -t lab-health .
docker run -p 23234:23234 lab-health
```

## Keys

| Key            | Action           |
| -------------- | ---------------- |
| `↑`/`k`, `↓`/`j` | move cursor    |
| `space` / `r`  | refresh now      |
| `q` / `ctrl+c` | quit             |

Checks also auto-refresh every 25s.

## Configuration

Right now, you place services to monitor in `internal/data.go`. Each service pings a
destination (`ping -c 4`); edit the list to change what's watched.
