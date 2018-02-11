package main

import (
	"fmt"
	"math/rand"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
)

// func getArfMail(file string) (*mail.Message, error) {
// 	p, err := filepath.Abs(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	f, err := os.Open(p)
// 	if err != nil {
// 		return nil, err
// 	}
// 	msg, err = mail.ReadMessage(f)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return mail.ReadMessage(msg.Body)
// }

func getEmlMail(file string) (*mail.Message, error) {
	p, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	return mail.ReadMessage(f)
}

func getDefaultMail(size int) (*mail.Message, error) {
	msg := "Date: Mon, 23 Jun 2015 11:40:36 -0400\n" +
		"From: RobotFrom <From@example.com>\n" +
		"To: RobotTo <To@example.com>\n" +
		"Subject: Test From Smtp-benchmark\n" +
		"\n"
	msg += randStringBytes(size*1000) +
		"\n"

	r := strings.NewReader(msg)
	return mail.ReadMessage(r)
}

func randStringBytes(n int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func dumpMail(m *mail.Message) {
	header := m.Header

	fmt.Printf(`
		From: %s,
		To: %s,
		Subject: %s, 
	`, header.Get("From"), header.Get("To"), header.Get("Subject"))
	// body, err := ioutil.ReadAll(m.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s", body)
}
