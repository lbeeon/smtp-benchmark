package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/mail"
	"net/smtp"
	"os"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	fWorker  = flag.Int("workers", 1, "Numbers of workers")
	fSendNum = flag.Int("nums", 1, "Numbers of mails")
	// if seconds == -1 ref nums
	// else send the periods
	fsendDuration = flag.Int("seconds", -1, "")
	fHost         = flag.String("host", "", "Target MAT")
	fEmlFile      = flag.String("eml", "", "EML file")
	fArfFile      = flag.String("arf", "", "ARF file")
	msg           *mail.Message
)

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func init() {
	flag.IntVar(fWorker, "w", 1, "Numbers of workers")
	flag.IntVar(fSendNum, "n", 1, "Numbers of mails")
	flag.IntVar(fsendDuration, "s", -1, "")
	flag.StringVar(fHost, "h", "", "")
	flag.StringVar(fEmlFile, "e", "", "")
	flag.StringVar(fArfFile, "f", "", "")
}

func sendMail(host string) int {
	c, err := smtp.Dial(host)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer c.Quit()
	c.Mail(msg.Header.Get("From"))
	// c.Mail("sender@example.org")

	rcptList, err := mail.ParseAddressList(msg.Header.Get("To"))
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range rcptList {
		c.Rcpt(v.Address)
	}

	// c.Rcpt("recipient@example.net")
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Println(err)
		return 0
	}
	// buf := bytes.NewBufferString("This is the email body. \n.")
	if _, err = io.Copy(wc, msg.Body); err != nil {
		log.Println(err)
		return 0
	}
	// if _, err = buf.WriteTo(wc); err != nil {
	// 	log.Println(err)
	// 	return 0
	// }
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
	fmt.Println("End:", e)
	fmt.Println("Total Success:", countSuccess)
	fmt.Println("Total Failure:", countFail)
	fmt.Println("Throughput:", float32(countSuccess)/float32(end-start))
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
	fmt.Println("Start From:", s)
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
		msg, err = getArfMail(*fEmlFile)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(*fArfFile) > 0 {
		msg, err = getArfMail(*fArfFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg, _ = getDefaultMail()
	}

	fmt.Printf(`
		Host: %s,
		Thread: %d,
		Numbers: %d,
		Periods: %d
		`, *fHost, *fWorker, *fSendNum, *fsendDuration)
	// fmt.Println(*fHost, *fWorker, *fSendNum, *fsendDuration, *fHost)
	exec(*fWorker, *fSendNum, *fsendDuration, *fHost)
}
