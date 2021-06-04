FROM golang:1.16 as gobuilder
WORKDIR /go/src/app
COPY . .
RUN go build -o /bruco ./cmd/bruco

# Final docker image
FROM python:latest
COPY --from=gobuilder /bruco /usr/local/bin/bruco