FROM golang:latest AS builder
ENV GOPROXY https://goproxy.cn,direct
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# RUN apk update --no-cache && apk add --no-cache tzdata
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/svcPodExporter

FROM alpine:latest@sha256:25fad2a32ad1f6f510e528448ae1ec69a28ef81916a004d3629874104f8a7f70
# FROM alpine:latest
LABEL date="2023-07-05"
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# RUN apk update --no-cache && apk add --no-cache tzdata
# RUN apk add curl
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/svcPodExporter /app/svcPodExporter
# COPY --from=builder /build/prometheusToZabbix.yml /app/prometheusToZabbix.yml
ENTRYPOINT ["/app/svcPodExporter"]
EXPOSE 8080
CMD ["--serviceaccount",  "-p", "8080"]