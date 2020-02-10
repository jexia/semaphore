flow "echo" {
    call "set" "setter.Set" {
        request {
            header {
                Cookie = "mnomnom"
            }
        }
    }

    call "get" "getter.Get" {
        request {
            header {
                Cookie = "mnomnom"
            }
        }
    }

    output {
        request {
            header {
                Cookie = "mnomnom"
            }
        }
    }
}