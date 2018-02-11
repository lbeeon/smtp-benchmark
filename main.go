package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"time"
)

var (
	fWorker  = flag.Int("workers", 1, "Numbers of workers")
	fSendNum = flag.Int("nums", 1, "Numbers of mails")
	// if seconds == -1 ref nums
	// else send the periods
	fsendDuration = flag.Int("seconds", -1, "")
	fBodySize     = flag.Int("size", 1, "")
	fHost         = flag.String("host", "", "Target MAT")
	fEmlFile      = flag.String("eml", "", "EML file")
	msg           *mail.Message
)

func init() {
	flag.IntVar(fWorker, "w", 1, "Numbers of workers")
	flag.IntVar(fSendNum, "n", 1, "Numbers of mails")
	flag.IntVar(fsendDuration, "s", -1, "")
	flag.IntVar(fBodySize, "b", 1, "")
	flag.StringVar(fHost, "h", "", "")
	flag.StringVar(fEmlFile, "e", "", "")
}

func sendMail(host string) int {
	c, err := smtp.Dial(host)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer c.Quit()
	c.Mail(msg.Header.Get("From"))

	rcptList, err := mail.ParseAddressList(msg.Header.Get("To"))
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range rcptList {
		c.Rcpt(v.Address)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Println(err)
		return 0
	}

	if _, err = io.Copy(wc, msg.Body); err != nil {
		log.Println(err)
		return 0
	}
	wc.Close()
	defer c.Quit()
	return 1
}

func sendWorker(job, done chan int, host string) {
	for _ = range job {
		done <- sendMail(host)
	}
}

func resultCollect(result chan int, start int64, seconds int64, nums int) {
	countSuccess, countFail := 0, 0
	if seconds == -1 {
		for {
			r := <-result
			if r == 1 {
				countSuccess++
			} else {
				countFail++
			}

			if countSuccess+countFail == nums {
				break
			}
		}
	} else {
		for {
			r := <-result
			if r == 1 {
				countSuccess++
			} else {
				countFail++
			}
			if time.Now().Unix() > start+seconds {
				break
			}
		}
	}
	e := time.Now()
	end := e.Unix()

	fmt.Printf(`
		End: %s,
		Total Success: %d,
		Total Failure: %d,
		Throughput: %f 
	`, e, countSuccess, countFail, float32(countSuccess)/float32(end-start))

	os.Exit(1)
}

func jobProducer(job chan int, count, seconds int) {
	if seconds > 0 {
		for i := 0; ; i++ {
			job <- i
		}
	} else {

		for i := 0; i < count; i++ {
			job <- i
		}
	}
}

func exec(worker int, nums int, seconds int, host string) {
	job := make(chan int)
	done := make(chan int)

	s := time.Now()
	start := s.Unix()

	fmt.Printf(`
		Start From: %s`, s)

	go resultCollect(done, start, int64(seconds), nums)

	for i := 0; i < worker; i++ {
		go sendWorker(job, done, host)
	}

	go jobProducer(job, nums, seconds)

	for {
		time.Sleep(1 * time.Second)
	}
}

func main() {
	var err error
	flag.Parse()
	if len(*fHost) == 0 {
		log.Println("Please provide host")
		os.Exit(0)
	}

	if len(*fEmlFile) > 0 {
		msg, err = getEmlMail(*fEmlFile)
		if err != nil {
			log.Fatalln("Eml:", err)
		}
	} else {
		msg, err = getDefaultMail(*fBodySize)
		if err != nil {
			log.Fatalln("Default:", err)
		}
	}

	fmt.Printf(`
		Host: %s,
		Thread: %d,
		Numbers: %d,
		Periods: %d,

		`, *fHost, *fWorker, *fSendNum, *fsendDuration)
	dumpMail(msg)

	exec(*fWorker, *fSendNum, *fsendDuration, *fHost)
}
