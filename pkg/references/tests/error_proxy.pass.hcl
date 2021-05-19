proxy "echo" {
    error {
        payload "com.error" {}
    }

    forward "" {}
}
