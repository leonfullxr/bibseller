// Package storage wraps the S3-compatible object store (MinIO in dev, R2 or
// Scaleway in prod - docs/ARCHITECTURE.md) used for private in-chat images. The
// bucket is private: objects are reached only through the API's authorized
// download route, never a public URL.
package storage

import (
	"context"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	mc     *minio.Client
	bucket string
}

// New builds a client. The endpoint may carry an http:// or https:// scheme
// (http -> insecure, matching the dev MinIO); a bare host[:port] defaults to
// secure for prod.
func New(endpoint, accessKey, secretKey, bucket string) (*Client, error) {
	endpoint = strings.TrimSpace(endpoint)
	secure, host := true, endpoint
	switch {
	case strings.HasPrefix(endpoint, "http://"):
		secure, host = false, strings.TrimPrefix(endpoint, "http://")
	case strings.HasPrefix(endpoint, "https://"):
		secure, host = true, strings.TrimPrefix(endpoint, "https://")
	}
	host = strings.TrimSuffix(host, "/") // a trailing slash makes an invalid minio host
	mc, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}
	return &Client{mc: mc, bucket: bucket}, nil
}

// EnsureBucket creates the (private) bucket if it does not exist. Called on boot;
// also the signal that the store is reachable.
func (c *Client) EnsureBucket(ctx context.Context) error {
	exists, err := c.mc.BucketExists(ctx, c.bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return c.mc.MakeBucket(ctx, c.bucket, minio.MakeBucketOptions{})
}

// Put stores an object. The bucket carries no public policy, so the object is
// private regardless of key.
func (c *Client) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	_, err := c.mc.PutObject(ctx, c.bucket, key, r, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}

// Get opens an object for reading, erroring if it does not exist (Stat surfaces
// a missing object now rather than mid-stream).
func (c *Client) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := c.mc.GetObject(ctx, c.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	if _, err := obj.Stat(); err != nil {
		_ = obj.Close()
		return nil, err
	}
	return obj, nil
}

// Delete removes an object. Used to clean up an image whose message row failed
// to commit, so no orphan is left in the bucket.
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.mc.RemoveObject(ctx, c.bucket, key, minio.RemoveObjectOptions{})
}

// DeleteMany removes multiple objects in one batched call (RemoveObjects)
// instead of one round trip per key - the retention job's per-batch cleanup
// (#99). Reports the first per-object error, if any; the rest are attempted
// regardless (a partial failure here means an orphaned object, not a retry-
// blocking one - the caller logs it, per purgeExpiredMessages).
func (c *Client) DeleteMany(ctx context.Context, keys []string) error {
	objectsCh := make(chan minio.ObjectInfo, len(keys))
	for _, key := range keys {
		objectsCh <- minio.ObjectInfo{Key: key}
	}
	close(objectsCh)

	var firstErr error
	for result := range c.mc.RemoveObjects(ctx, c.bucket, objectsCh, minio.RemoveObjectsOptions{}) {
		if result.Err != nil && firstErr == nil {
			firstErr = result.Err
		}
	}
	return firstErr
}

// IsNotFound reports whether err is a missing-object error from this store. The
// result comes from the error argument, not the receiver; the method form keeps
// callers decoupled from the minio package.
func (c *Client) IsNotFound(err error) bool {
	return minio.ToErrorResponse(err).Code == "NoSuchKey"
}
