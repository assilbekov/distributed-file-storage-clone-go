package main

type FileServerOpts struct {
	ListenAddr  string
	StorageRoot string
}

type FileServer struct {
	FileServerOpts
}
