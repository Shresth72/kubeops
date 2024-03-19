# Github Actions Terraform Workflow

## Setup IAM Provider

- Setup an OpenId Connect IAM Provider.
- Set `Provider URL` to `https://token.actions.githubusercontent.com`
- Set `Audience` to `sts.amazonaws.com`

## Setup S3 Bucket to store Terraform States

- Create S3 Bucket `github-oidc-terraform-aws-tfstates`
- Activate Encryption with `managed keys`

## Setup IAM Role for Github to Assume

- Create Role `github-oidc-terraform-aws-tfstates-role` with Custom Trust Policy

```json
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::YOUR_ACCOUNT_NUMBER:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringLike": {
          "token.actions.githubusercontent.com:sub": "repo:YOUR_GITHUB_USERNAME/YOUR_REPO_NAME:*"
        }
      }
    }
  ]
}
```

- Create policy `github-oidc-terraform-aws-tfstates-access` in the role for accessing S3 bucket to access terraform states

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:PutObject", "s3:GetObject", "s3:ListBucket"],
      "Resource": ["arn:aws:s3:::YOUR_BUCKET/*", "arn:aws:s3:::YOUR_BUCKET"]
    }
  ]
}
```

- Create another policy `github-oidc-terraform-aws-ssm-parameter-store` for deployment parameters (specific for each use case)

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "VisualEditor0",
      "Effect": "Allow",
      "Action": [
        "ssm:PutParameter",
        "ssm:LabelParameterVersion",
        "ssm:DeleteParameter",
        "ssm:UnlabelParameterVersion",
        "ssm:DescribeParameters",
        "ssm:GetParameterHistory",
        "ssm:ListTagsForResource",
        "ssm:GetParametersByPath",
        "ssm:GetParameters",
        "ssm:GetParameter",
        "ssm:DeleteParameters"
      ],
      "Resource": "*"
    }
  ]
}
```

- Save the ARN path for github actions