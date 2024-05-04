package storage

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/bnb-chain/greenfield-go-sdk/types"
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	key = strings.Replace(key, "\n", "", -1)

	size, err := s.head(key)
	if err == nil && size > 0 {
		object, stat, err := s.GnfdClient.GetObject(ctx, s.GetBucketName(), key, types.GetObjectOptions{})
		if err != nil {
			return nil, err
		}
		if stat.Size == 0 {
			return nil, nil
		}
		return io.ReadAll(object)
	}
	return nil, err
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
	if err != nil && (strings.Contains(err.Error(), "the specified object does not exist") || strings.Contains(err.Error(), "No such object")) {
		return false, nil
	}
	return false, err
}

func (s *GnfdStorage) put(key string, value []byte, isOverWrite bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if len(value) != 0 {
		exist, err := s.has(key)
		if err != nil {
			slog.Error("GnfdStorage.put", "error", err)
			return err
		}
		if exist {
			if isOverWrite {
				err = s.GnfdClient.DelegateUpdateObjectContent(ctx, s.GetBucketName(), key, int64(len(value)), bytes.NewReader(value), types.PutObjectOptions{IsUpdate: true})
			}
		} else {
			err = s.GnfdClient.DelegatePutObject(ctx, s.GetBucketName(), key, int64(len(value)), bytes.NewReader(value), types.PutObjectOptions{})
		}
		if err != nil {
			slog.Error("Gnfd put failed", "key", key, "err", err)
			return err
		}
	}
	return nil
}
