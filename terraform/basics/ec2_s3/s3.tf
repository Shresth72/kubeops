
# create s3 bucket
resource "aws_s3_bucket" "my_terraform-bucket" {
  bucket = "my-terraform-bucket"
  tags = {
    Name = "my-terraform-bucket"
  }
}

resource "aws_s3_bucket_ownership_controls" "bucket_ownership_controls" {
  bucket = aws_s3_bucket.my_terraform-bucket.bucket
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "bucket_acl" {
  depends_on = [aws_s3_bucket_ownership_controls.bucket_ownership_controls]

  bucket = aws_s3_bucket.my_terraform-bucket.id
  acl    = "private"
}

resource "aws_s3_bucket_policy" "bucket_policy" {
  bucket = aws_s3_bucket.my_terraform-bucket.id
  policy = data.aws_iam_policy_document.s3_read_access.json
}