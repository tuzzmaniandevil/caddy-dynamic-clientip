# Dynamic Client IP matcher for Caddy

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/Naereen/StrapDown.js/graphs/commit-activity) [![GoReportCard example](https://goreportcard.com/badge/github.com/tuzzmaniandevil/caddy-dynamic-clientip)](https://goreportcard.com/report/github.com/tuzzmaniandevil/caddy-dynamic-clientip) ![GitHub](https://img.shields.io/github/license/tuzzmaniandevil/caddy-dynamic-clientip?_t=3145938)

`dynamic_client_ip` is a `client_ip` matcher clone that with one difference, instead of providing the IP ranges up front, You specify an `IPRangeSource` to allow the IP ranges to be dynamicly loaded per request.