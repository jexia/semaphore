flow "echo" {
  resource "get" {
    rollback "getter" "Remove" {}
  }
}
