package storage

import (
	"context"
	"fmt"
	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage"
)

type GnfdStorage struct {
	gnfdClient client.Client
}

func NewStorage(chainID, rpcAddress, privateKey string) (*GnfdStorage, error) {
	account, err := types.NewAccountFromPrivateKey("gitd", privateKey)
	if err != nil {
		fmt.Println("New storage error: ", err)
		return nil, err
	}
	gnfdClient, err := client.New(chainID, rpcAddress, client.Option{DefaultAccount: account})
	if err != nil {
		fmt.Println("New storage error: ", err)
		return nil, err
	}

	block, err := gnfdClient.GetLatestBlock(context.Background())
	if err != nil {
		fmt.Println("New storage error: ", err)
		return nil, err
	}
	fmt.Println("New Greenfield storage success, chainID: ", block.ChainID, "height: ", block.Height)
	return &GnfdStorage{
		gnfdClient: gnfdClient,
	}, nil
}

func (s *GnfdStorage) NewEncodedObject() plumbing.EncodedObject {
	return &plumbing.MemoryObject{}
}

func (s *GnfdStorage) SetEncodedObject(object plumbing.EncodedObject) (plumbing.Hash, error) {
	panic("implement me")
}

func (s *GnfdStorage) EncodedObject(objectType plumbing.ObjectType, hash plumbing.Hash) (plumbing.EncodedObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) IterEncodedObjects(objectType plumbing.ObjectType) (storer.EncodedObjectIter, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) HasEncodedObject(hash plumbing.Hash) error {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) EncodedObjectSize(hash plumbing.Hash) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) SetReference(reference *plumbing.Reference) error {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) CheckAndSetReference(new, old *plumbing.Reference) error {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) Reference(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) IterReferences() (storer.ReferenceIter, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) RemoveReference(name plumbing.ReferenceName) error {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) CountLooseRefs() (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) PackRefs() error {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) SetShallow(hashes []plumbing.Hash) error {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) Shallow() ([]plumbing.Hash, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) SetIndex(index *index.Index) error {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) Index() (*index.Index, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) Config() (*config.Config, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) SetConfig(config *config.Config) error {
	//TODO implement me
	panic("implement me")
}

func (s *GnfdStorage) Module(name string) (storage.Storer, error) {
	//TODO implement me
	panic("implement me")
}
