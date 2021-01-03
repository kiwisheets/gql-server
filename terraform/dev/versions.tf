terraform {
  required_providers {
    nomad = {
      source = "hashicorp/nomad"
    }
    consul = {
      source = "hashicorp/consul"
    }
    vault = {
      source = "hashicorp/vault"
    }
    cloudflare = {
      source = "cloudflare/cloudflare"
    }
  }
  required_version = ">= 0.13"
}
