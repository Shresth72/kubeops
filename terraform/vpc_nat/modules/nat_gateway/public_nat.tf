# Elastic Ip for public subnet az1
resource "aws_eip" "eip_for_nat_gateway_az1" {
  domain = "vpc"
  tags = {
    Name = "EIP for NAT Gateway AZ1"
  }
}

# Elastic Ip for public subnet az2
resource "aws_eip" "eip_for_nat_gateway_az2" {
  domain = "vpc"
  tags = {
    Name = "EIP for NAT Gateway AZ2"
  }
}
resource "aws_nat_gateway" "nat_gateway_az1" {
  allocation_id = aws_eip.eip_for_nat_gateway_az1.id
  subnet_id     = var.public_subnet_az1_id

  tags = {
    Name = "NAT Gateway AZ1"
  }

  depends_on = [var.internet_gateway]
}

resource "aws_nat_gateway" "nat_gateway_az2" {
  allocation_id = aws_eip.eip_for_nat_gateway_az2.id
  subnet_id     = var.public_subnet_az2_id

  tags = {
    Name = "NAT Gateway AZ2"
  }

  depends_on = [var.internet_gateway]
}
