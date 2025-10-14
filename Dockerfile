FROM oven/bun:1.3 AS frontend
WORKDIR /src
COPY ./frontend ./
RUN bun install --frozen-lockfile

ENV NODE_ENV=production
# RUN bun run test # does not work currently in ci
RUN bun run build

FROM golang:1.25-trixie AS backend
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test ./...

ARG CGO_ENABLED=0
ENV CGO_ENABLED=${CGO_ENABLED}
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/server ./main.go


FROM alpine:3.22 AS application
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

VOLUME /var/run/docker.sock

RUN apk add --no-cache docker-cli

COPY --from=frontend /src/dist ./frontend/dist
COPY --from=backend /app/server ./server
COPY ./custom_images ./custom_images

RUN mkdir pb_data

EXPOSE 8080

CMD ["./server"]
