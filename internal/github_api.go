package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GithubApi stores values used when calling Github's APIs
type GithubApi struct {
	Owner   string
	Repo    string
	Auth    string
	Version string
}

// WorkflowRun stores the process list of the repos workflow runs
type WorkflowRun struct {
	ID         int
	Title      string
	Name       string
	Branch     string
	Status     string
	Conclusion string
	Actor      string
}

// workflowResponse stores the response of workflow run list
type workflowResponse struct {
	Count int `json:"total_count"`
	Runs  []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Branch     string `json:"head_branch"`
		Title      string `json:"display_title"`
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
		Actor      struct {
			Name string `json:"login"`
		} `json:"actor"`
	} `json:"workflow_runs"`
}

// workflowError struct to parse error response from Github's API
type workflowError struct {
	Doc string `json:"documentation_url"`
	Msg string `json:"message"`
}

// apiCall generic method to call API. You have to specify the http method, url and if its to process the response body
func (gh *GithubApi) apiCall(method, url string, processBody bool) (int, []byte) {
	client := http.Client{}
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 1000, []byte(err.Error())
	}
	request.Header = http.Header{
		"Accept":               {"application/vnd.github+json"},
		"X-GitHub-Api-Version": {gh.Version},
		"Authorization":        {gh.Auth},
	}
	response, err := client.Do(request)
	if err != nil {
		return 1001, []byte(err.Error())
	}
	if !processBody {
		return response.StatusCode, nil
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 1002, []byte(err.Error())
	}
	return response.StatusCode, body
}

// ListWorkflows calls API to retrieves the list of workflow run of the specified repos
func (gh *GithubApi) ListWorkflows() ([]WorkflowRun, error) {
	runs := make([]WorkflowRun, 0)
	page := 1
	for {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/actions/runs?per_page=100&page=%d",
			gh.Owner,
			gh.Repo,
			page,
		)
		fmt.Println(url)
		code, data := gh.apiCall("GET", url, true)
		// check for error code
		if code != 200 {
			var response workflowError
			err := json.Unmarshal(data, &response)
			if err != nil {
				return runs, err
			}
			return runs, fmt.Errorf(response.Msg)
		}
		// parse response
		var response workflowResponse
		err := json.Unmarshal(data, &response)
		if err != nil {
			return runs, err
		}
		// check if there is data to process
		if len(response.Runs) == 0 {
			return runs, nil
		}
		// process/filter response
		for _, run := range response.Runs {
			runs = append(runs, WorkflowRun{
				ID:         run.ID,
				Title:      run.Title,
				Name:       run.Name,
				Branch:     run.Branch,
				Status:     run.Status,
				Conclusion: run.Conclusion,
				Actor:      run.Actor.Name,
			})
		}
		// check if we reach the end, saving the need to make a another api call
		if len(runs) == response.Count {
			return runs, nil
		}
		page += 1
	}
}

// DeleteWorkflow calls API to delete workflow by it's id
func (gh *GithubApi) DeleteWorkflow(runID int) error {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/runs/%d",
		gh.Owner,
		gh.Repo,
		runID,
	)
	code, _ := gh.apiCall("DELETE", url, false)
	// check for error code
	if code != 204 {
		return fmt.Errorf("http error code: %d", code)
	}
	return nil
}
