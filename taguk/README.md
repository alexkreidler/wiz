# taguk

Taguk is a simple Golang library for building and exposing generic APIs, regardless of the transport/API layer.

It follows a code-first API design model due to Golang's strong and powerful type system with reflection. 

It will support the following transport protocols:
- HTTP
- gRPC
- CapnProto

It supports the following API definition languages:
- JSON (hyper)Schema
- OpenAPI
- Protobuf schema
- CapNProto schema

**Note: Advanced options** 
Once you generate your API definitions, if you want to remove yourself from the Taguk framework, you can simply
1. generate stubs in Go
2. call into your business library code (e.g. the Taguk system was just a wrapper for your library)

Thus you can customize your app more effectively to the transport protocol you've picked.

<!-- The library works like the test -->

The core of the library is entirely implemented in the taguk.go file, but extensions like JSON schema and protobuf generation are in other packages.

In general, the main design pattern is that all resources are represented by Structs that then get added to the server.

You can define any method on the structs, and these will be called into by the RPCs/HTTP requests. However, there are a few builtin methods that have special functionality for HTTP. These are:


| CRUD Role | Method Signature                                   | HTTP Path                             |
| --------- | -------------------------------------------------- | ------------------------------------- |
| Create    | Create(item Type)                                  | `POST /resource_name/create`          |
| Read      | GetAll()                                           | `GET /resource_name`                  |
| Read      | Get(id int64)                                      | `GET /resource_name/:id`              |
| Update    | Update(id int64, newItem Type)                     | `POST /resource_name/:id`             |
| Delete    | Delete(id int64)                                   | `DELETE /resource_name/:id`           |
| N/A       | CustomIndividualAction(id int64)                   | `GET /resource_name/:id/action_name`  |
| N/A       | CustomIndividualAction(id int64, data interface{}) | `POST /resource_name/:id/action_name` |
| N/A       | CustomResourceAction()                             | `GET /resource_name/action_name`      |
| N/A       | CustomResourceAction(data interface{})             | `POST /resource_name/action_name`     |

TODO: work on the above custom actions to ensure following of HTTP/REST best practices such as idempotency, etc.
This could require each action to provide additional information about itself.

## API Definition languages

Each API definition language supports a variety of features.

All support specifying required vs optional values.
gRPC and the other RPC libs also natively support streams whereas the HTTP IDLs don't,
even though this can be added to HTTP APIs through websockets.

For optional values, we look for a Go struct tag as follows:

```go
type Resource struct {
    OptionalField string `"required:false"`
    RequiredField int `"required:true"`
    AnotherRequiredField bool
}
```

As you can see, fields without annotations are required by default.

## References vs embedded results

In HTTP APIs, when one Resource type links to another, it can either return something like so:
```json
"sub_resource": {
    "id": "37824"
}
OR
"sub_resource": {
    "field": "actual_values",
    "etc": "etc"
}
```

This depends on the design of the API, usually several factors
1. how big the referenced resource is (e.g. hundres of KBs or a few bytes)
2. if the user can handle the N+1 fetching problem or not
3. whether or not the API wants to support the user modifying this subresource separately (solely returning it embedded makes this more complicated)
4. if the subresource is referenced by multiple other resources (e.g. then how does an update to the main resource that modifies the subresource get updated across other usages of the subresource)

Genereally, whether or not a field is embedded in the API response depends on whether it is a pointer in the Go code.

If it is a pointer, it is referenced with an ID, but if it is not, it is embedded. However, this default can be changed at the field level by using Go struct tags as follows:


```go
type Resource1 struct {...}
type Resource2 struct {...}
type MainResource struct {
    EmbeddedSubresource1 Resource1
    EmbeddedSubresource2 *Resource2 `"api_embed:true"`

    ReferencedSubresource1 *Resource1
    ReferencedSubresource2 Resource2 `"api_embed:false"`
}
```
