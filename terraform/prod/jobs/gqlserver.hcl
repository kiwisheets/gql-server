job "gql-server-prod" {
  datacenters = ["hetzner"]

  group "gql-server" {
    count = 2

    task "gql-server" {
      driver = "docker"

      config {
        image = "kiwisheets/gql-server:prod-${version}"

        volumes = [
          "secrets/db-password.secret:/run/secrets/db-password.secret",
          "secrets/jwt-secret-key.secret:/run/secrets/jwt-secret-key.secret",
          "secrets/hash-salt.secret:/run/secrets/hash-salt.secret"
        ]
      }

      env {
        APP_VERSION = "0.0.0"
        API_PATH = "/api/"
        PORT = 3000
        ENVIRONMENT = "production"
        POSTGRES_HOST = "$${NOMAD_UPSTREAM_IP_postgres}"
        POSTGRES_PORT = "$${NOMAD_UPSTREAM_PORT_postgres}"
        POSTGRES_DB = "kiwisheets"
        POSTGRES_USER = "kiwisheets"
        POSTGRES_PASSWORD_FILE = "/run/secrets/db-password.secret"
        POSTGRES_MAX_CONNECTIONS = 20
        REDIS_ADDRESS = "$${NOMAD_UPSTREAM_ADDR_redis}"
        JWT_SECRET_KEY_FILE = "/run/secrets/jwt-secret-key.secret"
        HASH_SALT = "/run/secrets/hash-salt.secret"
        HASH_MIN_LENGTH = 10
      }

      template {
        data = <<EOF
{{with secret "kv/data/prod"}}{{.Data.data.postgres_password}}{{end}}
        EOF
        destination = "secrets/db-password.secret"
      }

      template {
        data = <<EOF
{{with secret "kv/data/prod"}}{{.Data.data.jwt_secret}}{{end}}
        EOF
        destination = "secrets/jwt-secret-key.secret"
      }

      template {
        data = <<EOF
{{with secret "kv/data/prod"}}{{.Data.data.hash_salt}}{{end}}
        EOF
        destination = "secrets/hash-salt.secret"
      }

      vault {
        policies = ["gql-server-prod"]
      }

      resources {
        cpu    = 256
        memory = 256
      }
    }

    network {
      mode = "bridge"
      port "http" {
        to = 3000
      }
    }

    service {
      name = "gql-server-prod"
      port = "http"

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "gql-postgres-prod"
              local_bind_port = 5432
            }
            upstreams {
              destination_name = "gql-redis-prod"
              local_bind_port = 6379
            }
          }
        }
      }

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.gql-server-dev.rule=Host(`app.kiwisheets.com`) && PathPrefix(`/api/`)",
      ]

      check {
        type     = "http"
        path     = "/api/"
        interval = "2s"
        timeout  = "2s"
      }
    }
  }

  group "postgres" {
    count = 1

    volume "gql-postgres-prod" {
      type      = "csi"
      read_only = false
      source    = "${volume_id}"
    }

    task "postgres" {
      driver = "docker"

      volume_mount {
        volume      = "gql-postgres-prod"
        destination = "/var/lib/postgresql/data"
        read_only   = false
      }

      config {
        image = "postgres:latest"

        volumes = [
          "secrets/db-password.secret:/run/secrets/db-password.secret"
        ]
      }

      env {
        PGDATA = "/var/lib/postgresql/data/db"
        POSTGRES_DB = "kiwisheets"
        POSTGRES_USER = "kiwisheets"
        POSTGRES_PASSWORD_FILE = "/run/secrets/db-password.secret"
      }

      template {
        data = <<EOF
{{with secret "kv/data/prod"}}{{.Data.data.postgres_password}}{{end}}
        EOF
        destination = "secrets/db-password.secret"
      }

      vault {
        policies = ["gql-server-prod"]
      }
    }

    network {
      mode = "bridge"
    }

    service {
       name = "gql-postgres-prod"
       port = "5432"

       connect {
         sidecar_service {}
       }
     }
  }

  group "redis" {
    count = 1

    task "redis" {
      driver = "docker"

      config {
        image = "redis:latest"
      }
    }

    network {
      mode = "bridge"
    }

    service {
       name = "gql-redis-prod"
       port = "6379"

       connect {
         sidecar_service {}
       }
     }
  }
}
