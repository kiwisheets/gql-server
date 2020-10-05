terraform {
  backend "remote" {
    organization = "KiwiSheets"

    workspaces {
      name = "KiwiSheets-GraphQL-Server-Prod"
    }
  }
}

provider "hcloud" {
  token = var.hcloud_token
}

provider "nomad" {
  ca_file = "nomad-ca.pem"
}

data "nomad_plugin" "hcloud_volume" {
  plugin_id        = "hcloud-volume"
  wait_for_healthy = true
}

resource "hcloud_volume" "gql_postgres_prod" {
  name     = "gql-postgres-prod"
  size     = 10
  location = "nbg1"
}

resource "nomad_volume" "gql_postgres_prod" {
  depends_on            = [data.nomad_plugin.hcloud_volume]
  type                  = "csi"
  plugin_id             = "hcloud-volume"
  volume_id             = hcloud_volume.gql_postgres_prod.name
  name                  = hcloud_volume.gql_postgres_prod.name
  external_id           = hcloud_volume.gql_postgres_prod.id
  access_mode           = "single-node-writer"
  attachment_mode       = "file-system"
  deregister_on_destroy = true
}

resource "nomad_job" "gql-server" {
  jobspec = templatefile("${path.module}/jobs/gqlserver.hcl", {
    version   = var.image_version
    volume_id = nomad_volume.gql_postgres_prod.volume_id
  })
  detach = false
}
