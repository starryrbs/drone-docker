FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go mod download

COPY . .
RUN go build ./cmd/drone-docker
RUN chmod +x drone-docker

FROM docker:19.03.12

WORKDIR /

COPY --from=build /app/drone-docker /drone-docker

ENTRYPOINT ["/drone-docker"]