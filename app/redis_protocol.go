package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Type byte

const (
	SimpleString Type = '+'
	BulkString   Type = '$'
	Array        Type = '*'
)

type Value struct {
	typ   Type
	bytes []byte
	array []Value
}

func (v Value) String() string {
	if v.typ == BulkString || v.typ == SimpleString {
		return string(v.bytes)
	}

	return ""
}

func (v Value) Array() []Value {
	if v.typ == Array {
		return v.array
	}

	return []Value{}
}

func HandleWriteError(err error) {
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(3)
	}
}

func DecodeRESP(byteStream *bufio.Reader) (Value, error) {
	dataTypeByte, err := byteStream.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(dataTypeByte) {
	case "+":
		return DecodeSimpleString(byteStream)
	case "$":
		return DecodeBulkString(byteStream)
	case "*":
		return DecodeArray(byteStream)
	}

	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(dataTypeByte))
}

func DecodeSimpleString(byteStream *bufio.Reader) (Value, error) {
	readBytes, err := ReadUntilCRLF(byteStream)
	if err != nil {
		return Value{}, err
	}

	return Value{
		typ:   SimpleString,
		bytes: readBytes,
	}, nil
}

func DecodeBulkString(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := ReadUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("falied to read bulk string length: %s", err)
	}

	readBytes := make([]byte, count+2)

	if _, err := io.ReadFull(byteStream, readBytes); err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string contents: %s", err)
	}

	return Value{
		typ:   BulkString,
		bytes: readBytes[:count],
	}, nil
}

func DecodeArray(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := ReadUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("falied to read bulk string length: %s", err)
	}

	var array []Value
	for i := 0; i < count; i++ {
		value, err := DecodeRESP(byteStream)
		if err != nil {
			return Value{}, err
		}

		array = append(array, value)
	}

	return Value{
		typ:   Array,
		array: array,
	}, nil
}

func ReadUntilCRLF(byteStream *bufio.Reader) ([]byte, error) {
	readBytes := []byte{}

	for {
		b, err := byteStream.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		readBytes = append(readBytes, b...)

		if len(readBytes) <= 2 && readBytes[len(readBytes)-2] == '\r' {
			break
		}

	}
	return readBytes[:len(readBytes)-2], nil
}
