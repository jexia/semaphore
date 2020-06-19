endpoint "ComplexDataStructure" "http" {
	endpoint = "/"
	method = "POST"
	codec = "json"
}

flow "ComplexDataStructure" {
	input "proto.Checkout" {}

	output "proto.Checkout" {
		repeated "items" "input:items" {
            id = "{{ input:items.id }}"
            name = "{{ input:items.name }}"

            repeated "labels" "input:items.labels" {}
        }

        message "shipping" {
            time = "{{ input:shipping.time }}"

            message "address" {
                street = "{{ input:shipping.address.street }}"
                city = "{{ input:shipping.address.city }}"
            }
        }
	}
}
