job "gql-server-${env}" {
  datacenters = ["hetzner"]

  group "gql-server" {
    count = 1

    task "gql-server" {
      driver = "docker"

      config {
        image = "kiwisheets/gql-server:${image_tag}"

        volumes = [
          "secrets/db-password.secret:/run/secrets/db-password.secret",
          "secrets/jwt-secret-key.secret:/run/secrets/jwt-secret-key.secret",
          "secrets/hash-salt.secret:/run/secrets/hash-salt.secret"
        ]
      }

      env {
        APP_VERSION = "0.0.0"
        API_PATH = "/graphql"
        ALLOWED_ORIGINS = "${allowed_origins}"
        PORT = 3000
        ENVIRONMENT = "production"
        POSTGRES_HOST = "$${NOMAD_UPSTREAM_IP_gql-postgres-${env}}"
        POSTGRES_PORT = "$${NOMAD_UPSTREAM_PORT_gql-postgres-${env}}"
        POSTGRES_DB = "kiwisheets"
        POSTGRES_USER = "kiwisheets"
        POSTGRES_PASSWORD_FILE = "/run/secrets/db-password.secret"
        POSTGRES_MAX_CONNECTIONS = 20
        REDIS_ADDRESS = "$${NOMAD_UPSTREAM_ADDR_gql-redis-${env}}"
        JWT_SECRET_KEY_FILE = "/run/secrets/jwt-secret-key.secret"
        HASH_SALT_FILE = "/run/secrets/hash-salt.secret"
        HASH_MIN_LENGTH = 10
      }

      template {
        data = <<EOF
{{with secret "kv/data/${env}"}}{{.Data.data.postgres_password}}{{end}}
        EOF
        destination = "secrets/db-password.secret"
      }

      template {
        data = <<EOF
{{with secret "kv/data/${env}"}}{{.Data.data.jwt_secret}}{{end}}
        EOF
        destination = "secrets/jwt-secret-key.secret"
      }

      template {
        data = <<EOF
{{with secret "kv/data/${env}"}}{{.Data.data.hash_salt}}{{end}}
        EOF
        destination = "secrets/hash-salt.secret"
      }

      vault {
        policies = ["gql-server-${env}"]
      }

      resources {
        cpu    = 256
        memory = 256
      }
    }

    network {
      mode = "bridge"
      port "http" {
        to = "3000"
      }
    }

    service {
      name = "gql-server-${env}"
      port = "http"

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "gql-postgres-${env}"
              local_bind_port = 5432
            }
            upstreams {
              destination_name = "gql-redis-${env}"
              local_bind_port = 6379
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
      tags = [
        "traefik.enable=true",
        "traefik.http.routers.gql-server-${env}.rule=Host(`${host}`) && PathPrefix(`/graphql`)",
      ]

      check {
        type     = "http"
        path     = "/graphql"
        interval = "2s"
        timeout  = "2s"
      }
    }
  }

  group "postgres" {
    count = 1

    volume "gql-postgres-${env}" {
      type      = "csi"
      read_only = false
      source    = "${volume_id}"
    }

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

      volume_mount {
        volume      = "gql-postgres-${env}"
        destination = "/var/lib/postgresql/data"
        read_only   = false
      }

      template {
        data = <<EOF
{{with secret "kv/data/${env}"}}{{.Data.data.postgres_password}}{{end}}
        EOF
        destination = "secrets/db-password.secret"
      }

      vault {
        policies = ["gql-server-${env}"]
      }
    }
    
    network {
      mode = "bridge"
    }

    service {
       name = "gql-postgres-${env}"
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
       name = "gql-redis-${env}"
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
