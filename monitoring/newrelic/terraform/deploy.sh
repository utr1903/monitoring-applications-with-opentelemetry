#!/bin/bash

# Get commandline arguments
while (( "$#" )); do
  case "$1" in
    --dry-run)
      flagDryRun="true"
      shift
      ;;
    --destroy)
      flagDestroy="true"
      shift
      ;;
    *)
      shift
      ;;
  esac
done


# Initialize Terraform
terraform init

if [[ $flagDestroy != "true" ]]; then

  # Plan Terraform
  terraform plan \
    -var NEW_RELIC_ACCOUNT_ID=<YOUR_ACCOUNT_ID> \
    -var NEW_RELIC_API_KEY="<YOUR_API_KEY>" \
    -var NEW_RELIC_REGION="<YOUR_REGION>" \
    -out "./tfplan"

  # Apply Terraform
  if [[ $flagDryRun != "true" ]]; then
    terraform apply \
      -auto-approve \
      tfplan
  fi
else

  # Destroy Terraform
  terraform destroy \
    -auto-approve \
    -var NEW_RELIC_ACCOUNT_ID=<YOUR_ACCOUNT_ID> \
    -var NEW_RELIC_API_KEY="<YOUR_API_KEY>" \
    -var NEW_RELIC_REGION="<YOUR_REGION>"
fi
