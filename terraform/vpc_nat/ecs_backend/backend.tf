terraform {
  required_version = "~> 1.3"

  backend "s3" {
    bucket         = locals.bucket_name
    key            = "tf-infra/terraform.tfstate"
    region         = "us-west-2"
    profile        = "default"
    dynamodb_table = locals.table_name
    encrypt        = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}
