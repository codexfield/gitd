package storage

import (
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	ReferenceKey  = "refs/"
	ConfigKey     = "config/"
	ObjectKey     = "objects/"
	ObjectTypeKey = "types/"
)

func buildReferenceKey(name plumbing.ReferenceName) string {
	return ReferenceKey + name.String()
}

func buildObjectsKey(t plumbing.ObjectType, hash plumbing.Hash) string {
	return ObjectKey + t.String() + "/" + hash.String()
}

func buildObjectTypeKey(hash plumbing.Hash) string {
	return ObjectTypeKey + hash.String()
}
