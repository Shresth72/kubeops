locals {
  bucket_name = "${var.project_name}-tf-state"
  table_name  = "${var.project_name}-tf-state-lock-table"
}
