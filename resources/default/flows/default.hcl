endpoint "default" "http" {
	endpoint = "/"
}

flow "default" {
	output "maestro.Welcome" {
		message = "Welcome to Maestro"
	}
}
