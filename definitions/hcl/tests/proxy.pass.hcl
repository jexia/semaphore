proxy "echo" {
    forward "uploader" "File" {

    }
}

proxy "ping" {
    forward "uploader" {
        header {
            cookie = "mnomnom"
        }
    }
}