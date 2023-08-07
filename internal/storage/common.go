package storage

import (
	"cosmossdk.io/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"strings"
)

const (
	ReferenceKey  = "refs/"
	ConfigKey     = "config/"
	ObjectKey     = "objects/"
	ObjectTypeKey = "types/"
)

func parseReference(key string) (*plumbing.Reference, error) {
	if !strings.HasPrefix(key, ReferenceKey) {
		return nil, errors.Wrapf(git.ErrInvalidReference, "keys: %s", key)
	}
	spilts := strings.Split(key, "-")
	if len(spilts) != 2 {
		return nil, errors.Wrapf(git.ErrInvalidReference, "keys: %s", key)
	}

	return plumbing.NewReferenceFromStrings(spilts[0], spilts[1]), nil
}

func buildReferenceKey(name plumbing.ReferenceName) string {
	return ReferenceKey + name.String()
}

func buildObjectsKey(hash plumbing.Hash) string {
	return ObjectKey + hash.String()
}

func buildObjectTypeKey(hash plumbing.Hash) string {
	return ObjectTypeKey + hash.String()
}
