package main

import(
  "os"
  "log"
  "strings"
  "io/ioutil"
  "encoding/json"
  "github.com/nlopes/slack"
  "github.com/jessevdk/go-flags"
  "github.com/vranystepan/vault-renovator/renovator"
)

type ProgamOptions struct {
  // GCP stuff
  ProjectID string `long:"gcp-project" description:"" required:"yes"`
  Location string `long:"gcp-location" description:"" required:"yes"`
  Keyring string `long:"keyring" description:"" required:"yes"`
  Key string `long:"key" description:"" required:"yes"`
  Bucket string `long:"bucket" description:"" required:"yes"`
  // Vault stuff
  VaultAddr string `long:"vault-url" description:"" required:"yes"`
  ThresholdTTL int `long:"ttl-threshold" description:"" required:"no" default:"432000"` // 5 days
  IncrementTTL int `long:"ttl-increment" description:"" required:"no" default:"86400"` // 1 day
  // JSON file with remote filenames
  SpecsPath string `long:"token-specs" description:"" required:"yes"`
  // Slack stuff
  SlackWebhookUrl string `long:"slack-webhook-url" description:"" required:"no"`
}

type TokenFileNames struct {
  Names []string `json:"fileNames"`
}

func main() {
  //disable timestamp in the log output
  log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

  log.Print("Running vault-renovator ...")
  args := handleInput(os.Args[1:])
  decryptor := getDecryptor(args)
  tokens := getTokens(args)

  tokensPlainText, err := decryptTokens(tokens, decryptor) // start here
  if err != nil {
    log.Fatal(err)
  }

  client := renovator.NewClient(args.VaultAddr)
  statusList := []renovator.OutputRenewalStatus{}
  for _, v := range tokensPlainText {
    status := client.CheckOrRenew(v, args.ThresholdTTL, args.IncrementTTL)
    statusList = append(statusList, status)
  }

  // notify slack if SlackWebhookUrl has been set
  if(args.SlackWebhookUrl != ""){
    notifySlackFinal(statusList, args.SlackWebhookUrl)
  }
}

func notifySlackFinal(statusList []renovator.OutputRenewalStatus, url string) {
  attachments := statusListToAttachments(statusList)
  msg := slack.WebhookMessage{Text: "Vault token renewal status", Attachments: attachments}
  err := slack.PostWebhook(url, &msg)
  if err != nil {
    log.Fatal(err)
  }
}

func handleInput(args []string) *ProgamOptions {
  opts := ProgamOptions{}
  parser := flags.NewParser(&opts, (flags.HelpFlag | flags.PassDoubleDash))
  _, err := parser.ParseArgs(args);
  if err != nil {
    log.Fatal(err)
  }
  return &opts
}

func handleSpecsFile(path string) []string {
  jsonFile, err := os.Open(path)
  if err != nil {
    log.Fatal(err)
  }
  defer jsonFile.Close()

  byteValue, err := ioutil.ReadAll(jsonFile)
  if err != nil {
    log.Fatal(err)
  }

  fileNames := TokenFileNames{}
  json.Unmarshal(byteValue, &fileNames) // returns empty TokenFileNames on error
  return fileNames.Names
}

func getTokens(args *ProgamOptions) [][]byte {
  fileNames := handleSpecsFile(args.SpecsPath)
  provider := renovator.NewTokenProvider(args.Bucket, fileNames)
  tokens, err := provider.GetTokens()
  if err != nil {
    log.Fatal(err) //stop here
  }
  return tokens
}

func getDecryptor(args *ProgamOptions) *renovator.Decryptor {
  decryptor, err := renovator.NewDecryptor(args.ProjectID, args.Location, args.Keyring, args.Key)
  if err != nil {
    log.Fatal(err) //stop here
  }
  return decryptor
}

func decryptTokens(tokens [][]byte, decryptor *renovator.Decryptor) ([]string, error) {
  results := []string{}
  for _, v := range tokens {
    decodedToken, err := decryptor.Decrypt(v)
    if err != nil {
      return nil, err
    }
    results = append(results, strings.TrimSpace(string(decodedToken[:])))
  }
  return results, nil
}

func statusListToAttachments(list []renovator.OutputRenewalStatus) []slack.Attachment{
  attachments := []slack.Attachment{}
  for _, v := range list {
    color := ""
    switch v.StatusMessage {
      case "RENEWAL_DONE": color = "#008000"
      case "RENEWAL_NOT_NEEDED": color = "#008000"
      case "RENEWAL_FAILED": color = "#FF0000"
      default: color = "#808080"
    }
    msgBody, err := json.Marshal(v.TokenDetails)
    if err != nil {
      return nil
    }
    attachment := slack.Attachment{
      Color: color,
      Title: v.TokenDetails.Accessor,
      Text: string(msgBody[:]),
    }
    attachments = append(attachments, attachment)
  }
  return attachments
}
