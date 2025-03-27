FROM golang:1.24.1-alpine3.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY cmd/server/ ./cmd/server/
COPY internal/ ./internal/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -tags=linux -a -installsuffix cgo -o monintor-server ./cmd/server/

FROM alpine:3.21.3
ENV PORT=8080
#RUN apk --no-cache add ca-certificates
WORKDIR /root/
#EXPOSE 8080/tcp
COPY --from=builder /app/monintor-server ./
CMD ["/root/monintor-server"]
#RUN ./monintor-server