# Dynamic Client IP matcher for Caddy

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/Naereen/StrapDown.js/graphs/commit-activity) [![Go Reference](https://pkg.go.dev/badge/github.com/tuzzmaniandevil/caddy-dynamic-clientip.svg)](https://pkg.go.dev/github.com/tuzzmaniandevil/caddy-dynamic-clientip) [![GoReportCard example](https://goreportcard.com/badge/github.com/tuzzmaniandevil/caddy-dynamic-clientip?_t=3145938)](https://goreportcard.com/report/github.com/tuzzmaniandevil/caddy-dynamic-clientip) ![GitHub](https://img.shields.io/github/license/tuzzmaniandevil/caddy-dynamic-clientip?_t=3145938)

The `dynamic_client_ip` module is a clone of the `client_ip` matcher with one key difference: instead of providing IP ranges upfront, you specify an `IPRangeSource`. This allows IP ranges to be dynamically loaded per request.

## Installation

Build Caddy using [xcaddy](https://github.com/caddyserver/xcaddy):

```shell
xcaddy build --with github.com/tuzzmaniandevil/caddy-dynamic-clientip
```

## Usage

```caddyfile
:8880 {
	@denied dynamic_client_ip my_dynamic_provider
	abort @denied

    reverse_proxy localhost:8080
}
```

Example using the built-in static provider (But why though?)
```caddyfile
:8880 {
	@denied dynamic_client_ip static 12.34.56.0/24 1200:ab00::/32
	abort @denied

    reverse_proxy localhost:8080
}
```

## Development

Before diving into development, make sure to follow the [Extending Caddy](https://caddyserver.com/docs/extending-caddy#extending-caddy) guide. This ensures you're familiar with the Caddy development process and that your environment is set up correctly.

To run Caddy with this module:
```shell
xcaddy run
```

## License

The project is licensed under the [Apache License](LICENSE).