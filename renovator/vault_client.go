package renovator

import (
  "log"
  "net/url"
  "crypto/tls"
  "encoding/json"
  "gopkg.in/resty.v1"
)

type Client struct {
  VaultAddress  string
  VaultToken    string
  RestClient    *resty.Client
}

type TokenLookupData struct {
  Accessor      string `json:"accessor"`
  CreationTime  int `json:"creation_time"`
  CreationTTL   int `json:"creation_ttl"`
  DisplayName   string `json:"display_name"`
  ExpireTime    string `json:"expire_time"`
  IssueTime     string `json:"issue_time"`
  Renewable     bool `json:"renewable"`
  TTL           int `json:"ttl"`
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
  c.RestClient = resty.New()
  return c
}

// LookupSelf returns token details https://www.vaultproject.io/api/auth/token/index.html
func (c Client) LookupSelf() (*TokenLookupData, error) {
  resp, err := c.RestClient.R().
    SetHeader("X-Vault-Token", c.VaultToken).
    Get(c.VaultAddress + "/auth/token/lookup-self")
    if err != nil {
      return nil, err
    }
    checkStatusCode(resp.StatusCode(), resp.Body())

    lookupReponse := TokenLookupResponse{}
    json.Unmarshal(resp.Body(), &lookupReponse)
    return &lookupReponse.Data, nil
}

func (c Client) DisableTLS() {
  c.RestClient.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
}

func checkStatusCode(code int, body []byte) {
  if(code != 200) {
    log.Fatal("Wrong http status code: " + string(body[:]))
  }
}
