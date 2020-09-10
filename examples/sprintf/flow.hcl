endpoint "wtf" "http" {
  endpoint = "/wtf"
  method   = "POST"
  codec    = "json"
}

flow "wtf" {
  input "com.semaphore.WTFRequest" {}

  output "com.semaphore.GenericResponse" {
    message = "{{ sprintf('Hey %s! What is that %s?', input:name, input:subject) }}"
  }
}

endpoint "age" "http" {
  endpoint = "/age"
  method   = "POST"
  codec    = "json"
}

flow "age" {
  input "com.semaphore.AgeRequest" {}

  output "com.semaphore.GenericResponse" {
    message = "{{ sprintf('Hey %s! I know you are %d years old!', input:name, input:age) }}"
  }
}
