package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/valyala/fasthttp"
)

var lock sync.Mutex
var counter int64
var reqcout int64
var usernames []string

func reqSec() {
	before := counter
	time.AfterFunc(time.Second, func() {
		reqcout = counter - before
		reqSec()
	})
}

func threadPrint() {
	for true {
		fmt.Printf("Attempts: %d | R/S %d  \r", counter, reqcout)
		time.Sleep(time.Millisecond * 150)
	}
}

func setupClaim(usrname string, AuthToken string) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://api.twitter.com/1.1/account/update_profile.json?screen_name=" + usrname)
	req.Header.SetMethod("POST")
	req.Header.Add("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAAIWCCAAAAAAA2C25AxqI%2BYCS7pdfJKRH8Xh19zA%3D8vpDZzPHaEJhd20MKVWp3UR38YoPpuTX7UD2cVYo3YNikubuxd")
	req.Header.Add("X-CSRF-Token", "83368f29e6d092aacef9e4b10b0185ab")
	req.Header.Add("Cookie", "auth_token="+AuthToken+"; ct0=83368f29e6d092aacef9e4b10b0185ab")
	return req
}

func setupCheck(usrname string) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://cdn.syndication.twitter.com/timeline/profile.json?min_position=-1&suppress_response_codes=-1&screen_name=" + usrname)
	req.Header.SetMethod("HEAD")
	return req
}

func CheckExists(checkReq *fasthttp.Request, claimReq *fasthttp.Request, unames []string) {
	client := fasthttp.Client{}
	res := fasthttp.AcquireResponse()
	lenUserNames := int64(len(unames))
	i := int64(0)
	for true {
		username := unames[i]
		i = (i + 1) % lenUserNames
		checkReq.SetRequestURI("https://cdn.syndication.twitter.com/timeline/profile.json?min_position=-1&suppress_response_codes=-1&screen_name=" + username)
		client.Do(checkReq, res)
		cl := res.Header.ContentLength()
		if cl == len(username)+83 {
			lock.Lock()
			fmt.Println("Locked For: @" + username)
			claimReq.SetRequestURI("https://api.twitter.com/1.1/account/update_profile.json?screen_name=" + username)
			ClaimUsername(claimReq, &client, username)
		} else if cl == 18 {
			client.CloseIdleConnections()
		}
		counter++
	}
}

func ClaimUsername(prebuiltReq *fasthttp.Request, client *fasthttp.Client, username string) {
	bodyBytes := []byte("screen_name=" + username)
	prebuiltReq.SetBody(bodyBytes)
	res := fasthttp.AcquireResponse()
	client.Do(prebuiltReq, res)
	claimCode := res.StatusCode()
	if claimCode == 200 {
		fmt.Println("Claimed: @" + username + " after " + humanize.Comma(int64(counter)) + " attempts.            \n")
	} else if claimCode == 403 {
		fmt.Println("Missed: @" + username + " after " + humanize.Comma(int64(counter)) + " attempts.             \n")
	}
	os.Exit(0)
}

func readLines() ([]string, error) {
	file, err := os.Open("usernames.txt")
	if err != nil {
		panic("no username file")
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
func main() {
	fmt.Println("Twitter AC")
	var err error
	usernames, err = readLines()
	if err != nil {
		fmt.Println("sth bad happened!")
		fmt.Println(err.Error())
		return
	}
	fmt.Println("# of usernames: " + strconv.Itoa(len(usernames)))
	var AuthToken string
	fmt.Print("Auth Token: ")
	fmt.Scanln(&AuthToken)

	var threads int
	fmt.Print("Threads: ")
	fmt.Scanln(&threads)

	go reqSec()
	go threadPrint()

	fmt.Println()

	checkReq := setupCheck(usernames[0])
	claimReq := setupClaim(usernames[0], AuthToken)
	thradsToUsernames := len(usernames) / threads
	for i := 0; i < threads; i++ {
		checkReq = setupCheck(usernames[i%len(usernames)])
		claimReq = setupClaim(usernames[i%len(usernames)], AuthToken)
		if i == threads-2 {
			go CheckExists(checkReq, claimReq, usernames[thradsToUsernames*i:])
			continue
		}
		go CheckExists(checkReq, claimReq, usernames[thradsToUsernames*i:thradsToUsernames*(i+1)])
	}

	fmt.Scanln()
	fmt.Println("Exited successfully.               ")
}
