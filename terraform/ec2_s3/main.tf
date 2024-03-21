data "aws_ami" "aws_linux_2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "owner-alias"
    values = ["amazon"]
  }

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm"]
  }
}

# EC2 instance
resource "aws_instance" "ec2_instance" {
  ami           = data.aws_ami.aws_linux_2.id
  instance_type = var.ec2_instance_type
  key_name      = var.aws_key_name
  iam_instance_profile = aws_iam_instance_profile.ec2_instance_profile.name

  subnet_id                   = aws_subnet.my_subnet.id
  vpc_security_group_ids      = [aws_security_group.ec2_security_group.id]
  associate_public_ip_address = true
  #user_data                   = file("run.sh")

  tags = {
    Name = var.ec2_name
  }
}
