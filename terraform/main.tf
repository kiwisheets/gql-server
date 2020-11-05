terraform {
  backend "remote" {
    organization = "KiwiSheets"

    workspaces {
      prefix = "KiwiSheets-GraphQL-Server-"
    }
  }
}

provider "hcloud" {
  token = var.hcloud_token
}

provider "nomad" {}

provider "consul" {
  datacenter = "hetzner"
}

data "nomad_plugin" "hcloud_volume" {
  plugin_id        = "hcloud-volume"
  wait_for_healthy = true
}

resource "hcloud_volume" "gql_postgres" {
  name     = "gql-postgres-${var.environment}"
  size     = 12
  location = "nbg1"
}

resource "nomad_volume" "gql_postgres" {
  depends_on            = [data.nomad_plugin.hcloud_volume]
  type                  = "csi"
  plugin_id             = "hcloud-volume"
  volume_id             = hcloud_volume.gql_postgres.name
  name                  = hcloud_volume.gql_postgres.name
  external_id           = hcloud_volume.gql_postgres.id
  access_mode           = "single-node-writer"
  attachment_mode       = "file-system"
  deregister_on_destroy = true
}

resource "nomad_job" "gql_server" {
  jobspec = templatefile("${path.module}/jobs/gqlserver.hcl", {
    image_tag       = var.image_tag
    env             = var.environment
    allowed_origins = var.allowed_origins
    instance        = var.instance_count
    host            = var.host
    volume_id       = nomad_volume.gql_postgres.volume_id
  })
  detach = false
}
