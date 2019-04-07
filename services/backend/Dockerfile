# build stage
FROM golang:alpine as builder
RUN apk add git
WORKDIR /go/src/github.com/indiependente/barcode
ADD . .
WORKDIR /go/src/github.com/indiependente/barcode/services/backend
RUN GO111MODULE=on go mod download
RUN CGO_ENABLED=0 GO111MODULE=on go build -ldflags="-s -w" -o /service

# run stage
FROM scratch
WORKDIR /app
COPY --from=builder /service /app/
EXPOSE 8080
CMD [ "./service" ]
