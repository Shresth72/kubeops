provider "aws" {
  region  = var.region
  profile = "default"
}

module "vpc" {
  source                       = "../modules/vpc"
  region                       = var.region
  project_name                 = var.project_name
  vpc_cidr                     = var.vpc_cidr
  public_subnet_az1_cidr       = var.public_subnet_az1_cidr
  public_subnet_az2_cidr       = var.public_subnet_az2_cidr
  private_app_subnet_az1_cidr  = var.private_app_subnet_az1_cidr
  private_app_subnet_az2_cidr  = var.private_app_subnet_az2_cidr
  private_data_subnet_az1_cidr = var.private_data_subnet_az1_cidr
  private_data_subnet_az2_cidr = var.private_data_subnet_az2_cidr
}

module "nat_gateway" {
  source                     = "../modules/nat_gateway"
  vpc_id                     = var.vpc.vpc_id
  internet_gateway           = var.internet_gateway
  public_subnet_az1_id       = var.vpc.public_subnet_az1_id
  public_subnet_az2_id       = var.vpc.public_subnet_az2_id
  private_app_subnet_az1_id  = var.vpc.private_app_subnet_az1_id
  private_app_subnet_az2_id  = var.vpc.private_app_subnet_az2_id
  private_data_subnet_az1_id = var.vpc.private_data_subnet_az1_id
  private_data_subnet_az2_id = var.vpc.private_data_subnet_az2_id
}
