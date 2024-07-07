FROM golang:1.22 AS builder
WORKDIR /workspace
ENV GO111MODULE=on CGO_ENABLED=0
COPY . .
RUN go test ./... && go build -o /bin/server ./cmd/main.go


FROM alpine:latest AS release
ENV ENV=stg
COPY --from=builder /bin/server /bin/server
COPY ./config/ /config/
ENTRYPOINT ["./bin/server"]
