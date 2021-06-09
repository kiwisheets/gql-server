resource "cloudflare_access_application" "gql_server" {
  zone_id          = var.cloudflare_zone_id
  name             = "GQL Server Dev"
  domain           = "api-dev-gql-server.kiwisheets.com"
  session_duration = "24h"
}

resource "cloudflare_access_policy" "gql_server_policy" {
  application_id = cloudflare_access_application.gql_server.id
  zone_id        = var.cloudflare_zone_id
  name           = "Service Auth"
  decision       = "non_identity"
  precedence     = 1

  include {
    any_valid_service_token = true
  }
}
