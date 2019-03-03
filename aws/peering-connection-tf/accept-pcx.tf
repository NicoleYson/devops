######### Peering Connection #########
data "aws_caller_identity" "accepter" {
  provider = "aws.accepter"
}

resource "aws_vpc_peering_connection_accepter" "pcx" {
  provider                  = "aws.accepter"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.pcx.id}"
  auto_accept               = true

  tags {
    Side      = "Accepter"
    Name      = "${var.peering_connection_name}"
    Terraform = true
  }
}

######### Route Tables #########
data "aws_route_tables" "acc_rt" {
  provider = "aws.requester"
  vpc_id   = "${var.requester_vpc_id}"
}

# To get cidr block for route

data "aws_vpc" "requester" {
  provider = "aws.requester"
  id       = "${var.requester_vpc_id}"
}

# This places the peering connection on route tables specified by the user.
resource "aws_route" "defined_route" {
  provider       = "aws.accepter"
  count          = "${length(var.accepter_accepter_route_table_ids)}"
  route_table_id = "${element(var.accepter_accepter_route_table_ids, count.index)}"

  destination_cidr_block    = "${data.aws_vpc.accepter.cidr_block}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.pcx.id}"
}

# Places this peering connection on all route tables in the VPC
resource "aws_route" "all_route_tables" {
  provider = "aws.accepter"

  # If the user does not define any route table ids, then this data resource is utilized. Otherwise, this resource isn't used at all.
  count                     = "${length(var.accepter_route_table_ids) == 0 ? length(data.aws_route_tables.acc_rt.ids) : 0}"
  route_table_id            = "${data.aws_route_tables.acc_rt.ids[count.index]}"
  destination_cidr_block    = "${data.aws_vpc.accepter.cidr_block}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.pcx.id}"
}
