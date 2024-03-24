# VPC Module Variables
variable "region" {}
variable "project_name" {}

variable "vpc_cidr" {}

variable "public_subnet_az1_cidr" {}
variable "public_subnet_az2_cidr" {}

variable "private_app_subnet_az1_cidr" {}
variable "private_app_subnet_az2_cidr" {}

variable "private_data_subnet_az1_cidr" {}
variable "private_data_subnet_az2_cidr" {}

# NAT Gateway Module Variables
variable "internet_gateway" {}
variable "vpc_id" {}

variable "public_subnet_az1_id" {}
variable "public_subnet_az2_id" {}

variable "private_app_subnet_az1_id" {}
variable "private_app_subnet_az2_id" {}

variable "private_data_subnet_az1_id" {}
variable "private_data_subnet_az2_id" {}

