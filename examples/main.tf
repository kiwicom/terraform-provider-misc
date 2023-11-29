terraform {
  required_providers {
    kiwi = {
      source = "kiwicom/misc"
    }
  }
}

resource "misc_claim_from_pool" "master_cidr_subnets" {
  pool = [for i in range(900, 1024) : cidrsubnet("172.23.192.0/18", 10, i)]
  claimers = [
    "cluster1",
    "cluster2",
    "cluster3",
    "cluster4",
  ]
}

output "master_cidr_subnets" {
  value = misc_claim_from_pool.master_cidr_subnets.output
}
