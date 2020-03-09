# HTTP

## Service

A HTTP service accepts 

```hcl
service "placeholder" "http" "json" {
	host = "https://jsonplaceholder.typicode.com"
	schema = "proto.TODO"
}

flow "example" {
  call "simple" {
    request "placeholder" "User" {
      options {
        http.endpoint = "/users/1"
        http.method = "GET"
      }
    }
  }
}
```