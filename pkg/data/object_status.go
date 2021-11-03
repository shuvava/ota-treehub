package data

// ObjectStatus is enumeration of possible object states
type ObjectStatus int

const (
	// UPLOADED means data.Object was successfully uploaded to server
	UPLOADED ObjectStatus = iota
	// CLIENT_UPLOADING means data.Object was created, but content not fully uploaded to the server
	CLIENT_UPLOADING
	SERVER_UPLOADING
)
