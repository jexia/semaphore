flow "echo" {
    call "set" {
        request "setter" "Set" {
            header {
                Cookie = "mnomnom"
            }
        }
    }

    call "get" {
        request "getter" "Get" {
            header {
                Cookie = "mnomnom"
            }
        }
    }

    output {
        header {
            Cookie = "mnomnom"
        }
    }
}