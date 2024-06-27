package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func generateUniqueEmail() (string, error) {
	mx := big.NewInt(2000000)
	n, err := rand.Int(rand.Reader, mx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("user%d@example.com", n.Int64()), nil
}

func createMultipartForm(email string) (header string, bd []byte, er error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	err := writer.WriteField("email", email)
	if err != nil {
		return "", nil, err
	}
	err = writer.Close()
	if err != nil {
		return "", nil, err
	}
	return writer.FormDataContentType(), body.Bytes(), nil
}

func main() {
	targeter := func(tgt *vegeta.Target) error {
		if time.Now().Unix()%2 == 0 {
			tgt.Method = "GET"
			tgt.URL = "http://localhost:8080/rate"
			tgt.Body = nil
			tgt.Header = nil
		} else {

			uniqueEmail, err := generateUniqueEmail()
			if err != nil {
				return err
			}
			content, body, err := createMultipartForm(uniqueEmail)
			if err != nil {
				return err
			}
			tgt.Method = "POST"
			tgt.URL = "http://localhost:8080/subscribe"
			tgt.Body = body
			tgt.Header = make(http.Header)
			tgt.Header.Set("Content-Type", content)
		}
		return nil
	}

	duration := 2 * time.Minute
	initialRate := 10

	pacer := vegeta.LinearPacer{
		StartAt: vegeta.Rate{Freq: initialRate, Per: time.Second},
		Slope:   1,
	}

	attacker := vegeta.NewAttacker()

	metrics := &vegeta.Metrics{}

	results := attacker.Attack(targeter, pacer, duration, "Load Test")
	for res := range results {
		metrics.Add(res)
	}

	metrics.Close()
	file, err := os.Create("vegeta-report.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	res := vegeta.NewTextReporter(metrics)
	err = res.Report(file)
	if err != nil {
		log.Println(err)
	}
}
