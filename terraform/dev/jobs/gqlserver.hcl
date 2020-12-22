job "gql-server" {
  datacenters = ["${datacenter}"]

  update {
    auto_revert       = true
    auto_promote      = true
    canary            = 1
    max_parallel      = 3
    min_healthy_time  = "1m"
    healthy_deadline  = "5m"
    progress_deadline = "10m"
  }

  group "gql-server" {
    count = 2

    task "gql-server" {
      driver = "docker"

      config {
        image = "kiwisheets/gql-server:${image_tag}"

        volumes = [
          "secrets/db-password.secret:/run/secrets/db-password.secret",
          "secrets/jwt-private-key.secret:/run/secrets/jwt-private-key.secret",
          "secrets/hash-salt.secret:/run/secrets/hash-salt.secret"
        ]
      }

      env {
        APP_VERSION = "0.0.0"
        API_PATH = "/graphql"
        PORT = 3000
        ENVIRONMENT = "production"
        POSTGRES_HOST = "$${NOMAD_UPSTREAM_IP_gql-postgres}"
        POSTGRES_PORT = "$${NOMAD_UPSTREAM_PORT_gql-postgres}"
        POSTGRES_DB = "kiwisheets"
        POSTGRES_USER = "kiwisheets"
        POSTGRES_PASSWORD_FILE = "/run/secrets/db-password.secret"
        POSTGRES_MAX_CONNECTIONS = 20
        REDIS_ADDRESS = "$${NOMAD_UPSTREAM_ADDR_gql-redis}"
        JWT_EC_PRIVATE_KEY_FILE = "/run/secrets/jwt-private-key.secret"
        HASH_SALT_FILE = "/run/secrets/hash-salt.secret"
        HASH_MIN_LENGTH = 10
      }

      template {
        data = "{{with secret \"secret/data/gql-server\"}}{{.Data.data.postgres_password}}{{end}}"
        destination = "secrets/db-password.secret"
      }

      template {
        data = "{{with secret \"secret/data/jwt-private\"}}{{.Data.data.key}}{{end}}"
        destination = "secrets/jwt-private-key.secret"
        change_mode = "noop"
      }

      template {
        data = "{{with secret \"secret/data/gql-server\"}}{{.Data.data.hash_salt}}{{end}}"
        destination = "secrets/hash-salt.secret"
      }

      vault {
        policies = ["gql-server"]
      }

      resources {
        cpu    = 256
        memory = 256
      }
    }

    network {
      mode = "bridge"
      port "health" {}
    }

    service {
      name = "gql-server"
      port = 3000

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "gql-postgres"
              local_bind_port = 5432
            }
            upstreams {
              destination_name = "gql-redis"
              local_bind_port = 6379
            }
            expose {
              path {
                path           = "/health"
                protocol        = "http"
                local_path_port = 3000
                listener_port   = "health"
              }
            }
          }
        }

        sidecar_task {
          resources {
            cpu    = 20
            memory = 32
          }
        }
      }

      check {
        type     = "http"
        path     = "/health"
        port     = "health"
        interval = "2s"
        timeout  = "2s"
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
        data = "{{with secret \"secret/data/gql-server\"}}{{.Data.data.postgres_password}}{{end}}"
        destination = "secrets/db-password.secret"
      }

      vault {
        policies = ["gql-server"]
      }
    }
    
    network {
      mode = "bridge"
    }

    service {
       name = "gql-postgres"
       port = "5432"

       connect {
         sidecar_service {}

         sidecar_task {
          resources {
            cpu    = 20
            memory = 32
          }
        }
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
      name = "gql-redis"
      port = "6379"

      connect {
        sidecar_service {}

        sidecar_task {
          resources {
            cpu    = 20
            memory = 32
          }
        }
      }
    }
  }
}
