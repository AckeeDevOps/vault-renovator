# Work In Progress 👨‍💻

This repository will become a Docker image for the automatic Vault token
renewal. I work hard to finish it ASAP as I really need this in my production infrastructure 😂

## Purpose
This code is meant to be running inside the Kubernetes (GKE) pod as a CronJob. This solution is exclusively using Google Cloud Platform, porting to different platforms would most likely require a great slaughter.

Please note that vault-renovator is not trying to separate you from the all the hassle with the  Vault. This tool does only one thing - it renews the specified Vault tokens based on the schedule specified in the CronJob manifest file.
