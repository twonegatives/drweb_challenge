package drweb

type FileEncoder interface {
	Encode(contents []byte) []byte
}

type Storage interface {
	Save(f *File) (string, error)
	Load(filename string) (*File, error)
	Delete(filename string) error
}

type FileSaveHooks interface {
	Before(file *File, args ...interface{}) error
	After(file *File, args ...interface{}) error
}

type File struct {
	Body        []byte
	Storage     Storage
	HooksOnSave FileSaveHooks
	Encoder     FileEncoder
}

func (f *File) GetFilename() string {
	return string(f.Encoder.Encode(f.Body)[:])
}

func (f *File) Save() (string, error) {
	if err := f.HooksOnSave.Before(f); err != nil {
		return "", err
	}

	filepath, err := f.Storage.Save(f)

	if err != nil {
		return "", err
	}

	if err := f.HooksOnSave.After(f); err != nil {
		return "", err
	}

	return filepath, nil
}
