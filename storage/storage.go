package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage"
)

type GnfdStorage struct {
	GnfdClient client.IClient
	RepoName   string
	Account    *types.Account
}

func (s *GnfdStorage) AddAlternate(remote string) error {
	//TODO implement me
	panic("implement me")
}

func NewStorage(chainID, rpcAddress, privateKey, bucketName string) (*GnfdStorage, error) {
	//fmt.Println("ChainID: ", chainID, " rpcAddress: ", rpcAddress, " privateKey: ", privateKey, " RepoName: ", RepoName)
	account, err := types.NewAccountFromPrivateKey("gitd", privateKey)
	if err != nil {
		fmt.Println("New account from private key error: ", err)
		return nil, err
	}
	gnfdClient, err := client.New(chainID, rpcAddress, client.Option{DefaultAccount: account})
	if err != nil {
		fmt.Println("New greenfield storage client error: ", err)
		return nil, err
	}

	_, err = gnfdClient.GetLatestBlock(context.Background())
	if err != nil {
		fmt.Println("Get latest block error: ", err)
		return nil, err
	}

	return &GnfdStorage{
		GnfdClient: gnfdClient,
		RepoName:   bucketName,
		Account:    account,
	}, nil
}

func (s *GnfdStorage) GetBucketName() string {
	return s.RepoName
}

func (s *GnfdStorage) NewEncodedObject() plumbing.EncodedObject {
	return &plumbing.MemoryObject{}
}
func (s *GnfdStorage) SetEncodedObject(obj plumbing.EncodedObject) (plumbing.Hash, error) {
	r, err := obj.Reader()
	if err != nil {
		return obj.Hash(), err
	}

	c, err := io.ReadAll(r)
	if err != nil {
		return obj.Hash(), err
	}

	// {type} {length}\x00
	var buffer bytes.Buffer
	buffer.Write(obj.Type().Bytes())
	buffer.Write([]byte(" "))
	buffer.WriteString(strconv.Itoa(len(c)))
	buffer.Write([]byte{0x00})
	buffer.Write(c)
	err = s.put(buildObjectsKey(obj.Hash()), buffer.Bytes(), false)
	return obj.Hash(), err
}

func (s *GnfdStorage) EncodedObject(t plumbing.ObjectType, h plumbing.Hash) (plumbing.EncodedObject, error) {
	var err error

	// get object key
	rec, err := s.get(buildObjectsKey(h))
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, plumbing.ErrObjectNotFound
	}

	return parseObjectFromBlob(t, rec)
}

func parseObjectFromBlob(t plumbing.ObjectType, blob []byte) (plumbing.EncodedObject, error) {
	s := bytes.IndexByte(blob, 32) // first space
	i := bytes.IndexByte(blob, 0)  // first null value

	if s == -1 || i == -1 {
		return nil, fmt.Errorf("invalid buffer format")
	}

	// get type of object
	typeVal := string(blob[:s])
	actualType, err := plumbing.ParseObjectType(typeVal)
	if err != nil {
		return nil, fmt.Errorf("parse object type failed")
	}

	if t != plumbing.AnyObject && actualType != t {
		return nil, fmt.Errorf("the object type mismatch, %v - %v", t, actualType)
	}

	// get length of object
	lengthBytes := blob[s+1 : i]
	objectLength, err := strconv.Atoi(string(lengthBytes))
	if err != nil {
		return nil, err
	}

	actualLength := len(blob) - (i + 1)

	// verify length
	if objectLength != actualLength {
		return nil, fmt.Errorf("length mismatch: expected %d bytes but got %d instead", objectLength, actualLength)
	}
	o := &plumbing.MemoryObject{}
	o.SetType(actualType)
	o.SetSize(int64(objectLength))

	_, err = o.Write(blob[i+1:])
	if err != nil {
		return nil, err
	}
	return o, nil
}

type EncodedObjectIter struct {
	t             plumbing.ObjectType
	s             *GnfdStorage
	nextKey       string
	encodeObjects []plumbing.EncodedObject
	limitSizeOnce uint64
}

