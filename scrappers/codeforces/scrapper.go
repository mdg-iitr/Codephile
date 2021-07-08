package codeforces

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"

	. "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scrappers/common"
)

type Scrapper struct {
	Handle  string
	Context context.Context
}

func (s Scrapper) GetProfileInfo() types.ProfileInfo {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	var profile types.ProfileInfo
	requestUrl := "http://codeforces.com/api/user.info?handles=" + s.Handle
	var err error
	for attempt := 1; attempt < 10; attempt++ {
		time.Sleep(time.Second * time.Duration(attempt))
		data, _ := common.HitGetRequest(requestUrl)
		if data == nil {
			log.Println(errors.New("GetRequest failed. Please check connection status"))
			hub.CaptureException(errors.New("GetRequest failed. Please check connection status"))
		}
		err = json.Unmarshal(data, &profile)
		if err == nil {
			break
		}
		log.Println(err.Error())
		hub.AddBreadcrumb(&sentry.Breadcrumb{
			Category: "JSON parse error",
			Message:  string(data),
		}, nil)
	}
	if err != nil {
		hub.CaptureException(err)
	}
	return profile
}

// Calls the codeforces submission API and return the response in same format
func callCodeforcesAPI(handle string, afterIndex int, hub *sentry.Hub) (types.CodeforcesSubmissions, error) {
	url := "http://codeforces.com/api/user.status?handle=" + handle + "&from=" + strconv.Itoa(afterIndex) + "&count=50"
	fmt.Println(url)
	data, statusCode := common.HitGetRequest(url)
	if data == nil {
		return types.CodeforcesSubmissions{}, errors.New("GetRequest failed. Please check connection status")
	}
	var codeforcesSubmission types.CodeforcesSubmissions
	err := json.Unmarshal(data, &codeforcesSubmission)
	if err != nil {
		hub.AddBreadcrumb(&sentry.Breadcrumb{
			Category: "JSON parse error",
			Message:  string(data),
		}, nil)
		// Dont unnecessarily report error when API limit exceeds
		if statusCode != 503 {
			hub.CaptureException(err)
		}
		log.Println(err.Error())
		return types.CodeforcesSubmissions{}, err
	}
	return codeforcesSubmission, nil
}

//Get submissions of a user after an index.
//Returns an error if unsuccessful
//On receiving the error caller should return empty submission list
func getCodeforcesSubmissionParts(handle string, afterIndex int, hub *sentry.Hub) ([]types.Submission, error) {
	codeforcesSubmission, _ := callCodeforcesAPI(handle, afterIndex, hub)
	if codeforcesSubmission.Status != "OK" {
		log.Println("Codeforces submission could not be retrieved. Retrying...")
		var newCodeforcesSub types.CodeforcesSubmissions
		for attempt := 1; attempt < 10; attempt++ {
			time.Sleep(time.Second * time.Duration(attempt))
			newCodeforcesSub, _ = callCodeforcesAPI(handle, afterIndex, hub)
			if newCodeforcesSub.Status == "OK" {
				codeforcesSubmission = newCodeforcesSub
				break
			}
		}
		if newCodeforcesSub.Status == "" || newCodeforcesSub.Status == "FAILED" {
			hub.CaptureException(errors.New("codeforces API repeatedly returned FAILED"))
			return nil, errors.New("codeforces API repeatedly returned FAILED")
		}
	}
	submissions := make([]types.Submission, len(codeforcesSubmission.Result))
	for i, result := range codeforcesSubmission.Result {
		problem := result["problem"].(map[string]interface{})
		var status string
		switch result["verdict"].(string) {
		case "FAILED":
			status = StatusWrongAnswer
		case "OK":
			status = StatusCorrect
		case "PARTIAL":
			status = StatusPartial
		case "COMPILATION_ERROR":
			status = StatusCompilationError
		case "RUNTIME_ERROR":
			status = StatusRuntimeError
		case "WRONG_ANSWER":
			status = StatusWrongAnswer
		case "TIME_LIMIT_EXCEEDED":
			status = StatusTimeLimitExceeded
		case "MEMORY_LIMIT_EXCEEDED":
			status = StatusMemoryLimitExceeded
		default:
			status = StatusWrongAnswer
		}
		submissions[i].Status = status
		submissions[i].Language = result["programmingLanguage"].(string)
		submissions[i].Name = problem["name"].(string)
		if problem["contestId"] != nil {
			submissions[i].URL = "http://codeforces.com/problemset/problem/" + strconv.Itoa(int(problem["contestId"].(float64))) + "/" + problem["index"].(string)
		} else {
			submissions[i].URL = ""
		}
		submissions[i].CreationDate = time.Unix(int64(result["creationTimeSeconds"].(float64)), 0)
		if problem["points"] != nil {
			submissions[i].Points = int(problem["points"].(float64))
		}
		if problem["rating"] != nil {
			submissions[i].Rating = int(problem["rating"].(float64))
		}
		for _, x := range problem["tags"].([]interface{}) {
			submissions[i].Tags = append(submissions[i].Tags, x.(string))
		}
	}
	return submissions, nil
}

func (s Scrapper) GetSubmissions(after time.Time) []types.Submission {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	var oldestSubIndex, current int
	var oldestSubFound = false
	var subs []types.Submission
	//Fetch submission until oldest submission not found
	for !oldestSubFound {
		newSub, err := getCodeforcesSubmissionParts(s.Handle, current+1, hub)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		//Check for repetition of previous fetched submission
		if len(newSub) != 0 {
			for i, sub := range newSub {
				subs = append(subs, sub)
				oldestSubIndex = current + i + 1
				if sub.CreationDate.Equal(after) || sub.CreationDate.Before(after) {
					oldestSubFound = true
					break
				}
			}
			//50 submissions per page
			current += 50
		} else {
			break
		}
	}
	subs = subs[0:oldestSubIndex]
	return subs
}

func (s Scrapper) CheckHandle() (bool, error) {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	data, _ := common.HitGetRequest("http://codeforces.com/api/user.info?handles=" + url.PathEscape(s.Handle))
	var i interface{}
	err := json.Unmarshal(data, &i)
	if err != nil {
		for attempt := 1; attempt < 5; attempt++ {
			time.Sleep(time.Millisecond * time.Duration(attempt) * 100)
			data, _ := common.HitGetRequest("http://codeforces.com/api/user.info?handles=" + s.Handle)
			var i interface{}
			err = json.Unmarshal(data, &i)
			if err == nil {
				break
			}
			hub.AddBreadcrumb(&sentry.Breadcrumb{
				Category: "JSON parse error",
				Message:  string(data),
			}, nil)
			log.Println(err.Error())
		}
	}
	if err != nil {
		hub.CaptureException(err)
		return false, err
	}
	return i.(map[string]interface{})["status"] != "FAILED", nil
}
