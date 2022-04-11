package storage

import (
	"bytes"
	"fmt"
	"io"
	"maintainman/config"
	"maintainman/logger"
	"os"
	"path/filepath"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/spf13/viper"
)

var (
	Storage           IStorage
	ImageCacheStorage IStorage
	s3Conn            *s3.S3
)

type IStorage interface {
	Exist(id string) bool
	Load(id string, fn func(io.Reader) error) error
	Save(id, format string, fn func(io.Writer) error) error
	LoadBytes(id string) ([]byte, error)
	SaveBytes(id, format string, data []byte) error
	Delete(id string) error
}

func init() {
	Storage = initStorage(config.AppConfig)
	ImageCacheStorage = initStorage(config.ImageConfig)
}

func initStorage(config *viper.Viper) (storage IStorage) {
	storageType := config.GetString("storage.driver")
	switch storageType {
	case "local":
		path := config.GetString("storage.local.path")
		clean := config.GetBool("storage.local.clean")
		storage = newLocalStorage(path, clean)
	case "s3":
		if s3Conn == nil {
			initS3Conn(config)
		}
		bucket := config.GetString("storage.s3.bucket")
		storage = newS3Storage(bucket)
	default:
		panic(fmt.Errorf("support local and s3 only"))
	}
	return storage
}

func initS3Conn(config *viper.Viper) {
	accessKey := config.GetString("storage.s3.access_key")
	secretKey := config.GetString("storage.s3.secret_key")
	auth, err := aws.GetAuth(accessKey, secretKey)
	if err != nil {
		panic(fmt.Errorf("aws cannot get auth: %v", err))
	}

	bucket := config.GetString("storage.s3.bucket")
	if bucket == "" {
		panic(fmt.Errorf("s3 bucket not set"))
	}

	region, ok := aws.Regions[config.GetString("storage.s3.region")]
	if !ok {
		region = aws.EUWest
	}

	s3Conn = s3.New(auth, region)
}

func newLocalStorage(path string, clean bool) *LocalStorage {
	if clean {
		os.RemoveAll(path)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	return &LocalStorage{
		path: path,
	}
}

func newS3Storage(bucket string) *S3Storage {
	bucketObj := s3Conn.Bucket(bucket)
	return &S3Storage{
		bucket: bucketObj,
	}
}

type LocalStorage struct {
	path string
}

func (s *LocalStorage) Exist(id string) bool {
	fullPath := filepath.Join(s.path, id)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (s *LocalStorage) Load(id string, fn func(io.Reader) error) (err error) {
	fullPath := filepath.Join(s.path, id)
	reader, err := os.Open(fullPath)
	defer reader.Close()
	if err != nil {
		return fmt.Errorf("failed to load %s: %v", id, err)
	}
	err = fn(reader)
	return err
}

func (s *LocalStorage) LoadBytes(id string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if err := s.Load(id, func(reader io.Reader) error {
		_, err := io.Copy(buffer, reader)
		return err
	}); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *LocalStorage) Save(id, format string, fn func(io.Writer) error) error {
	// Open file for writing, overwrite if it already exists
	fullPath := filepath.Join(s.path, id)
	writer, err := os.Create(fullPath)
	defer writer.Close()
	if err != nil {
		return err
	}
	if err := fn(writer); err != nil {
		return err
	}
	return nil
}

func (s *LocalStorage) SaveBytes(id, format string, data []byte) error {
	return s.Save(id, format, func(writer io.Writer) error {
		_, err := writer.Write(data)
		return err
	})
}

func (s *LocalStorage) Delete(id string) error {
	fullPath := filepath.Join(s.path, id)
	return os.Remove(fullPath)
}

// S3Storage is a storage implementation using Amazon S3
type S3Storage struct {
	bucket *s3.Bucket
}

func (s *S3Storage) Exist(id string) bool {
	resp, err := s.bucket.List(id, "/", "", 10)
	if err != nil {
		logger.Logger.Errorf("Error while listing S3 bucket: %v\n", err)
		return false
	}
	if resp == nil {
		logger.Logger.Error("Error while listing S3 bucket: empty response")
	}

	for _, element := range resp.Contents {
		if element.Key == id {
			return true
		}
	}

	return false
}

func (s *S3Storage) Load(id string, fn func(io.Reader) error) (err error) {
	rc, err := s.bucket.GetReader(id)
	defer rc.Close()
	if err != nil {
		return err
	}
	err = fn(rc)
	return err
}

func (s *S3Storage) LoadBytes(id string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if err := s.Load(id, func(reader io.Reader) error {
		_, err := io.Copy(buffer, reader)
		return err
	}); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *S3Storage) Save(id, format string, fn func(io.Writer) error) error {
	buffer := bytes.NewBuffer(nil)
	if err := fn(buffer); err != nil {
		return err
	}
	return s.SaveBytes(id, format, buffer.Bytes())
}

func (s *S3Storage) SaveBytes(id, format string, data []byte) error {
	return s.bucket.Put(id, data, format, s3.Private)
}

func (s *S3Storage) Delete(id string) error {
	return s.bucket.Del(id)
}