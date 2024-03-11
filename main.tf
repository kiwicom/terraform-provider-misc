terraform {
  required_providers {
    misc = {
      source = "kiwicom/misc"
    }
  }
}

resource "misc_cert_pack_squash" "bla" {
  hosts = [
    "host.com",
    "asdasdas.asdsadasd.com",
    "*.asdsadasd.com",
    "asdasd.asdasdas.asdsadasd.com",
    "host.com",
  ]
}

output "cert" {
  value = misc_cert_pack_squash.bla
}


resource "kustomization_resource" "" {
  manifest = ""
}