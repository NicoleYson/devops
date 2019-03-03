######### Peering Connection #########
resource "aws_vpc_peering_connection" "pcx" {
  provider = "aws.requester"

  peer_owner_id = "${data.aws_caller_identity.accepter.account_id}"
  peer_vpc_id   = "${var.accepter_vpc_id}"
  vpc_id        = "${var.requester_vpc_id}"

  tags {
    Name      = "${var.peering_connection_name}"
    Side      = "Requester"
    Terraform = true
  }
}

data "aws_caller_identity" "requester" {
  provider = "aws.requester"
}

######### Route Tables #########
data "aws_route_tables" "req_rt" {
  provider = "aws.requester"
  vpc_id   = "${var.requester_vpc_id}"
}

# To get cidr block for route
data "aws_vpc" "accepter" {
  provider = "aws.accepter"
  id       = "${var.accepter_vpc_id}"
}

# Places this peering connection on all route tables in the VPC
resource "aws_route" "all_route_tables" {
  provider = "aws.requester"

  # If the user does not define any route table ids, then this data resource is utilized. Otherwise, this resource isn't used at all.
  count                     = "${length(var.requester_route_table_ids) == 0 ? length(data.aws_route_tables.req_rt.ids) : 0}"
  route_table_id            = "${data.aws_route_tables.req_rt.ids[count.index]}"
  destination_cidr_block    = "${data.aws_vpc.accepter.cidr_block}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.pcx.id}"
}

# This places the peering connection on route tables specified by the user.
resource "aws_route" "defined_route" {
  provider       = "aws.requester"
  count          = "${length(var.requester_route_table_ids)}"
  route_table_id = "${element(var.requester_route_table_ids, count.index)}"

  destination_cidr_block    = "${data.aws_vpc.accepter.cidr_block}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.pcx.id}"
}
