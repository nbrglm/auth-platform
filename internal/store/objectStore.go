package store

import (
	"context"
	"io"
)

// ObjectStore defines the interface for an object storage service.
//
// The default implementation is S3/MinIO/S3-like, but it can be extended to support other object storage services.
type ObjectStore interface {
	// Init initializes the object store.
	// It returns an error if the initialization fails.
	Init(ctx context.Context) error

	// GetBucketName returns the name of the bucket used by the object store.
	GetBucketName() string

	// UploadObject uploads an object to the specified bucket with the given key prefixed with `private/`.
	// It returns an error if the upload fails.
	//
	// Returns the complete key of the uploaded object, and error if the upload fails.
	//
	// This method does not add any ACL/BucketPolicy.
	UploadObject(ctx context.Context, key string, file io.Reader, contentType, cacheControl string) (string, error)

	// UploadPublicObject uploads an object to the specified bucket with the given key prefixed with 'public/',
	// making it publicly accessible (as the store sets BucketPolicy).
	//
	// Returns the complete key of the uploaded object, and error if the upload fails.
	//
	// Use this method with caution, as it can expose sensitive data if not handled properly.
	UploadPublicObject(ctx context.Context, key string, file io.Reader, contentType, cacheControl string) (string, error)

	// TODO: Implement this method to download an object from the bucket IF NEED BE.
	// DownloadObject downloads an object from the specified bucket with the given key.
	// It returns the data as a byte slice and an error if the download fails.
	// DownloadObject(ctx context.Context, key string) ([]byte, error)

	// DeleteObject deletes an object from the specified bucket with the given key.
	// It returns an error if the deletion fails.
	// NOTE: This method expects the key to be prefixed with 'private/' or 'public/' as per the upload methods.
	DeleteObject(ctx context.Context, key string) error

	// TODO: Implement this method to return a list of object keys in the bucket IF NEED BE.
	// ListObjects lists all objects in the specified bucket.
	// It returns a slice of object keys and an error if the listing fails.
	// ListObjects(bucket string) ([]string, error)

	// GetObjectURL returns the URL of an object in the specified bucket with the given key.
	// It returns the URL as a string and an error if the retrieval fails.
	// Note: This method does not return a pre-signed URL, so it may not be accessible if the object is private.
	GetObjectURL(ctx context.Context, key string) (string, error)

	// GetObjectURLWithExpiry returns a pre-signed URL for an object in the specified bucket with the given key,
	// valid for the specified expiry duration in seconds.
	GetObjectURLWithExpiry(ctx context.Context, key string, expiry int64) (string, error)
}

var Objects ObjectStore

func InitS3Store(ctx context.Context) error {
	// Initialize the S3 store
	Objects = NewS3Store()
	return Objects.Init(ctx)
}
