// configure the max size of a segment's store and index
// centralize the log's configuration with config struct
package log

type Config struct {
	Segment struct {
		MaxStoreBytes uint64
		MaxIndexBytes uint64
		InitialOffset uint64
	}
}
