package blobs

// Type all supported db types
type Type string

var (
	// LocalFs defines local filesystem as blob storage
	LocalFs Type = "localfs"
)
