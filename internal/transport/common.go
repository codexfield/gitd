package transport

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type session struct {
	storer   storer.Storer
	caps     *capability.List
	asClient bool
}

func (s *session) checkSupportedCapabilities(cl *capability.List) error {
	for _, c := range cl.All() {
		if !s.caps.Supports(c) {
			return fmt.Errorf("unsupported capability: %s", c)
		}
	}

	return nil
}
