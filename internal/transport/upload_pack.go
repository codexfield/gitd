package transport

import (
	"context"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
)

type upSession struct {
	session
}

func (s *upSession) UploadPack(ctx context.Context, request *packp.UploadPackRequest) (*packp.UploadPackResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *upSession) AdvertisedReferences() (*packp.AdvRefs, error) {
	//TODO implement me
	panic("implement me")
}

func (s *upSession) AdvertisedReferencesContext(ctx context.Context) (*packp.AdvRefs, error) {
	//TODO implement me
	panic("implement me")
}

func (s *upSession) Close() error {
	//TODO implement me
	panic("implement me")
}
