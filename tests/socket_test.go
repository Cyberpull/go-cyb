package tests

import (
	"os"
	"testing"
	"time"

	_ "cyberpull.com/go-cyb/env"

	"cyberpull.com/go-cyb/errors"
	"cyberpull.com/go-cyb/log"
	"cyberpull.com/go-cyb/socket"
)

var (
	socketClient *socket.Client
	socketServer *socket.Server
)

// Server Handlers
func socketRegisterHandlers() socket.ServerHandlerSubscriber {
	return func(subscriber *socket.ServerHandlerCollection) (err error) {
		log.Println("Registering test handlers...")

		subscriber.On("SUCCESS", "/testing", func(ctx *socket.Context) *socket.Output {
			return ctx.Success("Success Response")
		})

		subscriber.On("ERROR", "/testing", func(ctx *socket.Context) *socket.Output {
			return ctx.Error("Error Response")
		})

		return
	}
}

func TestSocket_StartServer(t *testing.T) {
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = errors.From(r)
		}

		if err != nil {
			t.Fatal(err)
		}
	}()

	socketServer = socket.NewServer(socket.ServerOptions{
		Host: os.Getenv("SOCKET_SERVER_HOST"),
		Port: os.Getenv("SOCKET_SERVER_PORT"),
		Name: "Socket Testing Server",
	})

	socketServer.Handlers(
		socketRegisterHandlers(),
	)

	go socketServer.Listen()

	time.Sleep(time.Second)

	err = socketServer.EnsureListening()
}

func TestSocket_StartClient(t *testing.T) {
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = errors.From(r)
		}

		if err != nil {
			t.Fatal(err)
		}
	}()

	socketClient = socket.NewClient(socket.ClientOptions{
		ServerHost: os.Getenv("SOCKET_SERVER_HOST"),
		ServerPort: os.Getenv("SOCKET_SERVER_PORT"),
		Name:       "Socket Testing Client",
	})

	go socketClient.Start()

	time.Sleep(time.Second)

	err = socketClient.EnsureStarted()
}

// Conclusion of socket tests

func TestSocket_StopClient(t *testing.T) {
	if err := socketClient.Stop(); err != nil {
		t.Fatal(err)
	}
}

func TestSocket_StopServer(t *testing.T) {
	if err := socketServer.Stop(); err != nil {
		t.Fatal(err)
	}
}
