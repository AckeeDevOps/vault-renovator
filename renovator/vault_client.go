package renovator

import (
  "log"
  "net/url"
  "crypto/tls"
  "encoding/json"
  "gopkg.in/resty.v1"
)

const msgRenewalNotNeeded =   "%s... does not need renewal. It will expire in ~%d days. " +
                              "It's %d days above the threshold."
const msgStartRenewal =       "%s... will expire in %d days. It's going to be renewed. " +
                              "Increment will be %d seconds."
const msgRenewalSuccessful =  "%s... has been renewed. New TTL is %d (it was %d before)."
const msgRenewalFailed =      "%s... has not been renewed. New TTL (%d) is not greater than " +
                              "the old one (%d)."

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

type RenewalRequest struct {
  Increment int `json:"increment"`
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
    log.Printf(msgStartRenewal, token[0:7], (tokenDetails.TTL / 60 / 60 / 24) , increment)

    // it seems that increment is not increment but total TTL
    // hence we're appending increment to the current TTL
    // tested with Vault v0.11.4
    tokenDetailsNew, err := c.renew(token, tokenDetails.TTL + increment)
    if err != nil {
      log.Fatal(err)
    }
    compareTTL(token, tokenDetails.TTL, tokenDetailsNew.TTL)

  } else {
    days := tokenDetails.TTL / 60 / 60 / 24
    aboveThreshold := (tokenDetails.TTL - threshold) / 60 / 60 / 24
    log.Printf(msgRenewalNotNeeded, token[0:7], days, aboveThreshold)
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
    checkStatusCode(resp.StatusCode(), resp.Body())

    lookupReponse := TokenLookupResponse{}
    json.Unmarshal(resp.Body(), &lookupReponse)
    return &lookupReponse.Data, nil
}

func (c Client) renew(token string, increment int) (*TokenLookupData, error) {
  resp, err := c.RestClient.R().
  SetHeader("X-Vault-Token", token).
  SetBody(RenewalRequest{Increment: increment}).
  Post(c.VaultAddress + "/v1/auth/token/renew-self")
  if err != nil {
    return nil, err
  }
  checkStatusCode(resp.StatusCode(), resp.Body())

  // get token details again
  tokenDetails, err := c.lookupSelf(token)
  if err != nil {
    log.Fatal(err)
  }
  return tokenDetails, nil
}

// DisableTLS disables TLS as deccribed at https://godoc.org/github.com/go-resty/resty#SetTLSClientConfig
func (c Client) DisableTLS() {
  c.RestClient.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
}

// checkStatusCodeGet checks 200 status code for GET and POST requests
// code 200 is valid for both endpoints used in this context
func checkStatusCode(code int, body []byte) {
  if(code != 200) {
    log.Fatal("Wrong http status code: " + string(body[:]))
  }
}

func compareTTL(token string, old int, new int) bool {
  msg := msgRenewalSuccessful
  if(new <= old) {
    msg = msgRenewalFailed
  }
  log.Printf(msg, token[0:7], new, old)
  return (new > old)
}
