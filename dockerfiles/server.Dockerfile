FROM --platform=${BUILDPLATFORM} golang:1.24 AS builder

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

ARG TARGETOS
ARG TARGETARCH
ARG BUILDPLATFORM

WORKDIR /app

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o cshort ./cmd/server

FROM alpine:3.21

ARG USERID=1001

RUN adduser -HD -u ${USERID} cshort

COPY ./db/migrations /db/migrations
COPY --from=builder /app/cshort .

USER ${USERID}:${USERID}

CMD [ "./cshort" ]
