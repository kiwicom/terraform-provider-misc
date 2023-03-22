terraform {
  required_providers {
    kiwi = {
      source  = "kiwicom/kiwi"
    }
  }
}

resource "kiwi_claim_from_pool" "server_addresses" {
  pool = [
    "10.0.0.1",
    "10.0.0.2",
    "10.0.0.3",
    "10.0.0.4",
    "10.0.0.5",
    "10.0.0.6",
  ]
  claimers = [
    "server1",
    "server2",
    "server3",
    "server4",
  ]
}

output "server_addresses" {
  value = kiwi_claim_from_pool.output
}
