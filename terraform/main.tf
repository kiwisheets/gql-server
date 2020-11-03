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
