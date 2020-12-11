flow "echo" {
  resource "set" {
    request "setter" "Set" {
      header {
        Cookie = "mnomnom"
      }
    }
  }

  resource "get" {
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
