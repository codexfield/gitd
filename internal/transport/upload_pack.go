package transport

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/packfile"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/revlist"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/utils/ioutil"
)

type upSession struct {
	session
}

func (s *upSession) UploadPack(ctx context.Context, req *packp.UploadPackRequest) (*packp.UploadPackResponse, error) {
	if req.IsEmpty() {
		return nil, transport.ErrEmptyUploadPackRequest
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	if s.caps == nil {
		s.caps = capability.NewList()
		if err := s.setSupportedCapabilities(s.caps); err != nil {
			return nil, err
		}
	}

	if err := s.checkSupportedCapabilities(req.Capabilities); err != nil {
		return nil, err
	}

	s.caps = req.Capabilities

	if len(req.Shallows) > 0 {
		return nil, fmt.Errorf("shallow not supported")
	}

	objs, err := s.objectsToUpload(req)
	if err != nil {
		return nil, err
	}

	pr, pw := ioutil.Pipe()
	e := packfile.NewEncoder(pw, s.storer, false)
	go func() {
		// TODO: plumb through a pack window.
		_, err := e.Encode(objs, 10)
		pw.CloseWithError(err)
	}()

	return packp.NewUploadPackResponseWithPackfile(req,
		ioutil.NewContextReadCloser(ctx, pr),
	), nil
}

func (s *upSession) objectsToUpload(req *packp.UploadPackRequest) ([]plumbing.Hash, error) {
	haves, err := revlist.Objects(s.storer, req.Haves, nil)
	if err != nil {
		return nil, err
	}

	return revlist.Objects(s.storer, req.Wants, haves)
}

func (s *upSession) AdvertisedReferences() (*packp.AdvRefs, error) {
	return s.AdvertisedReferencesContext(context.TODO())
}

func (s *upSession) AdvertisedReferencesContext(ctx context.Context) (*packp.AdvRefs, error) {
	ar := packp.NewAdvRefs()

	if err := s.setSupportedCapabilities(ar.Capabilities); err != nil {
		return nil, err
	}
	s.caps = ar.Capabilities
	if err := setReferences(s.storer, ar); err != nil {
		return nil, err
	}
	if err := setHEAD(s.storer, ar); err != nil {
		return nil, err
	}
	if s.asClient && len(ar.References) == 0 {
		return nil, transport.ErrEmptyRemoteRepository
	}
	return ar, nil
}

func (*upSession) setSupportedCapabilities(c *capability.List) error {
	if err := c.Set(capability.Agent, capability.DefaultAgent()); err != nil {
		return err
	}

	if err := c.Set(capability.OFSDelta); err != nil {
		return err
	}

	return nil
}

func (s *upSession) Close() error {
	return nil
}
