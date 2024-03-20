# Provider Block
provider "aws" {
  profile = "default"
  region  = "ap-southeast-1"
}

# Resource Block
resource "aws_instance" "app_server" {
  ami           = "ami-0c55b159cbfafe1f0" // aws linux 2
  instance_type = var.ec2_instance_type

  tags = {
    Name = var.instance_name
  }
}

