FROM golang:1.16
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go env -w GOPROXY=https://goproxy.cn
RUN go mod download
CMD ["bash", "/app/launch.sh"]