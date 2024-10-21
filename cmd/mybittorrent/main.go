package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

var _ = json.Marshal

func decodeList(bencodedString string) ([]interface{}, string, error) {
	decodedList := make([]interface{}, 0)
	bencodedString = bencodedString[1:] // Remove the leading 'l'

	for len(bencodedString) > 0 && bencodedString[0] != 'e' {
		decoded, rest, err := decodeBencode(bencodedString)
		if err != nil {
			return nil, "", err
		}
		decodedList = append(decodedList, decoded)
		bencodedString = rest
	}

	if len(bencodedString) == 0 {
		return nil, "", fmt.Errorf("unterminated list")
	}

	return decodedList, bencodedString[1:], nil // Remove the trailing 'e'
}

func decodeBencode(bencodedString string) (interface{}, string, error) {
	if len(bencodedString) == 0 {
		return nil, "", fmt.Errorf("empty string")
	}

	switch bencodedString[0] {
	case 'i':
		num, rest, err := decodeInt(bencodedString)
		return num, rest, err
	case 'l':
		list, rest, err := decodeList(bencodedString)
		return list, rest, err
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		str, rest, err := decodeStr(bencodedString)
		return str, rest, err
	case 'e':
		return nil, bencodedString[1:], nil // Handle empty list or end of list
	default:
		return nil, "", fmt.Errorf("unknown type: %c", bencodedString[0])
	}
}

func decodeInt(bencodedString string) (int64, string, error) {
	endIndex := strings.IndexByte(bencodedString, 'e')
	if endIndex == -1 {
		return 0, "", fmt.Errorf("invalid integer encoding")
	}
	num, err := strconv.ParseInt(bencodedString[1:endIndex], 10, 64)
	if err != nil {
		return 0, "", err
	}
	return num, bencodedString[endIndex+1:], nil
}

func decodeStr(bencodedString string) (string, string, error) {
	colonIndex := strings.IndexByte(bencodedString, ':')
	if colonIndex == -1 {
		return "", "", fmt.Errorf("invalid string encoding")
	}

	length, err := strconv.Atoi(bencodedString[:colonIndex])
	if err != nil {
		return "", "", err
	}

	if colonIndex+1+length > len(bencodedString) {
		return "", "", fmt.Errorf("string length mismatch")
	}

	str := bencodedString[colonIndex+1 : colonIndex+1+length]
	rest := bencodedString[colonIndex+1+length:]
	return str, rest, nil
}

func main() {
	command := os.Args[1]

	if command == "decode" {

		bencodedValue := os.Args[2]

		decoded, _, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}
		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
