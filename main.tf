terraform {
  backend "pg" {
    schema_name = "d6p10tscl758n6"
  }
}

variable "heroku_username" {
  description = "Heroku user name"
}

variable "heroku_api_key" {
  description = "Heroku api key"
}

variable "app_name" {
  description = "Name of the Heroku app provisioned"
}

variable "jwt_secret" {
  description = "JWT secret"
}

provider "heroku" {
  email   = "${var.heroku_username}"
  api_key = "${var.heroku_api_key}"
}

resource "heroku_app" "default" {
  name   = "${var.app_name}"
  region = "eu"
  stack  = "container"
}

resource "heroku_addon" "default" {
  app    = "${heroku_app.default.name}"
  plan   = "mongolab:sandbox"
}

resource "heroku_config" "default" {
  sensitive_vars = {
    JWT_SECRET = "${var.jwt_secret}"
  }
}

resource "heroku_app_config_association" "default" {
  app_id = "${heroku_app.default.id}"

  sensitive_vars = "${heroku_config.default.sensitive_vars}"
}
