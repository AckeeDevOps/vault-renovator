package renovator

import (
  "log"
  "context"
  "io/ioutil"
  "cloud.google.com/go/storage"
)

type TokenProvider struct {
  BucketName    string
  FilePaths     []string
  StorageClient *storage.Client
}

// NewTokenProvider constructor
func NewTokenProvider (bucketName string, filePaths []string) *TokenProvider {
  p := new(TokenProvider)
  p.BucketName = bucketName
  p.FilePaths = filePaths

  // initialize google storage client
  ctx := context.Background()
  client, err := storage.NewClient(ctx)
  if err != nil {
    log.Fatal(err)
  }

  p.StorageClient = client
  return p
}

func (e TokenProvider) GetTokens() ([][]byte, error) {
  res, err := e.getEncryptedBytes()
  if err != nil {
    return nil, err
  }
  return res, nil
}

func (e TokenProvider) getEncryptedBytes() ([][]byte, error) {
  results := [][]byte{}
  // go through all provided file paths
  for _, v := range e.FilePaths {
    ctx := context.Background()

    // obtain files with GCS client
    rc, err := e.StorageClient.Bucket(e.BucketName).Object(v).NewReader(ctx)
    if err != nil {
      return nil, err
    }
    defer rc.Close()

    data, err := ioutil.ReadAll(rc)
    if err != nil {
      return nil, err
    }

    // append data to the slice
    results = append(results, data)
  }
  return results, nil
}
