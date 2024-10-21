package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

var _ = json.Marshal

func decodeInt(bencodedString string) (int, error) {
	return strconv.Atoi(bencodedString[1 : len(bencodedString)-1])
}

func decodeStr(bencodedString string) (string, error) {
	var firstColonIndex int

	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}

	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
}

func decodeList(bencodedString string) ([]interface{}, error) {
	var decodedList []interface{}
	bencodedString = bencodedString[1 : len(bencodedString)-1]

	for len(bencodedString) > 0 {
		decoded, rest, err := decodeBencode(bencodedString)
		if err != nil {
			return nil, err
		}
		decodedList = append(decodedList, decoded)
		bencodedString = rest
	}
	return decodedList, nil
}

func decodeBencode(bencodedString string) (interface{}, string, error) {
	if len(bencodedString) == 0 {
		return nil, "", fmt.Errorf("empty string")
	}

	first := bencodedString[0]

	switch {
	case first == 'l':
		list, err := decodeList(bencodedString)
		if err != nil {
			return nil, "", err
		}
		return list, bencodedString[len(bencodedString):], nil
	case first >= '0' && first <= '9':
		str, err := decodeStr(bencodedString)
		if err != nil {
			return nil, "", err
		}
		return str, bencodedString[len(str)+len(strconv.Itoa(len(str)))+1:], nil
	case first == 'i':
		num, err := decodeInt(bencodedString)
		if err != nil {
			return nil, "", err
		}
		return num, bencodedString[len(strconv.Itoa(num))+2:], nil
	default:
		return nil, "", fmt.Errorf("unknown type: %s", string(first))
	}
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
