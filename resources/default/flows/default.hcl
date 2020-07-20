endpoint "default" "http" {
	endpoint = "/"
}

flow "default" {
	output "semaphore.Welcome" {
		message = "Welcome to Semaphore"
	}
}
