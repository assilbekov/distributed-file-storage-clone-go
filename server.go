package main

type FileServerOpts struct {
	ListenAddr  string
	StorageRoot string
}

type FileServer struct {
	FileServerOpts

	store *Store
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(StoreOpts{Root: opts.StorageRoot}),
	}
}
