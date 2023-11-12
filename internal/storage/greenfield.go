package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bnb-chain/greenfield-go-sdk/types"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
)

func (s *GnfdStorage) list(prefix, startAfter string, limit uint64) ([]string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	listResult, err := s.GnfdClient.ListObjects(ctx, s.GetBucketName(),
		types.ListObjectsOptions{
			Prefix:     prefix,
			MaxKeys:    limit,
			StartAfter: startAfter})
	if err != nil {
		return nil, "", err
	}

	var names []string
	for _, object := range listResult.Objects {
		names = append(names, object.ObjectInfo.ObjectName)
	}
	return names, listResult.MaxKeys, nil
}

func (s *GnfdStorage) head(key string) (int64, error) {
	object, err := s.GnfdClient.HeadObject(context.Background(), s.GetBucketName(), key)
	if err != nil {
		return 0, err
	}
	return int64(object.ObjectInfo.PayloadSize), nil
}

func (s *GnfdStorage) get(key string) ([]byte, error) {
	//fmt.Println("get key:", key)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	objectDetails, err := s.GnfdClient.HeadObject(ctx, s.GetBucketName(), key)
	if err != nil {
		return nil, err
	}

	if objectDetails.ObjectInfo.PayloadSize == 0 {
		return []byte(""), nil
	}

	object, status, err := s.GnfdClient.GetObject(ctx, s.GetBucketName(), key, types.GetObjectOptions{})
	_ = status
	if err != nil {
		return nil, err
	}
	val, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}
	//fmt.Println("get key: ", key, "value", string(val))
	return val, nil
}

func (s *GnfdStorage) delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	_, err := s.GnfdClient.HeadObject(ctx, s.GetBucketName(), key)
	if err != nil {
		return err
	}
	_, err = s.GnfdClient.DeleteObject(ctx, s.GetBucketName(), key, types.DeleteObjectOption{})
	if err != nil {
		return err
	}
	return nil
}

func (s *GnfdStorage) has(key string) (bool, error) {
	object, err := s.GnfdClient.HeadObject(context.Background(), s.GetBucketName(), key)
	if err == nil && object != nil {
		return true, nil
	}
	return false, err
}

func (s *GnfdStorage) put(key string, value []byte, isOverWrite bool) error {
	//fmt.Println("RepoName: ", s.GetBucketName(), " key: ", key, "isOverwrite: ", isOverWrite)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	object, err := s.GnfdClient.HeadObject(ctx, s.GetBucketName(), key)
	if err != nil && !strings.Contains(err.Error(), "No such object") {
		return err
	}

	if err == nil && object.ObjectInfo != nil {
		if isOverWrite {
			if object.ObjectInfo.ObjectStatus == storagetypes.OBJECT_STATUS_SEALED {
				_, err2 := s.GnfdClient.DeleteObject(ctx, s.GetBucketName(), key, types.DeleteObjectOption{})
				if err2 != nil {
					return err2
				}
			} else {
				_, err2 := s.GnfdClient.CancelCreateObject(ctx, s.GetBucketName(), key, types.CancelCreateOption{})
				if err2 != nil {
					return err2
				}
			}
			time.Sleep(3 * time.Second)
		} else {
			return nil
		}
	}

	retry_create := true
	for i := 0; i < 3; i++ {
		if retry_create {
			_, err = s.GnfdClient.CreateObject(
				ctx,
				s.GetBucketName(),
				key,
				bytes.NewReader(value),
				types.CreateObjectOptions{},
			)

			if err != nil {
				fmt.Println("Create Object failed, err: ", err)
				return err
			}
		}
		if len(value) != 0 {
			err = s.GnfdClient.PutObject(ctx, s.GetBucketName(), key, int64(len(value)), bytes.NewReader(value), types.PutObjectOptions{})
			if err != nil {
				if strings.Contains(err.Error(), "invalid payload data integrity hash") {
					_, err2 := s.GnfdClient.CancelCreateObject(ctx, s.GetBucketName(), key, types.CancelCreateOption{})
					if err2 != nil {
						return err2
					}
					time.Sleep(3 * time.Second)
				} else {
					retry_create = false
					fmt.Println("Put Object to greenfield failed, err: ", err)
				}
			} else {
				break
			}
		}
	}

	return err
}
