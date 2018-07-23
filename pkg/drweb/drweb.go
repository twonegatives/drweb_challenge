package drweb

type FileEncoder interface {
	Encode(contents []byte) []byte
}

type Storage interface {
	Save(contents []byte, filename string) (string, error)
	Load(filename string) ([]byte, error)
	Delete(filename string) error
}
