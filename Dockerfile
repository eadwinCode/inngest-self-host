FROM --platform=$BUILDPLATFORM golang:1.25.9-alpine AS build
RUN apk add build-base
WORKDIR /app
COPY vendor vendor
COPY . .
ARG TARGETARCH
ARG TARGETOS
ARG VERSION=dev
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOFLAGS=-mod=vendor GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -ldflags="-s -w -X github.com/inngest/inngest/pkg/inngest/version.Version=${VERSION}" \
    -o /go/bin/inngest ./cmd/

FROM alpine:3.16 AS inngest
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates
COPY --from=build /go/bin/inngest /bin/inngest
CMD ["inngest"]
