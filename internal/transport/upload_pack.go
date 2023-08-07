package transport

import (
	"context"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type upSession struct {
	session
}

func (s *upSession) UploadPack(ctx context.Context, request *packp.UploadPackRequest) (*packp.UploadPackResponse, error) {
	//TODO implement me
	panic("implement me")
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

func setReferences(s storer.Storer, ar *packp.AdvRefs) error {
	//TODO: add peeled references.
	iter, err := s.IterReferences()
	if err != nil {
		return err
	}
	return iter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() != plumbing.HashReference {
			return nil
		}

		ar.References[ref.Name().String()] = ref.Hash()
		return nil
	})
}

func setHEAD(s storer.Storer, ar *packp.AdvRefs) error {
	ref, err := s.Reference(plumbing.HEAD)
	if err == plumbing.ErrReferenceNotFound {
		return nil
	}

	if err != nil {
		return err
	}

	if ref.Type() == plumbing.SymbolicReference {
		if err := ar.AddReference(ref); err != nil {
			return nil
		}

		ref, err = storer.ResolveReference(s, ref.Target())
		if err == plumbing.ErrReferenceNotFound {
			return nil
		}

		if err != nil {
			return err
		}
	}

	if ref.Type() != plumbing.HashReference {
		return plumbing.ErrInvalidType
	}

	h := ref.Hash()
	ar.Head = &h

	return nil
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
