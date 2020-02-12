flow "echo" {
    call "get" "getter.Get" {
        rollback "getter.Remove" {
        }
    }
}