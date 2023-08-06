package transport

import (
	"cosmossdk.io/errors"
	"fmt"
	"gitd/internal/storage"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	"os"
	"strconv"
)

type GnfdTransport struct {
	gnfdStorage *storage.GnfdStorage
	asClient    bool
}

var DefaultClient = NewClient()

const (
	GnfdProtocol = "gnfd"

	EnvChainID    = "GREENFIELD_CHAIN_ID"
	EnvPrivateKey = "GREENFIELD_PRIVATE_KEY"
)

func InstallGreenfieldTransport() {
	client.InstallProtocol(GnfdProtocol, DefaultClient)
}

func NewClient() transport.Transport {
	return &GnfdTransport{
		asClient: true,
	}
}

func NewServer() transport.Transport {
	return &GnfdTransport{
		asClient: false,
	}
}

func (t *GnfdTransport) NewUploadPackSession(ep *transport.Endpoint, auth transport.AuthMethod) (transport.UploadPackSession, error) {
	s, err := t.LoadStorage(ep)
	if err != nil {
		return nil, err
	}
	return &upSession{
		session: session{storer: s, asClient: false},
	}, nil
}
func (t *GnfdTransport) NewReceivePackSession(ep *transport.Endpoint, auth transport.AuthMethod) (transport.ReceivePackSession, error) {
	s, err := t.LoadStorage(ep)
	if err != nil {
		return nil, err
	}
	return &rpSession{
		session:   session{storer: s, asClient: false},
		cmdStatus: map[plumbing.ReferenceName]error{},
	}, nil
}

func (t *GnfdTransport) LoadStorage(endpoint *transport.Endpoint) (storer.Storer, error) {
	if t.gnfdStorage != nil {
		return t.gnfdStorage, nil
	}

	// TODO: refine the config
	chainID := os.Getenv(EnvChainID)
	if chainID == "" {
		panic(fmt.Sprintf("Please set the environment variable: %s", EnvChainID))
	}

	privateKey := os.Getenv(EnvPrivateKey)
	if privateKey == "" {
		panic(fmt.Sprintf("Please set the enviroment variable: %s", EnvPrivateKey))
	}

	// Endpoint scheme:
	// if not loaded, init a greenfield client
	if endpoint.Protocol != GnfdProtocol {
		return nil, transport.ErrRepositoryNotFound
	}

	rpcAddress := "https://" + endpoint.Host + ":" + strconv.Itoa(endpoint.Port)
	newStorage, err := storage.NewStorage(chainID, rpcAddress, privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "New greenfield storage failed.")
	}
	return newStorage, nil
}
