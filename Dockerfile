# Should be started with:
# docker run -ti -v $(pwd):/work fclairamb/ci-info

# Preparing the build environment
FROM golang:1.20.0-alpine3.16 AS builder
# RUN apk add --update --no-cache bash ca-certificates curl git
RUN mkdir -p /build
WORKDIR /build

# Building
COPY . .
RUN CGO_ENABLED=0 go build -mod=readonly -ldflags='-w -s' -v -o ci-info 

# Preparing the final image
# FROM alpine:3.16.2
FROM scratch
WORKDIR /work
COPY --from=builder /build/ci-info /bin/ci-info
ENTRYPOINT [ "/bin/ci-info" ]
