FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache --virtual .build-deps \
        ca-certificates \
        tzdata \
        gcc \
        g++ &&  \
        go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download

COPY . .

RUN go build -ldflags "-s -w" -o treehole

FROM alpine

WORKDIR /app

COPY --from=builder /app/treehole /app/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV TZ=Asia/Shanghai

EXPOSE 8000

ENTRYPOINT ["./treehole"]