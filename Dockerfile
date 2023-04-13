FROM golang AS builder
WORKDIR /app
COPY . .
RUN go get .
RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o k8s-dashboard-impersonation-proxy .

FROM scratch
COPY --from=builder --chmod=0755 /app/k8s-dashboard-impersonation-proxy /
EXPOSE 8080
CMD ["/k8s-dashboard-impersonation-proxy"]
