package internal

import (
	"fmt"
	"os/exec"
	"time"
)

const workers = 10

type service struct {

	name           string
	lastState      healthState
	state                        healthState
	healthcheck    healthcheck
	atStateSince   time.Time
	atStateElasped time.Duration
}

func (s service) init(name string, healthcheck healthcheck) service {
	return service{
		name:           name,
		state:            Unknown,
		lastState:      Unknown,
		atStateSince:   time.Now(),
		atStateElasped: time.Duration(0),
		healthcheck:    healthcheck,
	}
}

// if  diff from LastState; update AtStateSince (i.e. at THIS/CURRENT state since
func (s *service) updateAtStateSince() {
	if s.state != s.lastState {
		s.atStateSince = time.Now()
		s.lastState = s.state
	}
	s.atStateElasped = time.Since(s.atStateSince)
}

func (s service) processCommand() {
	// []string -- --> exec.cmd type
	h := s.healthcheck
	args := append([]string{}, h.cmd[1:]...)
	args = append(args, h.dest)
	command := exec.Command(h.cmd[0], args...)

	_, err := command.Output()

	if err != nil {
		s.state = Unhealthy
		s.updateAtStateSince()

		h.msg = fmt.Sprintf("destination: %s unhealthy. Error: %s", h.dest, err)
		return
	}

	s.state = Healthy
	s.updateAtStateSince()
	h.msg = fmt.Sprintf("destination: %s OK (healthy)", h.dest)
}

// Needs to be converted to PER SERVICE proceessing, also paralel such that results trickle back to main independently
func processHeathchecks(services []service) []service {
	in := make(chan service, workers)
	out := make(chan service, workers)

	// one goroutine PER worker (queue workers to unload in)
	for i := 0; i < workers; i++ {
		go func() {
			for s := range in { //waits for in to have something; stops when closed AND empty

				s.processCommand()
				out <- s
			}
		}()
	}

	// one goroutine
	go func() {
		for _, s := range services { // # of health checks determines how many get pulled out of in by above workers
			in <- s
	}
		close(in) // in can still be READ
	}()

	results := make([]service, 0, len(services))
		for range services {
			s := <-out
			results = append(results, s) //waits, but grabs exactly 'healthchecks' number of vals
	}

	return results
}
