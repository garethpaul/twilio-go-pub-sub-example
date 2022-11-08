# Twilio-Go Pub Sub Example
## What is it?
This repo contains Terraform and Golang code for showcasing of Pub/Sub message flow using IaC and principle of least privilege on GCP with Twilio Go.

## Deploying:
Export `PROJECT_ID` as environmental variable before running docker-compose build and push. The solution assumes you're utilizing GCR as your container registry and you should enable GCR API in the project on beforehand.

Terraform requires following variables to be passed:
```
TF_VAR_project_id=<your-project-id>
```

## Running:
`terraform init`
`terraform plan`

WARNING: doing this will spend real money so be careful
`terraform apply`


## Testing it out
Curl the publisher API with following command:
```
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" -X POST "$(terraform output -raw publisher_url)/generate-messages"ex
```

You should then see data flow from the publisher to the processor.