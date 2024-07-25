#!/bin/bash

ORG_NAME=$TF_ORG_NAME
if [ -z $ORG_NAME ]
then
  ORG_NAME="octopus"
fi

TOKEN=$(cat ~/.terraform.d/credentials.tfrc.json| jq -r '.credentials."app.terraform.io".token')
if [ -z "$TOKEN" ]
then
  TOKEN=$TF_TOKEN_app_terraform_io
  if [ -z $TOKEN ]
  then
    echo 'please run `terraform login` or set the TF_TOKEN_app_terraform_io environment variable'
    exit 2
  fi
fi

VERSION="2.4.0"

curl \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
  --request POST \
  --data @platform.json \
  https://app.terraform.io/api/v2/organizations/$ORG_NAME/registry-providers/private/$ORG_NAME/nexus/versions/$VERSION/platforms

