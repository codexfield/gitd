package storage

import (
	"context"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/sync/semaphore"
	"sync"
)

const (
	ReferenceKey = "refs/"
	ConfigKey    = "config/"
	ObjectKey    = "objects/"
)

type GnfdStorageContext struct {
	context.Context
	Sem *semaphore.Weighted
	Wg  sync.WaitGroup
}

func buildReferenceKey(name plumbing.ReferenceName) string {
	return ReferenceKey + name.String()
}

func buildObjectsKey(hash plumbing.Hash) string {
	return ObjectKey + hash.String()
}
