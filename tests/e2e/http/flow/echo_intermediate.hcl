
services {
    select "semaphore.typetest.*" {
        host = "http://127.0.0.1:8081/"
    }
}

endpoint "typetest" "http" {
    endpoint = "/json"
    method   = "POST"
    codec    = "json"
}

flow "typetest" {
    input {
        payload "semaphore.typetest.Request" {}
    }

    resource "echo" {
        request "semaphore.typetest.External" "Post" {
            enum    = "{{ input:data.enum }}"
            string  = "{{ input:data.string }}"
            integer = "{{ input:data.integer }}"
            double  = "{{ input:data.double }}"
            numbers = "{{ input:data.numbers }}"
        }
    }

    output {
        payload "semaphore.typetest.Data" {
            enum    = "{{ echo:enum }}"
            string  = "{{ echo:string }}"
            integer = "{{ echo:integer }}"
            double  = "{{ echo:double }}"
            numbers = "{{ echo:numbers }}"
        }
    }
}

