package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Holds values to check the bytes starting at the most significat bit to the least significat bit
var lookUp = []byte{128, 64, 32, 16, 8, 4, 2, 1}
var encryption byte

// Will be a simple XOR on each byte if a password is provided
var password = flag.String("password", "", "Password to encrypt your message")

func main() {
	encode := flag.Bool("e", false, "encode a message into the bmp file")
	decode := flag.Bool("d", false, "Decode a message from the bmp file")
	message := flag.String("m", "", "Message to be encoded")
	picture := flag.String("p", "", "BMP file name to encode or decode")
	output := flag.String("o", "", "Output file name")

	flag.Parse()
	if *password != "" {
		encryption = generatePassword(*password)
	}

	if *encode {
		if *message == "" || *picture == "" {
			fmt.Println("Error you need to specify \"-p\" and \"-m\"")
		} else {
			encodeMessage(*picture, *message, *output)
		}
	} else if *decode {
		if *picture == "" {
			fmt.Println("Error you need to specidy \"-p\"")
		} else {
			decodeMessage(*picture)
		}
	} else if *picture == "" {
		fmt.Println("Error please specify -p")
	} else {
		printLength(*picture)
	}
}

func encodeMessage(fileName string, message string, output string) {
	picture, err := ioutil.ReadFile(fileName)
	m := []byte(message)
	if len(message)*8 > len(picture)-54 {
		fmt.Println("Error " + fileName + " is not large enough to hold this message")
		return
	}
	if err != nil {
		fmt.Println("Error, could not open " + fileName)
		return
	}

	// Loop through the message
	for i := 0; i < len(m); i++ {
		// Determine where in the picture we are going to start
		index := 55 + (i * 8)
		b := m[i]
		if *password != "" {
			b = b ^ encryption
		}
		// Loop through each bit of the message using our lookup table to get which bit is set
		for j := 0; j < 8; j++ {
			if b&lookUp[j] == 0 {
				picture[index+j] = setLSB(0, picture[index+j])
			} else {
				picture[index+j] = setLSB(1, picture[index+j])
			}
		}
		// If this is the last run in the loop set the next 8 bits to 0 to terminate the message
		if i == len(m)-1 {
			for j := 8; j < 16; j++ {
				picture[index+j] = setLSB(0, picture[index+j])
			}
		}
	}
	// Write the encoded message back
	if output != "" {
		if !strings.HasSuffix(output, ".bmp") {
			output += ".bmp"
		}
		f, _ := os.Create(output)
		f.Write(picture)
		defer f.Close()
	} else {
		f, _ := os.Create("test.bmp")
		f.Write(picture)
		defer f.Close()
	}
}

// Read the LSB of each byte of the specifed file until we reach the end of the message
func decodeMessage(fileName string) {
	picture, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error could not open " + fileName)
	}
	// message will hold the message that is printed
	message := ""
	// Loop through the bmp file starting at byte 55
	for i := 55; i < len(picture)-9; i += 8 {
		var letter byte
		for j := 0; j < 8; j++ {
			b := picture[i+j]
			if b%2 == 0 {
				letter &^= 1
			} else {
				letter |= 1
			}
			// Bit shift left 1 unless it is our last bit
			if j != 7 {
				letter = letter << 1
			}
		}
		// If letter is 0 that means it has reached the end of the file
		if letter == 0 {
			break
		}
		if *password != "" {
			letter = letter ^ encryption
		}
		// Append the decoded character to message
		message += string(letter)
	}
	fmt.Println(message)
}

// Print out the number of characters the specifed file can hold
func printLength(fileName string) {
	picture, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error could not open " + fileName)
	} else {
		fmt.Printf("%s can hold %v characters\n", fileName, (len(picture)/8)-54)
	}
}

func setLSB(b byte, val byte) byte {
	if b != 0 {
		val |= 1
	} else {
		val &^= 1
	}
	return val
}

func generatePassword(password string) byte {
	byteArray := []byte(password)
	var code byte
	for _, v := range byteArray {
		code += v
	}
	return code
}
