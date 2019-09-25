FROM golang:1.11-alpine as builder
WORKDIR /todo-api

COPY go.mod . 
COPY go.sum .

RUN apk add --no-cache ca-certificates git

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 go install -installsuffix "static" ./cmd/api/...

FROM alpine
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /go/bin /bin
ENTRYPOINT ["/bin/api"]