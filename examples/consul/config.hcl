log_level = "$LOG_LEVEL"
protobuffers = ["./schemas/*.proto"]

include = ["flow.hcl"]

http {
  address = ":8080"
}

discovery "consul" {
  address = "http://localhost:8500"
}

// It's also possible to define a named discovery server client.
// By default, the provider type is taken from the block title.
// But we can use any title for the block and set provider manually.
// In this case, we have to refer to the discovery resolver by the custom name: "myAwesomeConsul".
//
// discovery "myAwesomeConsul" {
//   address = "http://localhost:8500"
//   provider = "consul"
// }

service "com.semaphore" "awesome-dogs" {
  transport = "http"
  codec     = "json"
  host      = "http://awesome-dogs"
  resolver  = "consul"

  method "List" {
    response = "com.semaphore.Dogs"
    request = "com.semaphore.Void"

    options {
      endpoint = "/"
      method = "GET"
    }
  }
}