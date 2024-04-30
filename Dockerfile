FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build
EXPOSE 50051
ADD go.mod .
ADD go.sum .

COPY ./pkg/services-proto ./pkg/services-proto
COPY ./pkg/support-go ./pkg/support-go

RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/http-server http-server.go

# nomal
#FROM alpine
# smaller
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app

COPY ./manifest/config /app/manifest/config

COPY --from=builder /app/http-server /app/http-server

CMD ["./http-server"]
