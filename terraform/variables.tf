variable "image_tag" {
  type        = string
  description = "image version"
}

variable "environment" {
  type        = string
  description = "prod or dev environment"
}

variable "instance_count" {
  type        = number
  description = "numbe of server instances to launch"
}
