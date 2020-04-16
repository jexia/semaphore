# HTTP

Provides HTTP client and server implementations.

```hcl
endpoint "mock" "http" {
	endpoint = "/"
	method = "GET"
	codec = "json"
	read_timeout = "5s"
	write_timeout = "5s"
}
```

Services could be defined inside the HCL definitions.

```hcl
service "mock" "http" {
	host = "https://service.prod.svc.cluster.local"

	options {
		flush_interval = "1s"
		timeout = "60s"
		keep_alive = "60s"
		max_idle_conns = "100"
	}
}
```

Or in schema definitions such as proto.

```proto
rpc Mock(Empty) returns (Empty) {
	option (maestro.http) = {
		endpoint: "/endpoint"
		method: "GET"
	};
};
```

Object properties available inside the request object could be referenced inside a endpoint.

```proto
rpc Mock(Empty) returns (Empty) {
	option (maestro.http) = {
		endpoint: "/endpoint/:id"
		method: "GET"
	};
};
```