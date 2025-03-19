variable "REGISTRY" {
  default = "ghcr.io/thecrabilia"
}

group "default" {
  targets = ["server", "client"]
}

target "_common" {
  platforms = ["linux/amd64", "linux/arm64"]
  labels = {
    "org.opencontainers.image.source" = "https://github.com/TheCrabilia/chaos-shortener"
  }
}

target "server" {
  inherits = ["_common"]
  context = "."
  dockerfile = "./dockerfiles/server.Dockerfile"
  tags = ["${REGISTRY}/chaos-shortener/server:latest"]
}

target "client" {
  inherits = ["_common"]
  context = "."
  dockerfile = "./dockerfiles/client.Dockerfile"
  tags = ["${REGISTRY}/chaos-shortener/client:latest"]
}
