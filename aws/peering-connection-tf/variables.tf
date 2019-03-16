variable "peering_connection_name" {
  description = "The name of the peering connection.s"
}

variable "accepter_vpc_id" {
  description = "The ID of the VPC with which you are creating the VPC Peering Connection."
}

variable "accepter_route_table_ids" {
  description = "ID of the route tables you would like the routes to be added to. By default, it's added to all route tables."
  default     = []
  type        = "list"
}

variable "requester_vpc_id" {
  description = "The ID of the VPC requesting the peering connection."
}

variable "requester_route_table_ids" {
  description = "ID of the route tables you would like the routes to be added to. By default, it's added to all route tables."
  default     = []
  type        = "list"
}

variable "requester_cidr_block" {
  description = "Requester CIDR block to be peered."
  default = ""
}

variable "accepter_cidr_block" {
  description = "Accepter CIDR block to be peered."
  default = ""
}