package internal

type healthState int
const (
	Healthy healthState = iota
	Stale // not UNHEALTHY, but hasnt been recently checked (checked at least once thou)
	Unhealthy
	Unknown // Never been checked
)

var names = map[healthState]string{
		Healthy: "Healthy",
		Stale: "Stale",
		Unhealthy: "Unhealthy",
		Unknown: "Unknown",
	}

func (hs healthState) String() string {
	return names[hs]
}

func (h healthcheck) init(dest string) healthcheck {
	//assumes the command is always ping; for now
	ping := []string{"ping", "-c", "4"}
	return healthcheck{
		cmd: ping,
		dest: dest,
	}
}

type healthcheck struct {
	cmd []string
	dest string
	msg string
}
