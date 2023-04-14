FROM golang:alpine # AS builder
WORKDIR /app
COPY . .
RUN go get .
RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o k8s-dashboard-impersonation-proxy .

## TODO: Restore once debugging is done
# FROM scratch
# COPY --from=builder --chmod=0755 /app/k8s-dashboard-impersonation-proxy /
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app/k8s-dashboard-impersonation-proxy"]
EXPOSE 8080
