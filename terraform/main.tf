terraform {
  backend "remote" {
    organization = "KiwiSheets"

    workspaces {
      name = "KiwiSheets-GraphQL-Server-Dev"
    }
  }
}

provider "nomad" {}

provider "consul" {
  ca_file    = var.consul_ca_file
  cert_file  = var.consul_cert_file
  key_file   = var.consul_key_file
  address    = var.consul_address
  datacenter = "hetzner"
}

resource "nomad_job" "gql_server" {
  jobspec = templatefile("${path.module}/jobs/gqlserver.hcl", {
    tag           = var.image_tag
    env           = var.environment
    instance      = var.instance_count
    domain_prefix = var.environment == "prod" ? "app" : "beta"
  })
  detach = false
}
