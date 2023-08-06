package storage

import (
	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage"
)

type Storage struct {
	client client.Client
}

func NewStorage(chainID, endpoint, privateKey string) (*Storage, error) {
	account, err := types.NewAccountFromPrivateKey("gitd", privateKey)
	if err != nil {
		return nil, err
	}
	gnfdClient, err := client.New(chainID, endpoint, client.Option{DefaultAccount: account})
	if err != nil {
		return nil, err
	}
	return &Storage{
		client: gnfdClient,
	}, nil
}

func (s *Storage) NewEncodedObject() plumbing.EncodedObject {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SetEncodedObject(object plumbing.EncodedObject) (plumbing.Hash, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) EncodedObject(objectType plumbing.ObjectType, hash plumbing.Hash) (plumbing.EncodedObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) IterEncodedObjects(objectType plumbing.ObjectType) (storer.EncodedObjectIter, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) HasEncodedObject(hash plumbing.Hash) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) EncodedObjectSize(hash plumbing.Hash) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SetReference(reference *plumbing.Reference) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) CheckAndSetReference(new, old *plumbing.Reference) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Reference(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) IterReferences() (storer.ReferenceIter, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) RemoveReference(name plumbing.ReferenceName) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) CountLooseRefs() (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) PackRefs() error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SetShallow(hashes []plumbing.Hash) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Shallow() ([]plumbing.Hash, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SetIndex(index *index.Index) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Index() (*index.Index, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Config() (*config.Config, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SetConfig(config *config.Config) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Module(name string) (storage.Storer, error) {
	//TODO implement me
	panic("implement me")
}
