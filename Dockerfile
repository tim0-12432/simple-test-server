FROM golang:1.25-trixie AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG CGO_ENABLED=0
ENV CGO_ENABLED=${CGO_ENABLED}
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/server ./main.go


FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/server ./server

RUN mkdir pb_data \
    && chown 65532:65532 pb_data

EXPOSE 8080

USER 65532:65532

CMD ["./server"]
