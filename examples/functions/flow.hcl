endpoint "string" "http" {
    endpoint = "/string"
    method   = "POST"
    codec    = "json"
}

flow "string" {
    input {
        payload "com.semaphore.GreetRequest" {}
    }

    output {
        payload "com.semaphore.GenericResponse" {
            message = "{{ sprintf('Hey %s! What is that %s?', input:name, input:subject) }}"
        }
    }
}

endpoint "numeric" "http" {
    endpoint = "/numeric"
    method   = "POST"
    codec    = "json"
}

flow "numeric" {
    input {
        payload "com.semaphore.AgeRequest" {}
    }

    output {
        payload "com.semaphore.GenericResponse" {
            message = "{{ sprintf('Hey %s! I know you are %d years old!', input:name, input:age) }}"
        }
    }
}

endpoint "json" "http" {
    endpoint = "/json"
    method   = "POST"
    codec    = "json"
}

flow "json" {
    input {
        payload "com.semaphore.MsgRequest" {}
    }

    output {
        payload "com.semaphore.GenericResponse" {
            message = "{{ sprintf('Hey %s! We have got your personal info in JSON: %json!', input:name, input:info) }}"
        }
    }
}
