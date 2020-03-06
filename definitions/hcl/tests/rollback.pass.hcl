flow "echo" {
    call "get" {
        rollback "getter" "Remove" {
        }
    }
}