package internal

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

func bubbleteaMiddleware(seedServices []service) wish.Middleware {
	teaHandler := func(s ssh.Session) *tea.Program {
		// EVERYTHING FROM HERE AND ON IS PER-CONNECTION
		m := initialModel(seedServices)
		p := tea.NewProgram(m, bubbletea.MakeOptions(s)...)
		return p
	}
	return bubbletea.MiddlewareWithProgramHandler(teaHandler)
}

func StartServer() {
	// SERVER/PROCESS-SCOPE until teaHandler is run
	s, err := wish.NewServer(
		ssh.AllocatePty(), // for reading input from user
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbleteaMiddleware(seedServices), // seed services data from data.go
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%s", host, port)
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-done
	log.Print("Stopping SSH server\n")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Fatal(err)
	}
}
