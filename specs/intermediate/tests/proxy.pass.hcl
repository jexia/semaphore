proxy "echo" {
    forward "uploader" "File" {

    }
}

proxy "ping" {
    forward "uploader" "File" {
        header {
            cookie = "mnomnom"
        }
    }
}