func NewEncodeObjectIter(t plumbing.ObjectType, s *GnfdStorage) *EncodedObjectIter {
	return &EncodedObjectIter{
		t:             t,
		s:             s,
		nextKey:       "",
		limitSizeOnce: 100,
	}
}

func (i *EncodedObjectIter) Next() (plumbing.EncodedObject, error) {
	if len(i.encodeObjects) == 0 {
		objects, maxKey, err := i.s.list(ObjectKey, i.nextKey, i.limitSizeOnce)
		if err != nil {
			return nil, err
		}
		for _, object := range objects {
			hash, found := strings.CutPrefix(object, ObjectKey)
			if !found {
				panic("Iter encode objects error, prefix not found")
			}

			encodeObject, err := i.s.EncodedObject(i.t, plumbing.NewHash(hash))
			if err != nil {
				return nil, err
			}
			i.encodeObjects = append(i.encodeObjects, encodeObject)
		}
		i.nextKey = maxKey
	}
	encodeObject := i.encodeObjects[0]
	i.encodeObjects = i.encodeObjects[1:]
	return encodeObject, nil
}

func (i *EncodedObjectIter) ForEach(cb func(obj plumbing.EncodedObject) error) error {
	for {
		obj, err := i.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		if err := cb(obj); err != nil {
			if err == storer.ErrStop {
				return nil
			}
			return err
		}
	}
}

func (i *EncodedObjectIter) Close() {}

func (s *GnfdStorage) IterEncodedObjects(objectType plumbing.ObjectType) (storer.EncodedObjectIter, error) {
	return NewEncodeObjectIter(objectType, s), nil
}

func (s *GnfdStorage) HasEncodedObject(hash plumbing.Hash) error {
	found, err := s.has(buildObjectsKey(hash))
	if err != nil {
		return err
	}
	if found {
		return nil
	} else {
		return plumbing.ErrObjectNotFound
	}
}

func (s *GnfdStorage) EncodedObjectSize(hash plumbing.Hash) (int64, error) {
	return s.head(buildObjectsKey(hash))

}

func (s *GnfdStorage) SetReference(reference *plumbing.Reference) error {
	var val []byte
	switch reference.Type() {
	case plumbing.HashReference:
		val = []byte(reference.Hash().String())
	case plumbing.SymbolicReference:
		val = []byte(fmt.Sprintf("ref: %s\n", reference.Target()))
	}
	return s.put(buildReferenceKey(reference.Name()), val, true)
}

func (s *GnfdStorage) CheckAndSetReference(new, old *plumbing.Reference) error {
	panic("implement me")
}

func (s *GnfdStorage) Reference(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	target, err := s.get(buildReferenceKey(name))
	if err != nil {
		if strings.Contains(err.Error(), "No such object") || strings.Contains(err.Error(), "the specified object does not exist") {
			return nil, plumbing.ErrReferenceNotFound
		} else {
			return nil, err
		}
	}

	return plumbing.NewReferenceFromStrings(name.String(), string(target[:])), nil
}

func (s *GnfdStorage) IterReferences() (storer.ReferenceIter, error) {
	refKeys, _, err := s.list(ReferenceKey, "", math.MaxUint64)
	if err != nil {
		fmt.Println("list failed, error: ", err, ", RepoName: ", s.GetBucketName())
		return nil, err
	}

	var refs []*plumbing.Reference
	for _, refName := range refKeys {
		refName, _ = strings.CutPrefix(refName, ReferenceKey)
		ref, err := s.Reference(plumbing.ReferenceName(refName))
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return storer.NewReferenceSliceIter(refs), nil
}

func (s *GnfdStorage) RemoveReference(name plumbing.ReferenceName) error {
	return s.delete(buildReferenceKey(name))
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
	rec, err := s.get(ConfigKey)
	if err != nil {
		return nil, err
	}

	if rec == nil {
		return config.NewConfig(), nil
	}

	c := &config.Config{}
	return c, json.Unmarshal(rec, c)
}

func (s *GnfdStorage) SetConfig(c *config.Config) error {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return s.put(ConfigKey, jsonBytes, true)

}

func (s *GnfdStorage) Module(name string) (storage.Storer, error) {
	//TODO implement me
	panic("implement me")
}
