FROM golang

RUN export GOPRIVATE=github.com/houko/wechatgpt && \
    export GOPROXY=https://goproxy.cn

RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list

RUN apt-get update && apt-get install -y default-mysql-server && rm -rf /var/lib/apt/lists/*


COPY . /root/build
COPY wget -O /root/build/go1.22.4.linux-amd64.tar.gz https://golang.org/dl/go1.22.4.linux-amd64.tar.gz

RUN rm -rf /usr/local/go && tar xvf /root/build/go1.22.4.linux-amd64.tar.gz -C /usr/local > /dev/null

WORKDIR /root/build

RUN echo "nameserver 8.8.8.8" > /etc/resolv.conf

RUN export GOPROXY=https://goproxy.cn &&  go build -o server main.go


CMD ["/root/build/init_sqlenv.sh"]
