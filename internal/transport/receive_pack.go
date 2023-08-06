package transport

import (
	"context"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
)

type rpSession struct {
	session
	cmdStatus map[plumbing.ReferenceName]error
	firstErr  error
	unpackErr error
}

func (s *rpSession) ReceivePack(ctx context.Context, request *packp.ReferenceUpdateRequest) (*packp.ReportStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (s *rpSession) AdvertisedReferences() (*packp.AdvRefs, error) {
	//TODO implement me
	panic("implement me")
}

func (s *rpSession) AdvertisedReferencesContext(ctx context.Context) (*packp.AdvRefs, error) {
	//TODO implement me
	panic("implement me")
}

func (s *rpSession) Close() error {
	//TODO implement me
	panic("implement me")
}
