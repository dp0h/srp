FROM golang:latest as builder

WORKDIR /build/app

COPY . /build/

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.revision=0.1`date -u +.%Y%m%d.%H%M%S`" -o srp .

FROM scratch

COPY --from=builder /build/app/srp /app/srp

ARG SRP_BUILDARG_CERT_FILE
ARG SRP_BUILDARG_KEY_FILE
COPY ${SRP_BUILDARG_CERT_FILE} /app/
COPY ${SRP_BUILDARG_KEY_FILE} /app/

ARG SRP_BUILDARG_CONFIG_FILE
COPY ${SRP_BUILDARG_CONFIG_FILE} /app/

ENTRYPOINT ["/app/srp"]
