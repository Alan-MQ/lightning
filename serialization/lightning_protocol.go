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
	result := make([]byte, 0, len(arguments))
	paramLens := make([]int, len(arguments), len(arguments))
	buffer := bytes.NewBuffer(buf)
	resultBuffer := bytes.NewBuffer(result)

	// buffer.Write(MAGIC_START[:])                                  // 写入开始魔数
	// binary.Write(buffer, binary.BigEndian, int32(len(arguments))) // 写入k 表示接下来会有k 个参数
	for i, arg := range arguments {
		switch v := arg.(type) {
		case bool:
			err := binary.Write(buffer, binary.BigEndian, v)
			types[i] = byte(L_BOOL)
			paramLens[i] = 1
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error writing bool to buffer, %d argument caused it err = %+c", i, err))
			}
		case string:
			types[i] = byte(L_STRING)
			fmt.Printf("string = %d, byte = %v\n", L_STRING, byte(L_STRING))
			paramLens[i] = len(v)
			_, err := buffer.WriteString(v)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error writing string to buffer, %d argument caused it err = %+v", i, err))
			}
		case int:
			paramLens[i] = 4 // 暂定一个int 四个字节 不考虑 int16, int32, int64
			types[i] = byte(L_INT)
			err := binary.Write(buffer, binary.BigEndian, int32(v))
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error writing int to buffer, %d argument caused it err = %+v", i, err))
			}
		}
	}
	// buffer.Write(types) // 写入k值 代表k个元素的类型
	// binary.Write(buffer, binary.BigEndian, paramLens) // 写入k个元素 代表接下来的k个参数分别占用多少字节
	// binary.Write(buffer, binary.BigEndian, arguments) // 写入k个参数
	// binary.Write(buffer, binary.BigEndian, MAGIC_END) // 写入结束魔数

	// <--4bytes magicStart ------ 4bytes --------------------   k bytes -------------------------------   x bytes ---- magicEnd
	//   MAGIC_START         k for params count          int array for next 4*k bytes             actual params    MAGIC_END

	resultBuffer.Write(MAGIC_START[:])                                  // 写入开始魔数
	binary.Write(resultBuffer, binary.BigEndian, int32(len(arguments))) // 写入k 表示接下来会有k 个参数
	resultBuffer.Write(types)                                           // 写入k值 代表k个元素的类型
	// binary.Write(resultBuffer, binary.BigEndian, paramLens)             // 写入k 个数字  表示接下来k个参数分别占用多少字节
	for _, v := range paramLens {
		binary.Write(resultBuffer, binary.BigEndian, uint32(v))
	}
	resultBuffer.Write(buffer.Bytes()) // 写入具体的参数信息
	resultBuffer.Write(MAGIC_END[:])   // 写入结束魔数
	return resultBuffer.Bytes(), nil
}

func Marshal(obj any) ([]byte, error) {
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)
	arguments := make([]any, 0, 10)

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("Must be struct")
	}

	for i := 0; i < val.NumField(); i++ {
		fmt.Printf("i = %d, val = %v\n", i, val.Field(i))
		if !typ.Field(i).IsExported() {
			continue
		}
		arguments = append(arguments, val.Field(i).Interface())
	}
	stream, err := MarshalArguments(arguments...)
	fmt.Printf("stream = %+v\n", stream)
	if err != nil {
		fmt.Printf("error = %v", err)
		return nil, err
	}
	fmt.Printf("arguments = %+v\n", arguments)
	return stream, nil
}

func UnMarshalArgument(obj any, data []byte) error {
	// TODO: 1. figure out what is the key to the obj
	// 2. if there is multiple string how to make sure the value matches to the right key?
	buffer := bytes.NewBuffer(data)
	// make sure it's lightning protocol
	header := buffer.Next(4)
	if !bytes.Equal(header, MAGIC_START[:]) {
		return errors.New("not a lightning protocol")
	}

	// 确认参数个数
	paramCount := int(binary.BigEndian.Uint32(buffer.Next(4)))
	lens := make([]int, paramCount, paramCount)
	fmt.Printf("paramCount = %d\n", paramCount)
	paramTypes := make([]byte, paramCount, paramCount)

	_, err := buffer.Read(paramTypes)
	if err != nil {
		fmt.Printf("error reading paramTypes, %v", err)
		return err
	}
	fmt.Printf("paramTypes = %+v\n", paramTypes)
	//paramLengths := make([]int, paramCount, paramCount)
	for i := range paramCount {
		k := int(paramTypes[i])
		fmt.Printf("i = %d, k = %d\n", i, k)
		paramTypes[i] = byte(k)
	}
	for i := range paramCount {
		currentParamLen := int(binary.BigEndian.Uint32(buffer.Next(4)))
		lens[i] = currentParamLen
	}
	fmt.Printf("lens = %+v\n", lens)

	return nil
}

func UnMarshal(obj any, data []byte) error {
	return UnMarshalArgument(obj, data)
}
