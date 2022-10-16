FROM golang:1.19 AS builder

WORKDIR /go/src/project
COPY go.* ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN go test ./...
RUN GOOS=linux CGO_ENABLED=0 go build -o server ./cmd/server

FROM alpine:3.16
WORKDIR /root
COPY --from=builder /go/src/project/server .

ENTRYPOINT ["./server"]