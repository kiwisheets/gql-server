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
    image_tag       = var.image_tag
    env             = var.environment
    allowed_origins = var.allowed_origins
    instance        = var.instance_count
    host            = var.host
  })
  detach = false
}
