package utils

import "os"

func ReadFile(filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fileInfo, _ := f.Stat()
	b := make([]byte, fileInfo.Size())
	_, err = f.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func WriteFile(filename string, data []byte) {
	fPublic, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer fPublic.Close()
	_, err = fPublic.Write(data)
	if err != nil {
		panic(err)
	}
}
