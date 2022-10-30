FROM golang:1.18 as builder

WORKDIR /build

COPY . .

RUN make build

FROM alpine:3.14

WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories    && \
    apk update -vU --allow-untrusted                                                && \
    apk add --no-cache bash tzdata openssl  && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime                             && \
    echo "Asia/Shanghai" > /etc/timezone                                            && \
    apk del tzdata                                                                 && \
    rm -rf /var/cache/apk/*  && rm -rf /root/.cache && rm -rf /tmp/*

COPY --from=builder /build/release/bin/gitlab_bot /user/local/bin

ENTRYPOINT ["/user/local/bin/gitlab_bot"]