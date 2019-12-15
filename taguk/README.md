# taguk

> Notice: Although Taguk will remain in development and likely become the system for developing future Wiz APIs, it is too hacky as of right now, so the 0.1.0 Processor API will just be built using an API IDL.

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

The core of the library is entirely implemented in the taguk.go file, but extensions like JSON schema and protobuf generation are in other packages, along with each transport server.

Examples: If you want to see examples of how to use taguk, just look at the tests in /test and /httpjson.

## Data model

The core Taguk data model is a set of **Resources**, implemented as Structs with Receiver methods, that are added by the user to the **Server**.

Each resource has a collection of **Items**, which are instances of that Resource. **Actions** can be taken on either the resource itself (e.g. like a service) or on an individual item.

This provides for a very simple and intuitive mapping to HTTP REST APIs, and it can also be mapped well to RPC based APIs by simply exposing all actions and returning data.

**Note: IDs** in Taguk the default ID type is int64, but in the future we may change this to a string e.g. UUIDs.

You can define any method on the structs, and these will be called into by the RPCs/HTTP requests. However, there are a few builtin **Specialized methods** that have special functionality for HTTP. These are:

## Specialized Methods


| CRUD Role | Method Signature                                   | HTTP Path                             |
| --------- | -------------------------------------------------- | ------------------------------------- |
| Create    | Create(item Type)                                  | `POST /resource_name` or `POST /resource_name/create`          |
| Read      | GetAll()                                           | `GET /resource_name`                  |
| Read      | Get(id int64)                                      | `GET /resource_name/:id`              |
| Update    | Update(id int64, newItem Type)                     | `POST /resource_name/:id`             |
| Delete    | Delete(id int64)                                   | `DELETE /resource_name/:id`           |
| N/A       | CustomIndividualAction(id int64)                   | `GET /resource_name/:id/action_name`  |
| N/A       | CustomIndividualAction(id int64, data interface{}) | `POST /resource_name/:id/action_name` |
| N/A       | CustomResourceAction()                             | `GET /resource_name/action_name`      |
| N/A       | CustomResourceAction(data interface{})             | `POST /resource_name/action_name`     |

As shown above, any function that takes an int64 as the first argument is assumed to be operating on an individual Item.

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

## Meaning

You may be wondering what taguk means. Well its: Tachyon API Generator Unified Kryptonite.

Just kidding, it doesn't mean anything.