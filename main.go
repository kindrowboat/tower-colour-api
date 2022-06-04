package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/karalabe/hid"
)

type ColourMessage struct {
	Red       byte      `json:"red"`
	Green     byte      `json:"green"`
	Blue      byte      `json:"blue"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

const logFileName string = "messages.txt"

func tomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		msg, err := readLastLogLine()
		if err != nil {
			log.Println(err.Error())
		}
		j, _ := json.Marshal(msg)
		w.Write(j)
	case "POST":
		decoder := json.NewDecoder(r.Body)
		colourMessage := &ColourMessage{}
		err := decoder.Decode(colourMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		colourMessage.CreatedAt = time.Now()
		changeColour(colourMessage.Red, colourMessage.Green, colourMessage.Blue)
		loggingErr := logMessage(*colourMessage)
		if loggingErr != nil {
			log.Printf(loggingErr.Error())
		}
		jsonMessage, _ := json.Marshal(colourMessage)
		w.Write(jsonMessage)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func changeColour(red, green, blue byte) error {

	devices := hid.Enumerate(0x187C, 0x0550)
	if len(devices) == 0 {
		return fmt.Errorf("Could not find HID device for colour control")
	}
	deviceInfo := devices[0]
	device, openErr := deviceInfo.Open()
	if openErr != nil {
		return openErr
	}
	defer device.Close()

	var hidError error
	_, hidError = device.Write([]byte{0x03, 0x20, 0x02}) //prepareTurn
	if hidError != nil {
		return hidError
	}
	_, hidError = device.Write([]byte{0x03, 0x26, 0x00, 0x00, 0x02, 0x00, 0x01}) //turnOn
	if hidError != nil {
		return hidError
	}
	_, hidError = device.Write([]byte{0x03, 0x21, 0x00, 0x01, 0xff, 0xff}) //reset
	if hidError != nil {
		return hidError
	}
	_, hidError = device.Write([]byte{0x03, 0x23, 0x01, 0x00, 0x02, 0x00, 0x01}) //colorSel
	if hidError != nil {
		return hidError
	}
	_, hidError = device.Write([]byte{0x03, 0x24, 0x00, 0x07, 0xd0, 0x00, 0xfa, red, green, blue}) //colorSet
	if hidError != nil {
		return hidError
	}
	_, hidError = device.Write([]byte{0x03, 0x21, 0x00, 0x03, 0x00, 0xff}) //update
	if hidError != nil {
		return hidError
	}

	return nil

}

func logMessage(msg ColourMessage) error {
	file, fileOpenErr := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
	if fileOpenErr != nil {
		return fileOpenErr
	}
	defer file.Close()
	cleanedMessage := strings.ReplaceAll(msg.Message, "\r\n", "")
	cleanedMessage = strings.ReplaceAll(cleanedMessage, "\n", "")
	cleanedMessage = strings.ReplaceAll(cleanedMessage, "\t", "")
	logLine := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\n", msg.CreatedAt.Format(time.RFC3339), msg.Red, msg.Green, msg.Blue, cleanedMessage)
	log.Printf(logLine)
	_, fileWriteErr := file.WriteString(logLine)
	if fileWriteErr != nil {
		return fileWriteErr
	}

	return nil
}

func readLastLogLine() (msg ColourMessage, err error) {
	lastLine, err := getLastLineWithSeek(logFileName)
	if err != nil {
		return
	}
	parts := strings.Split(lastLine, "\t")
	if len(parts) != 5 {
		err = fmt.Errorf("Malformed line %v", lastLine)
		return
	}
	msg.CreatedAt, _ = time.Parse(time.RFC3339, parts[0])
	iRed, _ := strconv.ParseUint(parts[1], 10, 8)
	msg.Red = byte(iRed)
	iGreen, _ := strconv.ParseUint(parts[2], 10, 8)
	msg.Green = byte(iGreen)
	iBlue, _ := strconv.ParseUint(parts[3], 10, 8)
	msg.Blue = byte(iBlue)
	msg.Message = strings.TrimRight(parts[4], "\n")
	return
}

func getLastLineWithSeek(filepath string) (string, error) {
	fileHandle, err := os.Open(filepath)

	if err != nil {
		return "", err
	}
	defer fileHandle.Close()

	line := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	for {
		cursor -= 1
		fileHandle.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		fileHandle.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
			break
		}

		line = fmt.Sprintf("%s%s", string(char), line)

		if cursor == -filesize { // stop if we are at the begining
			break
		}
	}

	return line, nil
}

func main() {
	if !hid.Supported() {
		log.Fatal("hid is not support")
	}

	if len(os.Args[1:]) == 3 {
		iRed, _ := strconv.ParseUint(os.Args[1], 10, 8)
		iGreen, _ := strconv.ParseUint(os.Args[2], 10, 8)
		iBlue, _ := strconv.ParseUint(os.Args[3], 10, 8)
		changeColour(byte(iRed), byte(iGreen), byte(iBlue))
		return
	}

	http.HandleFunc("/", tomHandler)
	addr := ":3010"
	log.Printf("Serving on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
