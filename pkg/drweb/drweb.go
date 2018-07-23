package drweb

type FileEncoder interface {
	Encode(contents []byte) []byte
}

type Storage interface {
	Save(contents []byte, filename string) (string, error)
	Load(filename string) ([]byte, error)
	Delete(filename string) error
}

type File struct {
	Body    []byte
	Encoder FileEncoder
	Storage Storage
}

func (f *File) encode() string {
	hashbytes := f.Encoder.Encode(f.Body)
	return string(hashbytes[:])
}

func (f *File) Save() (string, error) {
	filename := f.encode()
	return f.Storage.Save(f.Body, filename)
}
