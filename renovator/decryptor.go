package renovator

import (
  "fmt"
  "context"
  "encoding/base64"
  "golang.org/x/oauth2/google"
  cloudkms "google.golang.org/api/cloudkms/v1"
)

type Decryptor struct {
  KmsResource   string
  KmsService    *cloudkms.Service
}

func NewDecryptor(projectID string, kmsLocation string, kmsKeyring string, kmsKey string) (*Decryptor, error) {
  d := new(Decryptor)
  d.KmsResource = fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
    projectID, kmsLocation, kmsKeyring, kmsKey)

  ctx := context.Background()
  client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
  if err != nil {
    return nil, err
  }

  cloudkmsService, err := cloudkms.New(client)
  if err != nil {
    return nil, err
  }

  d.KmsService = cloudkmsService
  return d, nil
}

func (e Decryptor) Decrypt(ciphertext []byte) ([]byte, error) {
  req := &cloudkms.DecryptRequest{
    Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
  }

  resp, err := e.KmsService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(e.KmsResource, req).Do()
  if err != nil {
    return nil, err
  }

  decoded, err := base64.StdEncoding.DecodeString(resp.Plaintext)
  if err != nil {
    return nil, err
  }

  return decoded, nil
}
