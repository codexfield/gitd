package transport

import (
	"context"
	"errors"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/packfile"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/utils/ioutil"
	"io"
)

type rpSession struct {
	session
	cmdStatus map[plumbing.ReferenceName]error
	firstErr  error
	unpackErr error
}

var (
	ErrUpdateReference = errors.New("failed to update ref")
)

func (s *rpSession) ReceivePack(ctx context.Context, req *packp.ReferenceUpdateRequest) (*packp.ReportStatus, error) {
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

	if req.Packfile != nil {
		r := ioutil.NewContextReadCloser(ctx, req.Packfile)
		if err := s.writePackfile(r); err != nil {
			s.unpackErr = err
			s.firstErr = err
			return s.reportStatus(), err
		}
	}

	s.updateReferences(req)
	return s.reportStatus(), s.firstErr
}

func (s *rpSession) AdvertisedReferences() (*packp.AdvRefs, error) {
	return s.AdvertisedReferencesContext(context.TODO())
}

func (s *rpSession) AdvertisedReferencesContext(ctx context.Context) (*packp.AdvRefs, error) {
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

	// Git 2.41+ returns a zero-id plus capabilities when an empty
	// repository is being cloned. This skips the existing logic within
	// advrefs_decode.decodeFirstHash, which expects a flush-pkt instead.
	//
	// This logic aligns with plumbing/transport/internal/common/common.go.
	if ar.IsEmpty() {
		return nil, transport.ErrEmptyRemoteRepository
	}

	transport.FilterUnsupportedCapabilities(ar.Capabilities)

	return ar, nil
}

func (s *rpSession) setSupportedCapabilities(c *capability.List) error {
	if err := c.Set(capability.Agent, capability.DefaultAgent()); err != nil {
		return err
	}

	if err := c.Set(capability.OFSDelta); err != nil {
		return err
	}

	return nil
}

func (s *rpSession) Close() error {
	return nil

}

func (s *rpSession) writePackfile(r io.ReadCloser) error {
	if r == nil {
		return nil
	}

	if err := packfile.UpdateObjectStorage(s.storer, r); err != nil {
		_ = r.Close()
		return err
	}

	return r.Close()
}

func (s *rpSession) reportStatus() *packp.ReportStatus {
	if !s.caps.Supports(capability.ReportStatus) {
		return nil
	}

	rs := packp.NewReportStatus()
	rs.UnpackStatus = "ok"

	if s.unpackErr != nil {
		rs.UnpackStatus = s.unpackErr.Error()
	}

	if s.cmdStatus == nil {
		return rs
	}

	for ref, err := range s.cmdStatus {
		msg := "ok"
		if err != nil {
			msg = err.Error()
		}
		status := &packp.CommandStatus{
			ReferenceName: ref,
			Status:        msg,
		}
		rs.CommandStatuses = append(rs.CommandStatuses, status)
	}

	return rs
}

func referenceExists(s storer.ReferenceStorer, n plumbing.ReferenceName) (bool, error) {
	_, err := s.Reference(n)
	if err == plumbing.ErrReferenceNotFound {
		return false, nil
	}

	return err == nil, err
}

func (s *rpSession) updateReferences(req *packp.ReferenceUpdateRequest) {
	for _, cmd := range req.Commands {
		exists, err := referenceExists(s.storer, cmd.Name)
		if err != nil {
			s.setStatus(cmd.Name, err)
			continue
		}

		switch cmd.Action() {
		case packp.Create:
			if exists {
				s.setStatus(cmd.Name, ErrUpdateReference)
				continue
			}

			ref := plumbing.NewHashReference(cmd.Name, cmd.New)
			err := s.storer.SetReference(ref)
			s.setStatus(cmd.Name, err)
		case packp.Delete:
			if !exists {
				s.setStatus(cmd.Name, ErrUpdateReference)
				continue
			}

			err := s.storer.RemoveReference(cmd.Name)
			s.setStatus(cmd.Name, err)
		case packp.Update:
			if !exists {
				s.setStatus(cmd.Name, ErrUpdateReference)
				continue
			}

			ref := plumbing.NewHashReference(cmd.Name, cmd.New)
			err := s.storer.SetReference(ref)
			s.setStatus(cmd.Name, err)
		}
	}
}

func (s *rpSession) setStatus(ref plumbing.ReferenceName, err error) {
	s.cmdStatus[ref] = err
	if s.firstErr == nil && err != nil {
		s.firstErr = err
	}
}
