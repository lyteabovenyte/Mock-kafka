- [implementing distributed services with **Golang**](#implementing-distributed-services-with-golang)
    - [set of implemented features:](#set-of-implemented-features)
    - [implementation:](#implementation)
    - [UpComming features -updating...](#upcomming-features--updating)
    - [Notes](#notes)

### implementing distributed services with **Golang**

##### set of implemented features: 
- [x] commit log
- [x] networking with gRPC
- [x] encrypting connection, mutual TLS authentication, ACL based authorization using [Casbin](https://github.com/casbin/casbin) and peer-to-peer grpc connection
- [x] Observability using [zap](github.com/grpc-ecosystem/go-grpc-middleware/logging/zap), [ctxtags](github.com/grpc-ecosystem/go-grpc-middleware/tags) and [OpenCensus](go.opencensus.io) for tracing. all in gRPC interceptors ⚭
- [x] Server-to-Server Service Discovery using [Serf](https://www.serf.io)
- [x] Coordination and Consesusn with **Raft** algorithm and bind the cluster with **Serf** for Discovery Integration between servers
- [x] Client-side LoadBalancing with gRPC and End-to-End Discovery Integration 


##### implementation:
- commit log:
    - *hash-table* index approach for in-memory data structure using *write-ahead-log* and *LSM Tree engine* by fragmenting index, store and segments
    - implementing segment and stores in binary format that fits best for logs. it encodes the length of a string in bytes, followed by the raw strings (page 74 of designing data-intensive application for more info)
    - using store and index files approach for each log segment by in-memory data-structure for faster seeks. (already thinking about merging the old segments)
    - using [go-mmap](https://pkg.go.dev/github.com/go-mmap/mmap) library to memory map index file for performance issues.
    - tests for each segment and it's store and index files
- gRPC Services: *v2.0.0*
    - using bidirectional streaming APIs on the client and server side to stream the content between them.
    - using [status](https://godoc.org/google.golang.org/grpc/status), [codes](https://godoc.org/google.golang.org/grpc/codes) and [errdetails](https://godoc.org/google.golang.org/genproto/googleapis/rpc/errdetails) packages to customize error messages between client and server.
    - Dependency Inversion using Interfaces. (DIP principle). --> [wanna know more?](https://medium.com/@sumit-s/the-dependency-inversion-principle-dip-in-golang-fb0bdc503972)
- Security: *v3.0.0*
    - my approach to secure the system is based on three approach:
        - encrypt data in-flight to protect against man-in-the-middle attacks
        - authenticate to identify clients
        - authorize to determine the permission of the identified clients
    - base on the adoptiveness of mutual authentication in distributed services, I'm gonna try my best to adopt it too 🤠. interested? [learn how cloudflare adopt it](https://blog.cloudflare.com/how-to-build-your-own-public-key-infrastructure).
    - building access control list-based authorization & differentiate between authentication and authorization in case of varying levels of access and permissions.
    - using cloudflare's open source Certificate Authority (CA) called [CSFFL](https://blog.cloudflare.com/introducing-cfssl) for signing, verifying and bundling TLS certificates.
    - **v3.0.2**(Authentication) has compeleted *mutual communication* between client and server + containing tests.
    - **Authorization**:
            - using *Casbin*: [Casbin](https://github.com/casbin/casbin) supports enforcing authorization based on various [control models](https://github.com/casbin/casbin#supported-models)—including ACLs. Plus Casbin extendable and i's exciting to explore it.
    - *v4.0.0* --> encrypting connection, mutual TLS authentication, ACL based authorization using **casbin**
- Observability: *v4.0.0*
    - using [zap](github.com/grpc-ecosystem/go-grpc-middleware/logging/zap) for structured logs
    - using [request context tags](github.com/grpc-ecosystem/go-grpc-middleware/tags) to set value for request tags in context.
    it'll add a Tag object to the context that can be used by other middleware to add context about a request.
    - using [OpenCensus](go.opencensus.io) for tracing
- Service-to-Service Discovery: 
    - implementing *Membership* using [Serf](https://www.serf.io) on each service instance to discover other service instances.
    - implementing *Replication* to duplicate each server's data - already thinking about Consensus -_-.
    - after implementing our *replicator*, *membership*, *log* and *server* components, we'll implement and import an **Agent** type to run and sync these components on each instance. just like [Hachicorp Consul](https://github.com/hashicorp/consul).
    - updated on *v6.0.0*
- Coordinate the Service using Raft Consensus Algorithm:
    - using my own log library as Raft's log store and satisfy the LogStore interface that Raft requires.
    - using [Bolt](https://github.com/boltdb/bolt) which is an embedded and persisted key-value database for Go, as my stable store in Raft to store server's *current Term* and important metadata like the candidates the server voted for.
    - implemented Raft snapshots to recover and restore data efficiently, when necessary, like if our server’s EC2 instance failed and an autoscaling group(terraform terminology 🥸) brought up another instance for the Raft server. Rather than streaming all the data from the Raft leader, the new server would restore from the snapshot (which you could store in S3 or a similar storage service) and then get the latest changes from the leader.
    - again, using my own Log library as Raft's finite-state-machine(*FSM*), to replicate the logs across server's in the cluster.
    - *Discovery integration* and binding *Serf* and *Raft* to implement service discovery on Raft cluster by implementing *Join* and *Leave* methods to satisfy Serf's interface hence having a Membership in the cluster to be discovered. *A **Membership** service determines which nodes are currently active and live members of the cluster*
    - **Multiplexing on our Service to run multiple services on one port**
        - we identify the Raft connections from gRPC connections by making the Raft connection write a byte to identify them by, and multiplexing connection to different listeners to handle, and configured our agents to both manage Raft cluster connections and gRPC connection on our servers on a single port
- Client-Side LoadBalancing: *v7.0.0*
    - three major features our client needs at this point:
        - discover servers in the cluster
        - direct append calls to leaders and consume calls to followers
        - balance consume calls across followers
    - we will start by writing our own *resolver* and *picker* usign gRPC builder pattern. gRPC separates the server discovery, loadbalancing and client request and response handling. in our gRPC:
      - [resolver](https://google.golang.org/grpc/resolver)s discovers the servers and whether the server is the leader or not
      - *picker*s manage directing produce calls to the leader and balancing consume calls across the followers
    - Implementing End-to-End discovery and Balancing(client-side) with our agent.





##### UpComming features -updating...
- [ ] orchestration and deployment with [kubernetes](https://kuberenetes.io) + configuring with [Helm](https://helm.sh) and tune k8s controllers to handle our cluster as we desire
- [ ] provisioning resources on AWS by Infrastructure as Code principles using [Terraform](https://www.terraform.io)
- [ ] CI/CD using [Jenkins](https://www.jenkins.io) pipeline cluster-wide + github webhooks to automate deployment
- [ ] (final GOAL👾) machine learning models as the core functionality and agent to embrace and wield every aspect of the project

##### Notes
- [ ] merging segments together for key-value stores to keep the latest value for each key and truncate out-dated segments( compaction and merging, page 73 of designing data-intensive applications)
- [ ] as our logs are fixed-length size of records, we could use *binary-search* to find the right offset on the disk. in this approach we don't need any index file to store our log's offsets and gain performance