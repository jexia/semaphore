flow "mock" {
    before {
        resource "check" {
            request "com.maestro" "Fetch" {
                key = "value"
            }
        }

        resources {
            sample = "key"
        }
    }
}

proxy "mock" {
    before {
        resource "check" {
            request "com.maestro" "Fetch" {
                key = "value"
            }
        }

        resources {
            sample = "key"
        }
    }

    forward "" {}
}