package transport

import (
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type session struct {
	storer   storer.Storer
	caps     *capability.List
	asClient bool
}
