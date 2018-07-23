package drweb

type FileEncoder interface {
	Encode(contents []byte) []byte
}

type FileUploader interface {
	Upload(contents []byte, filename string) (string, error)
}
