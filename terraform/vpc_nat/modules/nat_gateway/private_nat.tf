/* --- Private AZ1 --- */

resource "aws_route_table" "private_route_table_az1" {
  vpc_id = var.vpc_id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat_gateway_az1.id
  }

  tags = {
    Name = "Private Route Table AZ1"
  }
}

resource "aws_route_table_association" "private_app_subnet_az1_association" {
  subnet_id      = var.private_app_subnet_az1_id
  route_table_id = aws_route_table.private_route_table_az1.id
}

resource "aws_route_table_association" "private_data_subnet_az1_association" {
  subnet_id      = var.private_data_subnet_az1_id
  route_table_id = aws_route_table.private_route_table_az1.id
}

/* --- Private AZ2 --- */

resource "aws_route_table" "private_route_table_az2" {
  vpc_id = var.vpc_id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat_gateway_az2.id
  }

  tags = {
    Name = "Private Route Table AZ2"
  }
}

resource "aws_route_table_association" "private_app_subnet_az2_association" {
  subnet_id      = var.private_app_subnet_az2_id
  route_table_id = aws_route_table.private_route_table_az2.id
}

resource "aws_route_table_association" "private_data_subnet_az2_association" {
  subnet_id      = var.private_data_subnet_az2_id
  route_table_id = aws_route_table.private_route_table_az2.id
}