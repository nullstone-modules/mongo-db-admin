variable "name" {
  description = "The name of the lambda function and role"
  type        = string
}

variable "tags" {
  description = "A map of tags that are applied to AWS resources"
  type        = map(string)
}

variable "host" {
  description = "The database cluster host to connect"
  type        = string
}

variable "port" {
  description = "The database cluster port to connect"
  type        = string
  default     = "27017"
}

variable "database" {
  description = "The initial database to connect. By default, uses 'admin'"
  type        = string
  default     = "admin"
}

variable "username" {
  description = "Mongo username"
  type        = string
}

variable "password" {
  description = "Mongo password"
  type        = string
}

variable "network" {
  description = "Network configuration"
  type = object({
    vpc_id : string
    security_group_ids : list(string)
    subnet_ids = list(string)
  })
  default = {
    vpc_id             = ""
    security_group_ids = []
    subnet_ids         = []
  }
}
