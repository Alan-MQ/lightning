package server

type Serializer interface {
	Marshal(obj any) ([]byte, error)
	UnMarshal(stream []byte, obj any) error
}
