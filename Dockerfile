FROM golang:1.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o k8s-dashboard-impersonation-proxy .

FROM scratch
COPY --from=builder /app/k8s-dashboard-impersonation-proxy /k8s-dashboard-impersonation-proxy
EXPOSE 8080
CMD ["./k8s-dashboard-impersonation-proxy"]
