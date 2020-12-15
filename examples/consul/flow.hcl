endpoint "ListAwesomeDogs" "http" {
    endpoint = "/"
    method   = "GET"
    codec    = "json"
}

flow "ListAwesomeDogs" {
    resource "list" {
        request "com.semaphore.awesome-dogs" "List" {}
    }

    output {
        payload "com.semaphore.Dogs" {
            dogs = "{{ list:dogs }}"
        }
    }
}
