terraform {
  # Store Terraform states in an S3 bucket
  backend "s3" {
    bucket  = "unique-bucket-name"
    key     = "vpc_nat/terraform.tfstate"
    region  = "us-west-2"
    profile = "default"
  }
}
