job "main-server-dev" {
  datacenters = ["hetzner"]

  group "gql-server" {
    count = 1

    task "gql-server" {
      driver = "docker"

      config {
        image = "kiwisheets/gql-server:develop"
        port_map = {
          http = 3000
        }

        volumes = [
          "secrets/db-password.secret:/run/secrets/db-password.secret",
          "secrets/jwt-secret-key.secret:/run/secrets/jwt-secret-key.secret"
        ]
      }

      env {
        APP_VERSION = "0.0.0"
        API_PATH = "/api/"
        PORT = 3000
        ENVIRONMENT = "production"
        POSTGRES_HOST = "${NOMAD_UPSTREAM_IP_postgres}"
        POSTGRES_DB = "kiwisheets"
        POSTGRES_USER = "kiwisheets"
        POSTGRES_PASSWORD_FILE = "/run/secrets/db-password.secret"
        POSTGRES_MAX_CONNECTIONS = 20
        REDIS_ADDRESS = "${NOMAD_UPSTREAM_ADDR_redis}"
        JWT_SECRET_KEY_FILE = "/run/secrets/jwt-secret-key.secret"
      }

      template {
        data = <<EOF
          {{with secret "kv/data/dev"}}{{.Data.data.postgres_password}}{{end}}
        EOF
        destination = "secrets/db-password.secret"
      }

      template {
        data = <<EOF
          {{with secret "kv/data/dev"}}{{.Data.data.jwt_secret}}{{end}}
        EOF
        destination = "secrets/jwt-secret-key.secret"
      }

      vault {
        policies = ["gql-server-dev"]
      }

      resources {
        cpu    = 256
        memory = 256

        network {
          mbits = 10
          port  "http" {}
        }
      }

      service {
        name = "gql-server-dev"
        port = "http"

        connect {
          sidecar_service {
            proxy {
              upstreams {
                destination_name = "postgres"
                local_bind_port = 5432
              }
              upstreams {
                destination_name = "redis"
                local_bind_port = 6379
              }
            }
          }
        }

        tags = [
          "traefik.enable=true",
          "traefik.http.routers.gql-server-dev.rule=Host(`beta.kiwisheets.com`) && PathPrefix(`/api/`)",
        ]

        check {
          type     = "http"
          path     = "/"
          interval = "2s"
          timeout  = "2s"
        }
      }
    }
  }

  group "postgres" {
    count = 1

    task "postgres" {
      driver = "docker"

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
          {{with secret "kv/data/dev"}}{{.Data.data.postgres_password}}{{end}}
        EOF
        destination = "secrets/db-password.secret"
      }

      vault {
        policies = ["gql-server-dev"]
      }
    }
    
    network {
      mode = "bridge"
    }

    service {
       name = "postgres"
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
       name = "redis"
       port = "6379"

       connect {
         sidecar_service {}
       }
     }
  }
}
