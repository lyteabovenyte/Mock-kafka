##### implementing distributed services with **Golang**

###### features: 
- [x] commit log
- [x] networking with gRPC
- [x] encrypting connection, mutual TLS authentication, ACL based authorization using [Casbin](https://github.com/casbin/casbin) and peer-to-peer grpc connection
- [x] Observability using [zap](github.com/grpc-ecosystem/go-grpc-middleware/logging/zap), [ctxtags](github.com/grpc-ecosystem/go-grpc-middleware/tags) and [OpenCensus](go.opencensus.io) for tracing. all in gRPC interceptors ⚭
- [ ] Server-to-Server Service Discovery


###### implementation:
- commit log:
    - using store and index files approach for each log segment
    - using [go-mmap](https://pkg.go.dev/github.com/go-mmap/mmap) library to memory map index file for performance issues.
    - tests for each segment and it's store and index files
- gRPC Services: *v2.0.0*
    - using bidirectional streaming APIs on the client and server side to stream the content between them.
    - using [status](https://godoc.org/google.golang.org/grpc/status), [codes](https://godoc.org/google.golang.org/grpc/codes) and [errdetails](https://godoc.org/google.golang.org/genproto/googleapis/rpc/errdetials) packages to customize error messages between client and server.
    - Dependency Inversion using Interfaces. (DIP principle). --> [wanna know more?](https://medium.com/@sumit-s/the-dependency-inversion-principle-dip-in-golang-fb0bdc503972)
- Security: *v3.0.0*
    - security in distributed services can be broken down into three steps:
        - encrypt data in-flight to protect against man-in-the-middle attacks
        - authenticate to identify clients
        - authorize to determine the permission of the identified clients
    - the adoptiveness of mutual authentication in distributed services. interested? [learn how cloudflare adopt it](https://blog.cloudflare.com/how-to-build-your-own-public-key-infrastructure).
    - building access control list-based authorization & differentiate between authentication and authorization in case of varying levels of access and permissions.
    - using cloudflare's open source Certificate Authority (CA) called [CSFFL](https://blog.cloudflare.com/introducing-cfssl) for signing, verifying and bundling TLS certificates.
    - **v3.0.2**(Authentication) has compeleted *mutual communication* between client and server + containing tests.
    - **Authorization**:
            - using *Casbin*: [Casbin](https://github.com/casbin/casbin) supports enforcing authorization based on various [control models](https://github.com/casbin/casbin#supported-models)—including ACLs. Plus Casbin is well adopted, tested, and extendable.
    - *v4.0.0* --> encrypting connection, mutual TLS authentication, ACL based authorization using **casbin**
- Observability: *v4.0.0*
    - using [zap](github.com/grpc-ecosystem/go-grpc-middleware/logging/zap) for structured logs
    - using [request context tags](github.com/grpc-ecosystem/go-grpc-middleware/tags) to set value for request tags in context.
    it'll add a Tag object to the context that can be used by other middleware to add context about a request.
    - using [OpenCensus](go.opencensus.io) for tracing
- Service-to-Service Discovery: 
    - implementing *Membership* using [Serf](https://www.serf.io) on each service instance to discover other service instances.
    - implementing *Replication* to duplicate each server's data
    - after implementing our *replicator*, *membership*, *log* and *server* components, we'll implement and import an **Agent** type to run and sync these components on each instance. just like [Hachicorp Consul](https://github.com/hashicorp/consul).
    - updated on *v6.0.0*