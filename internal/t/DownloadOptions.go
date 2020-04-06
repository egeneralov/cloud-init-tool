package t

type DownloadOptions struct {
	URL      string
	Filename string

	TimeoutDial int
	TimeoutTLS  int
	TimeoutHTTP int
}

func (DownloadOptions) New(URL string, path string) DownloadOptions {
	return DownloadOptions{
		URL:         URL,
		Filename:    path,
		TimeoutDial: 5,
		TimeoutTLS:  5,
		TimeoutHTTP: 90,
	}
}
