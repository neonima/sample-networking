FROM golang:1.20-alpine as builder
ARG APPNAME=ALO
RUN apk add make build-base
RUN mkdir -p /app
WORKDIR /app
ADD . .
RUN make tools
RUN make assets
RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build  -ldflags="-w -s -X main.Appname=$APPNAME" -o app cmd/server/main.go


FROM alpine:3.9 as ca
RUN adduser -D -g '' appuser
RUN apk add -U --no-cache ca-certificates

FROM gcr.io/distroless/static-debian11
COPY --from=ca /etc/passwd /etc/passwd
COPY --from=ca /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/app .


USER appuser
ENTRYPOINT   [ "./app" ]