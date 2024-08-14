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
    - working ...