package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	// "strconv"
	// "unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string) (interface{}, error) {
	if bencodedString[0] == 'i' && bencodedString[len(bencodedString)-1] == 'e' {
		// var firstColonIndex int
		//
		// for i := 0; i < len(bencodedString); i++ {
		// 	if bencodedString[i] == ':' {
		// 		firstColonIndex = i
		// 		break
		// 	}
		// }
		//
		// lengthStr := bencodedString[:firstColonIndex]
		//
		number, err := strconv.Atoi(bencodedString[1 : len(bencodedString)-1])
		if err != nil {
			return "", err
		}

		return number, nil
	} else {
		return "", fmt.Errorf(bencodedString)
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		// Uncomment this block to pass the first stage

		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)
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
