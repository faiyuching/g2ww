FROM golang:1.15.8-alpine AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache --no-progress add -U tzdata

COPY .  /go/src
WORKDIR /go/src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/g2ww *.go

FROM scratch

EXPOSE 2048

ENV TZ=Asia/Shanghai
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
COPY --from=builder /go/bin/g2ww /g2ww

ENTRYPOINT ["/g2ww"]