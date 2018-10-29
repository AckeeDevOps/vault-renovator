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

// NewClient creates a new Client object with pre-created resty client
func NewClient(vaultAddress string) *Client {
  c := new(Client)

  u, err := url.Parse(vaultAddress)
  if err != nil {
    log.Fatal(err)
  }

  c.VaultAddress = u.Scheme + "://" + u.Host
  c.RestClient = resty.New()
  return c
}

// CheckOrRenew checks the provided token, it tryes to renew the token only
// when the current TLL is beyond the TTL provided threshold
func (c Client) CheckOrRenew(token string, threshold int, increment int) {
  tokenDetails, err := c.lookupSelf(token)
  if err != nil {
    log.Fatal(err)
  }

  if(tokenDetails.TTL <= threshold) {
    c.renew(token)
  }
}

// LookupSelf returns token details https://www.vaultproject.io/api/auth/token/index.html
func (c Client) lookupSelf(token string) (*TokenLookupData, error) {
  resp, err := c.RestClient.R().
    SetHeader("X-Vault-Token", token).
    Get(c.VaultAddress + "/v1/auth/token/lookup-self")
    if err != nil {
      return nil, err
    }
    checkStatusCodeGet(resp.StatusCode(), resp.Body())

    lookupReponse := TokenLookupResponse{}
    json.Unmarshal(resp.Body(), &lookupReponse)
    return &lookupReponse.Data, nil
}

func (c Client) renew(token string) (*TokenLookupData, error) {
  return nil, nil // tbd
}

// DisableTLS disables TLS as deccribed at https://godoc.org/github.com/go-resty/resty#SetTLSClientConfig
func (c Client) DisableTLS() {
  c.RestClient.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
}

// checkStatusCodeGet checks 200 status code for GET requests
func checkStatusCodeGet(code int, body []byte) {
  if(code != 200) {
    log.Fatal("Wrong http status code: " + string(body[:]))
  }
}
