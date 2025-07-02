package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type QA struct {
	Input  string
	Output string
}

func TextToBinary(text string) string {
	binary := ""
	for _, c := range text {
		binary += fmt.Sprintf("%08b", c)
	}
	return binary
}

func BinaryToText(binary string) string {
	var text string
	for i := 0; i < len(binary); i += 8 {
		if i+8 > len(binary) {
			break
		}
		b := binary[i : i+8]
		n, err := strconv.ParseInt(b, 2, 64)
		if err != nil {
			continue
		}
		text += string(rune(n))
	}
	return text
}

func LoadCSV(path string) ([]QA, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // skip header
	if err != nil {
		return nil, err
	}

	var records []QA
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(line) < 2 {
			continue
		}
		records = append(records, QA{Input: line[0], Output: line[1]})
	}
	return records, nil
}

func GetResponse(input string, data []QA, csvPath string) string {
	input = strings.ToLower(input)
	if len(input) > 5 {
		input = input[:5] // ambil 5 huruf pertama
	}
	binaryInput := TextToBinary(input)

	for _, qa := range data {
		if qa.Input == binaryInput {
			return BinaryToText(qa.Output)
		}
	}

	// Tidak ditemukan → tanya user
	fmt.Println("Bot : Maaf, saya tidak mengerti.")
	fmt.Print("Bot : Apa jawaban yang benar untuk input ini? ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userAnswer := scanner.Text()
	binaryOutput := TextToBinary(userAnswer)

	err := AppendUnknownInput(csvPath, binaryInput, binaryOutput)
	if err != nil {
		fmt.Println("Gagal menyimpan input baru:", err)
	}

	return userAnswer
}

func AppendUnknownInput(path string, binaryInput, binaryOutput string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{binaryInput, binaryOutput})
	return err
}

func main() {
	dataPath := "data.csv"
	data, err := LoadCSV(dataPath)
	if err != nil {
		fmt.Println("Gagal memuat data:", err)
		return
	}

	fmt.Println("🤖 Chatbot AI Sederhana Siap")
	fmt.Println("Ketik `exit` untuk keluar.")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Kamu: ")
		scanner.Scan()
		text := scanner.Text()
		if strings.ToLower(text) == "exit" {
			fmt.Println("Bot: Sampai jumpa!")
			break
		}
		response := GetResponse(text, data, dataPath)
		data, _ = LoadCSV(dataPath)
		fmt.Println("Bot :", response)
	}
}
