variable "project_name" {
  type        = string
  description = "Globally used project name for this demo"
  default     = "eks_managed_nodes"
  validation {
    condition = (
      length(var.project_name) > 4 &&
      length(var.project_name) < 20 &&
      replace(var.project_name, " ", "") == var.project_name
    )
    error_message = "The project name must be alphanumeric and dashes only"
  }
}

variable "region" {
  type        = string
  description = "The AWS region to deploy to"
  default     = "us-west-2"
}

variable "cluster_version" {
  type        = string
  description = "The Kubernetes version for our clusters"
  default     = "1.22"
}

variable "cluster_instance_type" {
  type        = string
  description = "EC2 instance type for the EKS autoscaling group."
  default     = "t3.micro"
}

variable "cluster_asg_desired_capacity" {
  type        = number
  description = "The default number of EC2 instances our EKS cluster runs."
  default     = 1
}

variable "cluster_asg_min_size" {
  type        = number
  description = "The minimum number of EC2 instances our EKS cluster will have."
  default     = 1
}

variable "cluster_asg_max_size" {
  type        = number
  description = "The maximum number of EC2 instances our EKS cluster will have."
  default     = 3
}

variable "cluster_enabled_log_types" {
  type        = list(string)
  description = "The Kubernetes log types that will be enabled for the EKS cluster."
  default     = ["api", "audit", "authenticator", "controllerManager", "scheduler"]
}

variable "cluster_write_kubeconfig" {
  type        = bool
  description = "Specify if Terraform sound output the Kubernetes configuration file. "
  default     = false
}

variable "cidr" {
  type        = string
  description = "The CIDR block to use for this demo"
  default     = "10.0.0.0/16"
}
