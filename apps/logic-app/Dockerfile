#FROM golang:alpine AS development
#WORKDIR $GOPATH/src/
#COPY . .
#RUN go build -o app

FROM alpine:latest
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --update curl && rm -rf /var/cache/apk/*
WORKDIR /app/
#COPY --from=development /go/src/app .
COPY ./app .
EXPOSE 8080
ENTRYPOINT ["./app"]