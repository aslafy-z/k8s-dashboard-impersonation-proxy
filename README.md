# k8s-dashboard-impersonation-proxy

This is a tool that injects authorization and remaps impersonation headers to the Kubernetes Dashboard format.

![Kubernetes dashboard impersonation scheme](https://user-images.githubusercontent.com/8249283/54769658-03c87a00-4bd8-11e9-86ea-ea9165bb82da.png)

References:
- https://kubernetes.io/docs/reference/access-authn-authz/authentication/#user-impersonation
- https://github.com/kubernetes/dashboard/blob/master/docs/user/README.md#user-impersonation

## Usage with nginx and oauth2-proxy

oauth2-proxy reference: <https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/overview/>

oauth2-proxy with the `--set-xauthrequest` flag will set the following headers:

- `X-Auth-Request-Preferred-Username` holding the username
- `X-Auth-Request-Groups` holding the groups (comma separated)

This tool will remap these headers to the Kubernetes Dashboard impersonation headers:

- `Impersonate-User` holding the username
- `Impersonate-Group` holding the groups (one header per group)

Additionally, it will inject the `Authorization` header with the `Bearer` token sourced from a Kubernetes service account.

## Local development

```bash
$ go build
$ ./k8s-dashboard-impersonation-proxy
2023/04/12 16:31:24 Starting Server
```

```bash
$ curl http://localhost:8080
# should send request to target url
```

## Setup

Automatically built Docker image can be found at `ghcr.io/aslafy-z/k8s-dashboard-impersonation-proxy:latest`. Latest being the latest release, you can replace it with any Git tag.

TBD - Sample Kubernetes setup with nginx and oauth2-proxy.

## Contributing

Simply create an issue or a pull request.
