package storage

import "fmt"

// SimpleKey simple of key
type SimpleKey struct {
	Message string `json:"message,omitempty"`
}

// SimpleVal simple of value
type SimpleVal struct {
	Message []byte `json:"message,omitempty"`
}

// SimpleCompressor simple of compressor
type SimpleCompressor struct{}

func NewSimpleCompressor() *SimpleCompressor {
	return &SimpleCompressor{}
}

// Compress key
func (k *SimpleCompressor) Compress(key interface{}) string {
	return key.(*SimpleKey).Message
}

// SimpleSerializer simple of serailizer
type SimpleSerializer struct{}

func NewSimpleSerializer() *SimpleSerializer {
	return &SimpleSerializer{}
}

// Marshal key
func (s *SimpleSerializer) Marshal(val interface{}) (buf []byte, err error) {
	return val.(*SimpleVal).Message, nil
}

// Unmarshal key
func (s *SimpleSerializer) Unmarshal(buf []byte, val interface{}) (err error) {
	pv, ok := val.(*SimpleVal)
	if !ok {
		return fmt.Errorf("val [%v] is not a type of SimpleVal", val)
	}
	pv.Message = buf
	return nil
}
