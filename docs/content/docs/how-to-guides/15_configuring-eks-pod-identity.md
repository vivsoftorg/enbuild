---
title: Configuring EKS Pod Identity for AWS Services
description: How to configure EKS Pod Identity for ENBUILD AI and CTF services to access AWS services like Bedrock and S3.
---

# Configuring EKS Pod Identity for AWS Services

This guide explains how to configure EKS Pod Identity to enable ENBUILD's AI service to access AWS Bedrock and CTF service to access AWS S3.

## Overview

EKS Pod Identity provides fine-grained AWS IAM permissions to Kubernetes pods without requiring node-level permissions or OIDC provider setup. This is the recommended approach for granting AWS access to ENBUILD services running on EKS.

## Prerequisites

- EKS cluster with Pod Identity Agent addon installed
- AWS CLI configured with appropriate permissions
- Helm chart with service account support (v0.0.38+)

## Step 1: Enable Service Accounts in Helm Chart

Update your values file to enable service accounts:

```yaml
lightning_features:
  deploy_lightning:
    ai_lightning: true
  secure_lightning:
    ctf: true

enbuildAI:
  serviceAccount:
    create: true

enbuildCTF:
  serviceAccount:
    create: true
```

Deploy or upgrade the chart:

```bash
helm upgrade --install enbuild enbuild/enbuild \
  --namespace enbuild \
  -f values.yaml
```

## Step 2: Ensure EKS Pod Identity Agent is Installed

Check if the Pod Identity Agent is installed:

```bash
kubectl get pods -n kube-system | grep pod-identity
```

If not installed, add it:

```bash
aws eks create-addon \
  --cluster-name your-cluster-name \
  --addon-name eks-pod-identity-agent
```

## Step 3: Create IAM Roles

### Create Role for AI (Bedrock Access)

```bash
# Create IAM role with Pod Identity trust policy
aws iam create-role \
  --role-name enbuild-ai-bedrock \
  --assume-role-policy-document '{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": { "Service": "pods.eks.amazonaws.com" },
        "Action": ["sts:AssumeRole", "sts:TagSession"]
      }
    ]
  }'
```

Attach Bedrock permissions:

```bash
# Option 1: Full Bedrock access
aws iam attach-role-policy \
  --role-name enbuild-ai-bedrock \
  --policy-arn arn:aws:iam::aws:policy/AmazonBedrockFullAccess

# Option 2: Specific model access (recommended)
aws iam create-policy \
  --policy-name enbuild-ai-bedrock-specific \
  --policy-document '{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "bedrock:InvokeModel",
          "bedrock:ListFoundationModels"
        ],
        "Resource": "arn:aws:bedrock:*::foundation-model/*"
      }
    ]
  }'

aws iam attach-role-policy \
  --role-name enbuild-ai-bedrock \
  --policy-arn arn:aws:iam::123456789012:policy/enbuild-ai-bedrock-specific
```

### Create Role for CTF (S3 Access)

```bash
aws iam create-role \
  --role-name enbuild-ctf-s3 \
  --assume-role-policy-document '{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": { "Service": "pods.eks.amazonaws.com" },
        "Action": ["sts:AssumeRole", "sts:TagSession"]
      }
    ]
  }'

# Attach S3 permissions
aws iam attach-role-policy \
  --role-name enbuild-ctf-s3 \
  --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess
```

## Step 4: Create Pod Identity Associations

Associate the IAM roles with Kubernetes service accounts:

```bash
# For AI service (Bedrock)
aws eks create-pod-identity-association \
  --cluster-name your-cluster-name \
  --namespace enbuild \
  --service-account enbuild-ai \
  --role-arn arn:aws:iam::123456789012:role/enbuild-ai-bedrock

# For CTF service (S3)
aws eks create-pod-identity-association \
  --cluster-name your-cluster-name \
  --namespace enbuild \
  --service-account enbuild-ctf \
  --role-arn arn:aws:iam::123456789012:role/enbuild-ctf-s3
```

## Step 5: Verify Configuration

### Check Pod Identity Associations

```bash
aws eks list-pod-identity-associations --cluster-name your-cluster-name
```

### Verify Pod Has Credentials

After redeploying the pods, verify the AWS credentials are injected:

```bash
# Check environment variables in the pod
kubectl exec -it deployment/enbuild-enbuild-ai -- env | grep AWS

# Should show:
# AWS_ROLE_ARN=arn:aws:iam::123456789012:role/enbuild-ai-bedrock
# AWS_WEB_IDENTITY_TOKEN_FILE=/var/run/secrets/eks.amazonaws.com/serviceaccount/token
```

### Test AWS Access

```bash
# Test Bedrock access from AI pod
kubectl exec -it deployment/enbuild-enbuild-ai -- \
  aws bedrock list-foundation-models --region us-east-1
```

## Troubleshooting

### Pods Not Receiving Credentials

1. Check the service account name matches exactly:
   ```bash
   kubectl get sa enbuild-ai -n enbuild
   ```

2. Verify Pod Identity association:
   ```bash
   aws eks describe-pod-identity-association \
     --cluster-name your-cluster-name \
     --association-id <association-id>
   ```

### Access Denied Errors

1. Verify the IAM role has the correct trust policy
2. Check the attached policies grant the required permissions
3. Ensure the role ARN matches the association

### Credentials Not Rotating

Pod Identity credentials auto-rotate. If you encounter issues:

1. Restart the pods to get new credentials:
   ```bash
   kubectl rollout restart deployment/enbuild-enbuild-ai
   ```

## Disabling hostNetwork Mode

Previously, the CTF service used `hostNetwork: true` to access AWS via node metadata. This is now disabled in favor of Pod Identity. If you need to revert:

```yaml
# In values.yaml
enbuildCTF:
  hostNetwork: true  # Not recommended - use Pod Identity instead
```

## Additional Resources

- [EKS Pod Identity Documentation](https://docs.aws.amazon.com/eks/latest/userguide/pod-identities.html)
- [IAM Roles for Service Accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html)
- [Amazon Bedrock Permissions](https://docs.aws.amazon.com/bedrock/latest/userguide/security-iam.html)
