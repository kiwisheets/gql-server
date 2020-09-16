terraform {
  backend "remote" {
    organization = "KiwiSheets"

    workspaces {
      name = "KiwiSheets-GraphQL-Server-Dev"
    }
  }
}

provider "nomad" {
  ca_file = "nomad-ca.pem"
}

resource "nomad_job" "frontend" {
  jobspec = file("${path.module}/jobs/gqlserver.hcl")
  detach  = false
}
