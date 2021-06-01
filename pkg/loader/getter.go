package loader

type Getter interface {
	Download(string) (string, error)
	Cleanup()
}
