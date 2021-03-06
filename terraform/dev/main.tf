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
  datacenter = var.datacenter
}

provider "vault" {}

provider "cloudflare" {}

resource "random_password" "postgres_password" {
  length  = 32
  special = false
}

resource "random_password" "hash_salt" {
  length  = 32
  special = false
}

resource "vault_generic_secret" "gql_server" {
  path = "secret/gql-server"

  data_json = jsonencode({
    postgres_password = random_password.postgres_password.result
    hash_salt         = random_password.hash_salt.result
  })
}

resource "vault_pki_secret_backend_role" "jwt" {
  backend         = "pki"
  name            = "jwt"
  ttl             = "31536000"
  max_ttl         = "31536000"
  allow_localhost = true
  key_type        = "ec"
  key_bits        = 256
  generate_lease  = true
}

resource "vault_pki_secret_backend_cert" "jwt" {
  backend     = "pki"
  name        = vault_pki_secret_backend_role.jwt.name
  common_name = "localhost"
}

resource "vault_generic_secret" "jwt_public" {
  path = "secret/jwt-public"

  data_json = jsonencode({
    key = vault_pki_secret_backend_cert.jwt.certificate
  })
}

resource "vault_generic_secret" "jwt_private" {
  path = "secret/jwt-private"

  data_json = jsonencode({
    key = vault_pki_secret_backend_cert.jwt.private_key
  })
}

resource "vault_policy" "gql_server" {
  name = "gql-server"

  policy = <<EOT
path "secret/gql-server" {
  capabilities = ["read"]
}
path "secret/data/gql-server" {
  capabilities = ["read"]
}
path "secret/jwt-private" {
  capabilities = ["read"]
}
path "secret/data/jwt-private" {
  capabilities = ["read"]
}
EOT
}

resource "nomad_job" "gql_server" {
  jobspec = templatefile("${path.module}/jobs/gqlserver.hcl", {
    datacenter          = var.datacenter
    image_tag           = var.image_tag
    instances           = var.instance_count
    cloudflared_version = "2021.2.2"
  })
  detach = false
}

resource "consul_intention" "gql-server-tunnel" {
  source_name      = "tunnel-gql-server"
  destination_name = "gql-server"
  action           = "allow"
}

resource "consul_intention" "rabbit" {
  source_name      = "gql-server"
  destination_name = "rabbitmq"
  action           = "allow"
}
