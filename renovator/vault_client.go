package renovator

import (
  "log"
  "net/url"
  //"encoding/json"
  "gopkg.in/resty.v1"
)

type Client struct {
  VaultAddress  string
  VaultToken    string
}

type TokenLookupData struct {
  Accessor      string `json:"accessor"`
  CreationTime  int `json:"creation_time"`
  CreationTtl   int `json:"creation_ttl"`
  DisplayName   string `json:"display_name"`
  ExpireTime    string `json:"expire_time"`
  IssueTime     string `json:"issue_time"`
  Ttl           int `json:"ttl"`
}

type TokenLookupResponse struct {
  Data TokenLookupData `json:"data"`
}


func NewClient(vaultAddress string, vaultToken string) *Client {
  c := new(Client)

  u, err := url.Parse(vaultAddress)
  if err != nil {
    log.Fatal(err)
  }

  c.VaultAddress = u.Scheme + u.Host
  c.VaultToken = vaultToken
  return c
}

func (c Client) LookupSelf() (string, error) {
  _, err := resty.R().
    SetHeader("X-Vault-Token", c.VaultToken).
    Get(c.VaultAddress + "/auth/token/lookup-self")
    if err != nil {
      return "", err
    }
    return "", nil // delete me!
}
