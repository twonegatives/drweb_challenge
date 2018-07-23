package drweb

type FileEncoder interface {
	Encode(contents []byte) []byte
}

type Storage interface {
	Save(contents []byte, filename string) (string, error)
	Load(filename string) ([]byte, error)
	Delete(filename string) error
}

type FileSaveHooks interface {
	Before(file *File) error
	After(file *File, filename string, filepath string) error
}

type File struct {
	Body        []byte
	Encoder     FileEncoder
	Storage     Storage
	HooksOnSave FileSaveHooks
}

func (f *File) encode() string {
	hashbytes := f.Encoder.Encode(f.Body)
	return string(hashbytes[:])
}

func (f *File) Save() (string, error) {
	if err := f.HooksOnSave.Before(f); err != nil {
		return "", err
	}

	filename := f.encode()
	filepath, err := f.Storage.Save(f.Body, filename)

	if err != nil {
		return "", err
	}

	if err := f.HooksOnSave.After(f, filename, filepath); err != nil {
		return "", err
	}

	return filepath, nil
}
