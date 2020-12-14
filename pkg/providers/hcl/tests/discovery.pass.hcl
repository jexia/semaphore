discovery "consul" {
  address = "http://localhost:8500"
}

discovery "foobar" {
  provider = "consul"
  address = "http://localhost:8500"
}