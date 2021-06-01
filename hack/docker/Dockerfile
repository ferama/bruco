FROM golang:1.16 as gobuilder
WORKDIR /go/src/app
COPY . .
RUN go build -o /bruco .

# Final docker image
FROM ubuntu:latest
COPY --from=gobuilder /bruco /usr/local/bin/bruco