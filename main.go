package main

import (
	"bytes"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var messageDir string

func handleMessage(message string) {
	buff := bytes.Buffer{}
	uid, _ := uuid.NewV4()
	nowt := time.Now()
	strtime := nowt.Format("20060102150405")
	nano := nowt.UnixNano() % nowt.Unix()

	buff.WriteString(strtime)
	buff.WriteString(strconv.FormatInt(nano, 10))
	strtime = buff.String()

	buff.Reset()
	buff.WriteString(messageDir)
	buff.WriteString("/")
	buff.WriteString(uid.String())
	buff.WriteString("_")
	buff.WriteString(strtime)
	buff.WriteString(".writing")
	fileName := buff.String()

	buff.Reset()
	buff.WriteString(messageDir)
	buff.WriteString("/")
	buff.WriteString(uid.String())
	buff.WriteString("_")
	buff.WriteString(strtime)
	buff.WriteString(".xml")
	finalFileName := buff.String()

	// log.Printf("filename = [%s]\n", fileName)

	if _, err := os.Stat(messageDir); os.IsNotExist(err) {
		os.MkdirAll(messageDir, 0755)
	}

	err := ioutil.WriteFile(fileName, []byte(message), 0644)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = os.Rename(fileName, finalFileName)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("success create file %s\n", finalFileName)
	}
}

func messagesHttpRequestHandler(w http.ResponseWriter, r *http.Request) {
	result := "<response><success>true</success></response>"
	message := r.PostFormValue("logistics_interface")
	log.Printf("logistics_interface=[%s]\n", message)
	if strings.Index(message, "<Data>") == -1 ||
		strings.Index(message, "</Data>") == -1 {
		result = "<response><success>false</success><errorCode>S07</errorCode></response>"
	} else {
		go handleMessage(message)
	}
	fmt.Fprintln(w, result)
}

func main() {
	arguments := os.Args
	port := "8080"
	if len(arguments) < 2 {
		log.Printf("usage : %s <messageDir> [<port>]", arguments[0])
		return
	} else {
		if len(arguments) > 2 {
			port = arguments[2]
		} else {
			log.Println("No port number specified, useing the default 8080 port.")
		}
		messageDir = arguments[1]
	}

	log.Printf("start listen address *.%s", port)
	http.HandleFunc("/EcssTran/httpRequest/messagesHttpRequest", messagesHttpRequestHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
