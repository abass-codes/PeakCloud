package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectInfo struct {
	ETag string
	Size int64
}

type ObjectStore struct {
	client *minio.Client
	bucket string
}

func NewObjectStore(
	endpoint string,
	accessKey string,
	secretKey string,
	bucket string,
	useSSL bool,
) (*ObjectStore, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	return &ObjectStore{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *ObjectStore) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("check bucket: %w", err)
	}

	if exists {
		return nil
	}

	if err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("create bucket: %w", err)
	}

	return nil
}

func (s *ObjectStore) Put(
	ctx context.Context,
	objectKey string,
	reader io.Reader,
	size int64,
	contentType string,
) (*ObjectInfo, error) {
	info, err := s.client.PutObject(
		ctx,
		s.bucket,
		objectKey,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("put object: %w", err)
	}

	return &ObjectInfo{
		ETag: info.ETag,
		Size: info.Size,
	}, nil
}

func (s *ObjectStore) Get(
	ctx context.Context,
	objectKey string,
) (*minio.Object, error) {
	object, err := s.client.GetObject(
		ctx,
		s.bucket,
		objectKey,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}

	if _, err := object.Stat(); err != nil {
		_ = object.Close()
		return nil, fmt.Errorf("stat object: %w", err)
	}

	return object, nil
}

func (s *ObjectStore) Delete(ctx context.Context, objectKey string) error {
	if err := s.client.RemoveObject(
		ctx,
		s.bucket,
		objectKey,
		minio.RemoveObjectOptions{},
	); err != nil {
		return fmt.Errorf("remove object: %w", err)
	}

	return nil
}

func (s *ObjectStore) Ping(ctx context.Context) error {
	_, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("object storage health check: %w", err)
	}

	return nil
}
