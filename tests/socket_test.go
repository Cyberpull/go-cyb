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

		subscriber.On("DEMO", "/testing", func(ctx *socket.Context) *socket.Output {
			var data string

			if err := ctx.ParseData(&data); err != nil {
				return ctx.Error(err)
			}

			if data != "TestData" {
				return ctx.Error("FAILED")
			}

			return ctx.Success("SUCCESSFUL")
		})

		subscriber.On("UPDATE_DEMO", "/testing", func(ctx *socket.Context) *socket.Output {
			var data string

			if err := ctx.ParseData(&data); err != nil {
				return ctx.Error(err)
			}

			out := ctx.Success("NEW_UPDATE")

			if err = ctx.Update(out); err != nil {
				return ctx.Error(err)
			}

			return ctx.Success("SUCCESSFUL")
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
	resp, err := socket.MakeRequest[string](socketClient, "DEMO", "/testing", "TestData")

	if err != nil {
		t.Fatal(err)
		return
	}

	if resp != "SUCCESSFUL" {
		t.Fatalf(`Expected "SUCCESSFUL", got "%s" instead.`, resp)
	}
}

func TestSocket_SendFailedRequest(t *testing.T) {
	_, err := socket.MakeRequest[string](socketClient, "DEMO", "/testing", "BadTestData")

	if err == nil {
		t.Fatal(`Expected an error`)
		return
	}

	message := err.Error()

	if message != "FAILED" {
		t.Fatalf(`Expected "FAILED", got "%s" instead.`, message)
	}
}

func TestSocket_ReceiveUpdate(t *testing.T) {
	var err error

	errChan := make(chan error, 1)
	defer close(errChan)

	var requestData, updateData string

	socketClient.On("UPDATE_DEMO", "/testing", func(out *socket.Output) {
		errChan <- out.ParseData(&updateData)
	})

	requestData, err = socket.MakeRequest[string](socketClient, "UPDATE_DEMO", "/testing", "TestData")

	if err != nil {
		t.Fatal(err)
		return
	}

	if requestData != "SUCCESSFUL" {
		t.Fatalf(`Expected "SUCCESSFUL", got "%s" instead.`, requestData)
		return
	}

	if err = <-errChan; err != nil {
		t.Fatal(err)
		return
	}

	if updateData != "NEW_UPDATE" {
		t.Fatalf(`Expected "NEW_UPDATE", got "%s" instead.`, updateData)
	}
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
