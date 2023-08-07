package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	"io"
	"math"
	"strings"
	"time"
)

func (s *GnfdStorage) list(bucketName, prefix string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	listResult, err := s.gnfdClient.ListObjects(ctx, bucketName, types.ListObjectsOptions{Prefix: prefix, MaxKeys: math.MaxInt, EndPointOptions: &types.EndPointOptions{}})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, object := range listResult.Objects {
		names = append(names, object.ObjectInfo.ObjectName)
	}
	return names, nil
}

func (s *GnfdStorage) get(bucketName, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	_, err := s.gnfdClient.HeadObject(ctx, bucketName, key)
	if err != nil {
		return nil, err
	}

	object, status, err := s.gnfdClient.GetObject(ctx, bucketName, key, types.GetObjectOptions{})
	_ = status
	if err != nil {
		return nil, err
	}
	val, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (s *GnfdStorage) put(bucketName, key string, value []byte) error {
	fmt.Println("bucketName: ", bucketName, " key: ", key)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	object, err := s.gnfdClient.HeadObject(ctx, bucketName, key)
	if err != nil && !strings.Contains(err.Error(), "No such object") {
		return err
	}

	if err == nil && object != nil {
		return nil
		// return nil, do not overwrite.
		//_, err2 := s.gnfdClient.DeleteObject(ctx, bucketName, key, types.DeleteObjectOption{})
		//if err2 != nil {
		//	return err2
		//}
	}

	txHash, err := s.gnfdClient.CreateObject(
		ctx,
		bucketName,
		key,
		bytes.NewReader(value),
		types.CreateObjectOptions{},
	)
	if err != nil {
		fmt.Println("TxHash: ", txHash)
		return err
	}

	_, err = s.gnfdClient.WaitForTx(ctx, txHash)
	if err != nil {
		fmt.Println("TxHash: ", txHash, "err: ", err)
		return err
	}

	err = s.gnfdClient.PutObject(ctx, bucketName, key, int64(len(value)), bytes.NewReader(value), types.PutObjectOptions{})
	if err != nil {
		fmt.Println("PutObject err : ", err)
		return err
	}
	return nil
}
