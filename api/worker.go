package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/thaibui2308/terminal-app/models"
)

func pageUrl(pageNumber int) string {
	url := "http://www.ratemyprofessors.com/filter/professor/?&page=" + strconv.Itoa(pageNumber) + "&filter=teacherlastname_sort_s+asc&query=*%3A*&queryoption=TEACHER&queryBy=schoolId&sid=877"
	return url
}

func GetRatings(tid string) []models.Ratings {
	var ratings models.ProfessorRating

	response, err := http.Get("https://www.ratemyprofessors.com/paginate/professors/ratings?tid=" + tid + "&filter=&courseCode=&page=1")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	json.Unmarshal(responseData, &ratings)
	ratingList := ratings.Ratings
	return ratingList[0:2]
}

func FindInstructor(first, last string) (models.Professors, error) {

	start := time.Now()
	defer func() {
		fmt.Println("Execution Time: ", time.Since(start))
	}()

	found := make(chan models.Professors)
	fail := make(chan bool)

	from, to := searchInterval(last)

	if from == 0 && to == 0 {
		return models.Professors{}, errors.New("📊Cannot find this professor in the system")
	}

	// TODO: divide the requests into smaller chunk to cover all of the entries
	go sendRequest(first, last, found, fail, from, to)

	for {
		select {
		case m := <-found:
			return m, nil
		case <-fail:
			return models.Professors{}, errors.New("Cannot find this professor in the system")
		}
	}
}

func sendRequest(first, last string, found chan models.Professors, fail chan bool, from int, to int) {
	var wg sync.WaitGroup

	// loop through the professor array and run the search func on different thread
	for i := from; i <= to; i++ {
		wg.Add(1)
		go func(id int) {
			var pageResponse models.APIResponse
			page, err := http.Get(pageUrl(id))
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			responseData, err := ioutil.ReadAll(page.Body)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			json.Unmarshal(responseData, &pageResponse)
			professorList := pageResponse.Professors

			for _, v := range professorList {
				if (v.TLname == last) && (v.TFname == first) {
					found <- v
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fail <- false
}

func searchInterval(lastname string) (from, to int) {
	if lastname == "" {
		return 0, 0
	}
	firstLetter := strings.ToUpper(lastname[0:1])
	runes := []rune(firstLetter)
	asciiVal := int(runes[0])

	if asciiVal >= 65 && asciiVal <= 71 {
		from = 1
		to = 75
	} else if asciiVal >= 71 && asciiVal <= 80 {
		from = 75
		to = 150
	} else {
		from = 150
		to = 212
	}
	return from, to
}
