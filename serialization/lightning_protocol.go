package serialization

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

const (
	L_BOOL = iota
	L_INT
	L_FLOAT
	L_STRING
	L_ARRAY
	L_MAP
)

var MAGIC_START = [...]byte{33, 44, 55, 12}
var MAGIC_END = [...]byte{23, 84, 55, 29}

// <--4bytes magicStart ------ 4bytes --------------------   k bytes -------------------------------   x bytes ---- magicEnd
//   MAGIC_START         k for params quantity          int array for next x(from k bytes) bytes             actual params    MAGIC_END

type Lightning struct{}

func MarshalArguments(arguments ...any) ([]byte, error) {
	types := make([]byte, len(arguments), len(arguments))
	buf := make([]byte, 0, len(arguments))
	paramLens := make([]int, len(arguments), len(arguments))
	buffer := bytes.NewBuffer(buf)

	buffer.Write(MAGIC_START[:])                                  // 写入开始魔数
	binary.Write(buffer, binary.BigEndian, int32(len(arguments))) // 写入k 表示接下来会有k 个参数
	for i, arg := range arguments {
		switch v := arg.(type) {
		case bool:
			err := binary.Write(buffer, binary.BigEndian, v)
			types[i] = L_BOOL
			paramLens[i] = 1
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error writing bool to buffer, %d argument caused it err = %+c", i, err))
			}
		case string:
			types[i] = L_STRING
			paramLens[i] = len(v)
			_, err := buffer.WriteString(v)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error writing string to buffer, %d argument caused it err = %+v", i, err))
			}
		case int:
			paramLens[i] = 4 // 暂定一个int 四个字节 不考虑 int16, int32, int64
			types[i] = L_INT
			err := binary.Write(buffer, binary.BigEndian, int32(v))
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error writing int to buffer, %d argument caused it err = %+v", i, err))
			}
		}
	}
	binary.Write(buffer, binary.BigEndian, paramLens) // 写入k个元素 代表接下来的k个参数分别占用多少字节
	binary.Write(buffer, binary.BigEndian, arguments) // 写入k个参数
	binary.Write(buffer, binary.BigEndian, MAGIC_END) // 写入结束魔数
	return buffer.Bytes(), nil
}

func Marshal(obj any) ([]byte, error) {
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)
	arguments := make([]any, 0, 10)

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("Must be struct")
	}

	// 这里拿到所有key 对应的value
	for i := 0; i < val.NumField(); i++ {
		fmt.Printf("i = %d, val = %v\n", i, val.Field(i))
		if !typ.Field(i).IsExported() {
			continue
		}
		arguments = append(arguments, val.Field(i).Interface())
	}
	stream, err := MarshalArguments(arguments...)
	fmt.Printf("stream = %+v", stream)
	if err != nil {
		fmt.Printf("error = %v", err)
		return nil, err
	}
	fmt.Printf("arguments = %+v", arguments)
	return stream, nil
}

func UnMarshal(obj any, data *any) error {
	// 这里需要把二进制数据解析到结构体里， 但是结构体应该有哪些key呢？
	// 每次解析的对象都是固定的吗？
	return nil
}
