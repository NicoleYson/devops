output "pcx_id_accepter" {
  value = "${aws_vpc_peering_connection_accepter.pcx.id}"
}

output "pcx_id_requester" {
  value = "${aws_vpc_peering_connection.pcx.id}"
}
