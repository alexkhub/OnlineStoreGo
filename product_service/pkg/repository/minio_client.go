package repository

import (
	"context"
	"fmt"
	"log"

	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
)

type MinioClient struct {
	minIO *minio.Client
}

func NewMinioClient(minIO *minio.Client) *MinioClient {
	return &MinioClient{minIO: minIO}
}

func (m *MinioClient) GetOne(bucketName string, objectID string) (string, error) {

	url, err := m.minIO.PresignedGetObject(context.Background(), bucketName, objectID, time.Second*24*60*60, nil)
	if err != nil {
		return "", fmt.Errorf("get image URL error %s: %v", objectID, err)
	}
	return strings.Replace(url.String(), "minio", "localhost", 1), nil
}

func (m *MinioClient) GetMany(backetName string, objectIDs []string) (map[string]string, error) {

	urlCh := make(chan map[string]string, len(objectIDs))
	errCh := make(chan error, len(objectIDs))

	var wg sync.WaitGroup
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, objectID := range objectIDs {
		wg.Add(1)
		go func(objectID string) {
			defer wg.Done()

			url, err := m.GetOne(backetName, objectID)

			if err != nil {
				errCh <- err
				cancel() 
				return
			}
			urlCh <- map[string]string{
				objectID: url,
			}
		}(objectID)
	}

	go func() {
		wg.Wait()
		close(urlCh)
		close(errCh)
	}()

	var urls = make(map[string]string)
	var errs []error
	for url := range urlCh {

		for k, v := range url {
			urls[k] = v
		}
	}
	for opErr := range errCh {
		errs = append(errs, opErr)
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("errors when receiving objects: %v", errs)
	}

	return urls, nil
}

func (m *MinioClient) PresignedListObject(bucketName, prefix string, recursive bool) ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	defer cancel()

	objectCh := m.minIO.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	})
	var images []string

	for object := range objectCh {

		if object.Err != nil {
			return nil, object.Err

		}
		images = append(images, object.Key)
	}
	urlCh := make(chan string, len(images))
	errCh := make(chan error, len(images))
	for _, key := range images {
		wg.Add(1)
		go func(imageKey string) {
			defer wg.Done()

			url, err := m.GetOne(bucketName, imageKey)
			if err != nil {
				errCh <- err
				return
			}
			urlCh <- url

		}(key)
	}

	go func() {
		wg.Wait()
		close(urlCh)
		close(errCh)
	}()

	var urls []string
	var errs []error

	for opErr := range errCh {
		errs = append(errs, opErr)
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("errors when receiving objects: %v", errs)
	}

	for url := range urlCh {
		urls = append(urls, url)
	}

	return urls, nil

}

func (m *MinioClient) RemoveAllObjects(bucketName, prefix string, recursive bool) {
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for object := range m.minIO.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: recursive,
		}) {
			if object.Err != nil {
				log.Fatalln(object.Err)
			}
			objectsCh <- object
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}
	for rErr := range m.minIO.RemoveObjects(context.Background(), bucketName, objectsCh, opts) {
		log.Println("Error detected during deletion: ", rErr)
	}
}

func (m *MinioClient) RemoveOne(bucketName, objectID string) error {

	err := m.minIO.RemoveObject(context.Background(), bucketName, objectID, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
