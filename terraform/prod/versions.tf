terraform {
  required_providers {
    hcloud = {
      source = "hetznercloud/hcloud"
    }
    nomad = {
      source = "hashicorp/nomad"
    }
  }
  required_version = ">= 0.13"
}
