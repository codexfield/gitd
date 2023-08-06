package transport

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
)

type GreenfieldTransport struct {
}

var DefaultClient = NewClient()

func InstallGreenfieldTransport() {
	client.InstallProtocol("gnfd", DefaultClient)
}

func NewClient() transport.Transport {
	return &GreenfieldTransport{}
}

func (gt *GreenfieldTransport) NewUploadPackSession(*transport.Endpoint, transport.AuthMethod) (transport.UploadPackSession, error) {
	//TODO implement me
	panic("implement me")
}
func (gt *GreenfieldTransport) NewReceivePackSession(*transport.Endpoint, transport.AuthMethod) (transport.ReceivePackSession, error) {
	//TODO implement me
	panic("implement me")
}
