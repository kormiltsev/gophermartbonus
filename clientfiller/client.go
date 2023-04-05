package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// fill the table without cookies
	answer := "DONE"
	for {
		for i := 400000; i < 900000; i++ {
			switch WorkingWOcookies(fmt.Sprint(i)) {
			case 200:
				answer += fmt.Sprintf("200:%d\n", i)
				log.Println("200:", i)
			case 202:
				answer += fmt.Sprintf("202:%d\n", i)
				log.Println("202:", i)
			case 400:
				answer += fmt.Sprintf("400:%d\n", i)
			case 401:
				answer += fmt.Sprintf("401:%d\n", i)
			case 409:
				answer += fmt.Sprintf("409:%d\n", i)
			case 422:
				answer += fmt.Sprintf("422:%d\n", i)
			case 500:
				answer += fmt.Sprintf("500:%d\n", i)
			}

		}
	}
}

var Domain = "http://localhost:8080/"

func WorkingWOcookies(s string) int {
	body := bytes.NewBuffer([]byte(s))
	request, _ := http.NewRequest("POST", Domain+"api/user/orders", body)
	request.Header.Add("Authorization", "Bearer 2521ab64f32d7974945356ad627ed59d47")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
