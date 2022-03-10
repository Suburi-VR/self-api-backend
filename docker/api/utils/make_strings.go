package utils

import (
	"crypto/rand"
)

func MakeRandomStr(digit uint32) (string) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
			return ""
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
			// index が letters の長さに収まるように調整
			result += string(letters[int(v)%len(letters)])
	}
	return result
}

func CallId() (string) {

	callID := MakeRandomStr(3) + "-" + MakeRandomStr(3) + "-" + MakeRandomStr(3)

	return callID
}

func Password() (string) {

	password := MakeRandomStr(16)

	return password
}