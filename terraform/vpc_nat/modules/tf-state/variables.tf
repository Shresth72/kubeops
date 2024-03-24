# Terraform State Variables
variable "bucket_name" {
  description = "Remote S3 bucket"
  type        = string
  validation {
    condition     = length(var.bucket_name) > 0
    error_message = "Bucket name cannot be empty"
  }
}

variable "table_name" {
  description = "DynamoDB table name"
  type        = string
  validation {
    condition     = length(var.table_name) > 0
    error_message = "Table name cannot be empty"
  }
}
