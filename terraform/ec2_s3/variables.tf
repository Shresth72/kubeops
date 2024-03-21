variable "aws_key_name" {
  default = "myKeyName"
}

variable "vpc_cidr" {
  description = "Value of the CIDR range for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "vpc_name" {
  description = "Value of the name for the VPC"
  type        = string
  default     = "MyVPC"
}

variable "subnet_cidr" {
  description = "Value of the subnet cidr for the VPC"
  type        = string
  default     = "10.0.1.0/24"
}

variable "subnet_name" {
  description = "Value of the name for the subnet"
  type        = string
  default     = "MyTestSubnet"
}

variable "igw_name" {
  description = "Value of the name for the internet gateway"
  type        = string
  default     = "MyTestIGW"
}

variable "ec2_instance_type" {
  type    = string
  default = "t2.micro"
}

variable "ec2_name" {
  description = "Value of the name for the EC2 instance"
  type        = string
  default     = "MyTestEC2"
}
