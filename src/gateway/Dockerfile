FROM golang:1.13-alpine as builder

WORKDIR /gateway
COPY go.mod . 
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/gateway


FROM alpine
COPY --from=builder /go/bin/gateway /go/bin/gateway
ENTRYPOINT ["/go/bin/gateway"]