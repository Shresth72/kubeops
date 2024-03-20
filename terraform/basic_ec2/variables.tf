variable "instance_name" {
  description = "Value of the name tag for the EC2 instance"
  type        = string
  default     = "MyTerraformInstance"
} # terraform apply -var 'instance_name=MyInstanceName'

variable "ec2_instance_type" {
  description = "AWS EC2 instance type"
  type        = string
  default     = "t2.micro"
}