package drweb

type FileEncoder interface {
	Encode(contents []byte) []byte
}
