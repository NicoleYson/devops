### Peering Connection
This module is used to create peering connections between VPCs.
By default, this adds the peering connection to all the route tables in the given VPC.
This isn't necessarily reccomended and can be specified as route_table_ids and should be of type list.

### Example Usage


```hcl
data "aws_vpc" "vpc1" {
  filter {
    name = "tag:Name"
    values = ["<NAME-OF-YOUR-VPC>"]
  }
}

data "aws_vpc" "vpc2" {
  filter {
    name = "tag:Name"
    values = ["<NAME-OF-YOUR-VPC>"]
  }

  provider = "aws.<ALIAS-NAME>"
}

module "pcx" {
  source                  = "github.com/nicoleyson/devops/aws/peering-connection-tf"
  peering_connection_name = "vpc 1 to vpc2"
  accepter_vpc_id         = "${data.aws_vpc.vpc2.id}"
  accepter_route_table_ids = "[optional, if blank will add to all]"
  accepter_cidr_block = "[optional, if blank will peer entire vpc]"
  requester_vpc_id        = "${data.aws_vpc.vpc1.id}"
  requester_route_table_ids = "[optional, if blank will add to all]"
  requester_cidr_block = "[optional, if blank will peer entire vpc]"

  providers = {
    "aws.accepter"  = "aws.<ALIAS-NAME>"
    "aws.requester" = "aws"
  }
}
```

