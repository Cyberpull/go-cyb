package tests

import (
	"os"
	"strings"
	"testing"

	_ "cyberpull.com/go-cyb/env"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"cyberpull.com/go-cyb/errors"
	"cyberpull.com/go-cyb/log"
	"cyberpull.com/go-cyb/socket"
)

type SocketTestSuite struct {
	suite.Suite

	client *socket.Client
	server *socket.Server
}

func (s *SocketTestSuite) registerServerHandlers() socket.ServerHandlerSubscriber {
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

func (s *SocketTestSuite) registerServerAuth() socket.ServerAuthSubscriber {
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

func (s *SocketTestSuite) serverClientInit() socket.ServerClientInitHandler {
	return func(updater *socket.ServerClientUpdater) (err error) {
		updater.Update("INIT_UPDATE_DEMO", "/testing", "TestData::ClientInit")
		return
	}
}

func (s *SocketTestSuite) registerClientAuth() socket.ClientAuthSubscriber {
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

func (s *SocketTestSuite) clientUpdate() socket.ClientUpdateSubscriber {
	return func(collection *socket.ClientUpdateHandlerCollection) (err error) {
		collection.On("INIT_UPDATE_DEMO", "/testing", func(out *socket.Output) {
			var data string

			if err := out.ParseData(&data); err != nil {
				return
			}

			log.Printfln("Client Init Update: %s", data)
		})

		return
	}
}

func (s *SocketTestSuite) startServer() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.From(r)
		}
	}()

	s.server = socket.NewServer(socket.ServerOptions{
		Host: os.Getenv("SOCKET_SERVER_HOST"),
		Port: os.Getenv("SOCKET_SERVER_PORT"),
		Name: "Socket Testing Server",
	})

	s.server.Auth(
		s.registerServerAuth(),
	)

	s.server.ClientInit(
		s.serverClientInit(),
	)

	s.server.Handlers(
		s.registerServerHandlers(),
	)

	errChan := make(chan error)

	go s.server.Listen(errChan)

	err = <-errChan

	return
}

func (s *SocketTestSuite) startClient() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.From(r)
		}
	}()

	s.client = socket.NewClient(socket.ClientOptions{
		ServerHost: os.Getenv("SOCKET_SERVER_HOST"),
		ServerPort: os.Getenv("SOCKET_SERVER_PORT"),
		Name:       "Socket Testing Client",
	})

	s.client.Auth(
		s.registerClientAuth(),
	)

	s.client.Update(
		s.clientUpdate(),
	)

	errChan := make(chan error)

	go s.client.Start(errChan)

	err = <-errChan

	return
}

func (s *SocketTestSuite) SetupSuite() {
	// Start Socket Server
	require.NoError(s.T(), s.startServer())

	// Start Socket Client
	require.NoError(s.T(), s.startClient())
}

func (s *SocketTestSuite) TearDownSuite() {
	// Stop Socket Client
	require.NoError(s.T(), s.client.Stop())

	// Stop Socket Server
	require.NoError(s.T(), s.server.Stop())
}

func (s *SocketTestSuite) TestSendSuccessfulRequest() {
	resp, err := socket.MakeRequest[string](s.client, "DEMO", "/testing", "TestData")
	require.NoError(s.T(), err)

	assert.Equal(s.T(), "SUCCESSFUL", resp)
}

func (s *SocketTestSuite) TestSendFailedRequest() {
	_, err := socket.MakeRequest[string](s.client, "DEMO", "/testing", "BadTestData")
	require.Error(s.T(), err)

	assert.EqualError(s.T(), err, "FAILED")
}

func (s *SocketTestSuite) TestReceiveUpdate() {
	var err error

	errChan := make(chan error, 1)
	defer close(errChan)

	var requestData, updateData string

	s.client.On("UPDATE_DEMO", "/testing", func(out *socket.Output) {
		errChan <- out.ParseData(&updateData)
	})

	requestData, err = socket.MakeRequest[string](s.client, "UPDATE_DEMO", "/testing", "TestData")
	require.NoError(s.T(), err)

	assert.Equal(s.T(), "SUCCESSFUL", requestData)

	err = <-errChan
	require.NoError(s.T(), err)

	assert.Equal(s.T(), "NEW_UPDATE", updateData)
}

/********************************************/

func TestSocket(t *testing.T) {
	suite.Run(t, new(SocketTestSuite))
}
