package tests

import (
	"os"
	"strings"
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
func socketRegisterServerHandlers() socket.ServerHandlerSubscriber {
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

// Server Handlers
func socketRegisterServerAuth() socket.ServerAuthSubscriber {
	return func(ref *socket.ServerClientRef) (err error) {
		if _, err = ref.WriteStringln("TOKEN:"); err != nil {
			return
		}

		var data string

		if data, err = ref.ReadString('\n'); err != nil {
			return
		}

		data = strings.TrimSpace(data)
		data = strings.ToLower(data)

		if data != "testing" {
			err = errors.Newf(`Expected "testing", got "" instead.`, 500, data)
			return
		}

		_, err = ref.WriteStringln("SUCCESS")

		return
	}
}

func socketRegisterClientAuth() socket.ClientAuthSubscriber {
	return func(ref *socket.ClientRef) (err error) {
		var data string

		if data, err = ref.ReadString('\n'); err != nil {
			return
		}

		if data = strings.TrimSpace(data); data != "TOKEN:" {
			err = errors.Newf(`Expected "TOKEN:", got "%s" instead.`, 500, data)
			return
		}

		if _, err = ref.WriteStringln("testing"); err != nil {
			return
		}

		if data, err = ref.ReadString('\n'); err != nil {
			return
		}

		if data = strings.TrimSpace(data); data != "SUCCESS" {
			err = errors.New("Authorization failed", 403)
		}

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

	socketServer.Auth(
		socketRegisterServerAuth(),
	)

	socketServer.Handlers(
		socketRegisterServerHandlers(),
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

	socketClient.Auth(
		socketRegisterClientAuth(),
	)

	go socketClient.Start()

	time.Sleep(time.Second)

	err = socketClient.EnsureStarted()
}

func TestSocket_SendSuccessfulRequest(t *testing.T) {
	//
}

func TestSocket_SendFailedRequest(t *testing.T) {
	//
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
