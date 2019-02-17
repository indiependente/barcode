FROM golang:latest as builder
ADD . /usr/src/app
WORKDIR /usr/src/app
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o service .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /usr/src/app/service .
EXPOSE 8080
CMD [ "./service" ]
