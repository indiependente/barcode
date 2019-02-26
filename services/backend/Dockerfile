# build stage
FROM golang:latest as builder
WORKDIR /go/src/github.com/indiependente/barcode
ADD . .
WORKDIR /go/src/github.com/indiependente/barcode/services/backend
RUN CGO_ENABLED=0 GO111MODULE=on go build -o service

# run stage
FROM scratch
WORKDIR /app
COPY --from=builder /usr/src/app/service .
EXPOSE 8080
CMD [ "./service" ]
