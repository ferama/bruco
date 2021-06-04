FROM golang:1.16 as gobuilder
WORKDIR /go/src/app
COPY . .
RUN go build -o /bruco-controller ./cmd/bruco-controller

# Final docker image
FROM ubuntu:latest
COPY --from=gobuilder /bruco-controller /usr/local/bin/bruco-controller

CMD [ "bruco-controller" ]