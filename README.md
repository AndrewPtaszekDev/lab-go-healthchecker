# Lab Healthchecker

☞ _Made for fun to learn go concurrency, Elm-style TUIs, and serving ssh apps._

A small SSH-served TUI that checks [my Lab's](https://github.com/AndrewPtaszekDev/lab-showcase) service-health status live.

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

*Right now*, you place services to monitor in `internal/data.go`. Each service pings a
destination (`ping -c 4`); edit the list to change what's watched.

## Internals

* `server.go`: entry point (`StartServer`). Boots the SSH server and, per connection, hands off to a fresh program seeded with the service list.
* `data.go`: the seed list of services to monitor, passed into each session.
* `healthcheck.go`: defines the `healthState` enum (Healthy/Stale/Unhealthy/Unknown) and the `healthcheck` (command + destination); where the command can be set to any 'healthcheck' command)
* `service.go`: a `service` pairs a name with its healthcheck and current state.
  `processHeathchecks` fans the services out to a worker pool and streams results back on a channel.
* `model.go`: state + `Update()`/`View()`. It drives refreshes, drains the results channel from `service.go`, and renders the dashboard.
* `styles.go`: styles used by `model.go`'s `View`.

Flow: `server.go` starts a session → `model.go` triggers a refresh → `service.go` runs
the checks concurrently → results stream back into the model → `View()` (with `styles.go`)
renders them.
