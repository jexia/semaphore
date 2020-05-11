# GraphQL

The GraphQL transport allows to expose flows as query objects.
All flows could be exposed with `graphql` endpoints and the following optional options.

```hcl
endpoint "flow" "graphql" {
    path = "user.address"
    name = "address"
    base = "mutation"
}
```
