package renovator

type ProgramOptions struct {
  // GCP stuff
  ProjectID       string  `long:"gcp-project" description:"" required:"yes"`
  Location        string  `long:"gcp-location" description:"" required:"yes"`
  Keyring         string  `long:"keyring" description:"" required:"yes"`
  Key             string  `long:"key" description:"" required:"yes"`
  Bucket          string  `long:"bucket" description:"" required:"yes"`

  // Vault stuff
  VaultAddr       string  `long:"vault-url" description:"" required:"yes"`
  ThresholdTTL    int     `long:"ttl-threshold" description:"" required:"no" default:"432000"`
  IncrementTTL    int     `long:"ttl-increment" description:"" required:"no" default:"86400"`

  // JSON file with remote filenames
  SpecsPath       string  `long:"token-specs" description:"" required:"yes"`

  // Slack stuff
  SlackWebhookUrl string  `long:"slack-webhook-url" description:"" required:"no"`
}
