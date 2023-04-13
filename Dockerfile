FROM golang:1.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o oauth2-proxy-k8s-impersonation .

FROM scratch
COPY --from=builder /app/oauth2-proxy-k8s-impersonation /oauth2-proxy-k8s-impersonation
EXPOSE 8080
CMD ["./oauth2-proxy-k8s-impersonation"]
