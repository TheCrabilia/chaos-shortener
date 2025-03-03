variable "REGISTRY" {
  default = "ghcr.io/thecrabilia"
}

group "default" {
  targets = ["server", "client"]
}

target "_common" {
  platforms = ["linux/amd64", "linux/arm64"]
}

target "server" {
  inherits = ["_common"]
  context = "."
  dockerfile = "./dockerfiles/server.Dockerfile"
  tags = ["${REGISTRY}/chaos-shortener:latest"]
}

target "client" {
  inherits = ["_common"]
  context = "."
  dockerfile = "./dockerfiles/client.Dockerfile"
  tags = ["${REGISTRY}/cs-client:latest"]
}
