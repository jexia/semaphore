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