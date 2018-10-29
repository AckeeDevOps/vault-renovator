# Work In Progress 👨‍💻

This repository will become a Docker image for the automatic Vault token
renewal. I work hard to finish it ASAP as I really need this in my production infrastructure 😂

## Purpose

This code is meant to be running inside the Kubernetes (GKE) pod as a CronJob. This solution is exclusively using Google Cloud Platform, porting to different platforms would most likely require a great slaughter.

This tools is processing Vault tokens (encrypted with KMS and saved in the GCS bucket) and renews them if needed based on the provided specification. Please note that vault-renovator uses these tokens even for the login and renewal (so there's no need of additional 'admin' token or account) so tokens should have sufficient privileges: `lookup-self` and `renew-self`. Both capabilities are covered by the `default` policy so most likely you don't have to pay extra attention here.

```hcl
path "auth/token/lookup-self" {
  capabilities = ["read"]
}

path "auth/token/renew-self" {
  capabilities = ["update"]
}
```

Please note that vault-renovator is not trying to separate you from the all the hassle with the  Vault. This tool does only one thing - it renews the specified Vault tokens based on the schedule specified in the CronJob manifest file.

## Command line parameters

|param|Description|
| ------------------|---------------------------|
|**gcp-project**|ID of your GCP project, please note that ID might not be the same as name|
|**gcp-location**|[Location](https://cloud.google.com/compute/docs/regions-zones/) of your KMS keyring|
|**keyring**|Name of the KMS keyring|
|**key**|Name of the KMS key (subordinate of keyring)|
|**bucket**|Name of the GCS bucket|
|**vault-url**|Base url of your Vault instance|
|**token-specs**|File path of the JSON file with the list of token files|
|**slack-webhook-url**|*(optional)* Slack Webhook url for the notifications|
|**ttl-threshold**|*(optional)* Threshold in seconds, the default value is equivalent of 5 days|
|**ttl-increment**|*(optional)* TTL increment is seconds, the default value is equivalent of 1 day|

## Sample usage from the command line

```
go run main.go \
  --gcp-project=test-cloud-1234 \
  --gcp-location=europe-west2 \
  --keyring=keyring-name01 \
  --key=key-name01 \
  --bucket=my-gcs-bucket-name \
  --vault-url=https://vault.test.co.uk \
  --token-specs=/tmp/specs.json \
  --slack-webhook-url=https://hooks.slack.com/services/your-secret-webhook-url \
  --ttl-threshold=432000 \
  --ttl-increment=86400
```

## Deploying in vault-renovator in Kubernetes environment

1. Enable Google Cloud Key Management Service
2. Create key ring and key in the KMS section
3. In the GCP IAM section crate a new Service Account
4. Generate a new key for the previously created Service Account
5. Set proper IAM binding for the previously created SA ([source](https://codelabs.developers.google.com/codelabs/vault-on-gke/index.html?index=..%2F..%2Fcloud#5))

  ```bash
  gcloud kms keys add-iam-policy-binding ${GCP_KEY_NAME} \
    --location ${GCP_KEY_RING_LOCATION} \
    --keyring ${GCP_KEY_RING_NAME} \
    --member "serviceAccount:${SERVICE_ACCOUNT}" \
    --role roles/cloudkms.cryptoKeyEncrypterDecrypter
  ```

6. Create a new GCS bucket, set correct permissions for the previously created SA
7. Generate your tokens from the Vault CLI
8. Put each token in the single file
9. Encrypt these files

  ```bash
  gcloud kms encrypt \
    --key ${GCP_KEY_NAME} \
    --keyring ${GCP_KEY_RING_NAME} \
    --location ${GCP_KEY_RING_LOCATION} \
    --plaintext-file /tmp/plaintext-token.01.txt \
    --ciphertext-file /tmp/encrypted-token.01.bin
  ```

10. Upload **encrypted** files to the previously created GCS bucket
11. Create specification file

  ```json
  {
    "fileNames": [
      "encrypted-token.01.bin"
    ]
  }
  ```

12. Create a ConfigMap from previously created JSON file

  ```
  kubectl create configmap vault-renovator-config \
    --from-file=./file.json
  ```

13. Create a Secret from downloaded SA key

  ```
  
  ```
