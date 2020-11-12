variable "datacenter" {
  type = string
}

variable "hcloud_token" {
  type        = string
  description = "Hetzner Cloud API Token"
}

variable "image_tag" {
  type        = string
  description = "image version"
}

variable "instance_count" {
  type        = number
  description = "number of server instances to launch"
}

variable "allowed_origins" {
  type        = string
  description = "the allowed origins for CORS"
}

variable "host" {
  type        = string
  description = "API host"
}
