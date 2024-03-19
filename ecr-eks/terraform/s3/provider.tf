terraform {
    backend "s3" {
        bucket = "terraform-remote-state-ecr-eks"
        key    = "terraform.tfstate"
        region = "us-west-2"
    }
}