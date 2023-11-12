package storage

import (
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	ReferenceKey = "refs/"
	ConfigKey    = "config/"
	ObjectKey    = "objects/"
)

func buildReferenceKey(name plumbing.ReferenceName) string {
	return ReferenceKey + name.String()
}

func buildObjectsKey(hash plumbing.Hash) string {
	return ObjectKey + hash.String()
}
