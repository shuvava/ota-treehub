package data

// ObjectStatus is enumeration of possible object states
type ObjectStatus int

const (
	// Uploaded means data.Object was successfully uploaded to server
	Uploaded ObjectStatus = iota
	// ServerUploading means data.Object was created, but content not fully uploaded to the server
	ServerUploading
)
