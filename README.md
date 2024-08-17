##### implementing distributed services with **Golang**

###### features: 
- [x] commit log
- [ ] networking with gRPC



###### implementation:
- commit log:
    - using store and index files approach for each log segment
    - using [go-mmap](https://pkg.go.dev/github.com/go-mmap/mmap) library to memory map index file for performance issues.
    - test for each segment and it's store and index files
- gRPC Services:
    - using bidirectional streaming APIs on the client and server side to stream the content between them.
    - using [status](https://godoc.org/google.golang.org/grpc/status), [codes](https://godoc.org/google.golang.org/grpc/codes) and [errdetails](https://godoc.org/google.golang.org/genproto/googleapis/rpc/errdetials) packages to customize error messages between client and server.
    - Dependency Inversion using Interfaces. (DIP principle). --> (wanna know more?)[https://medium.com/@sumit-s/the-dependency-inversion-principle-dip-in-golang-fb0bdc503972]
- Security:
    - security in distributed services can be broken down into three steps:
        - encrypt data in-flight to protect against man-in-the-middle attacks
        - authenticate to identify clients
        - authorize to determine the permission of the identified clients
    - the adoptiveness of mutual authentication in distributed services. interested? (learn how cloudflare adopt it)[https://blog.cloudflare.com/how-to-build-your-own-public-key-infrastructure].
    - building access control list-based authorization & differentiate between authentication and authorization in case of varying levels of access and permissions.
    - using cloudflare's open source Certificate Authority (CA) called (CSFFL)[https://blog.cloudflare.com/introducing-cfssl] for signing, verifying and bundling TLS certificates.
    - 