package main

import (
	"bytes"
	"encoding/json"
	"favro"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {

	issues := make(chan Issue)
	var wg sync.WaitGroup
	fmt.Print("Enter JIRA Password: ")
	jiraPassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Print("\nEnter Favro API key: ")
	favroApiKey, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	cards := getAllCardsFromFavro(string(favroApiKey))

	go getIssuesFromJIRA(string(jiraPassword), issues)
	go func() {
		for issue := range issues {
			wg.Add(1)
			done := make(chan bool)
			if _, exists := cards[issue.Key]; !exists {
				card := convertJIRAToFavro(issue)
				go submitCardToFavro(string(favroApiKey), card, done)
			}
			<-done
			wg.Done()
		}
	}()

	wg.Wait()
	fmt.Println("Successfully printed cards")
	fmt.Scanln()
}

func convertJIRAToFavro(issue Issue) favro.PostCard {
	url := fmt.Sprintf("[%s](https://jira.softwire.com/jira/browse/%s): %s", issue.Key, issue.Key, issue.Fields.Description)
	var tags []favro.Tag
	colors := map[string]string{
		"Permissions": "red",
		"Airthings":   "green",
		"Bookings":    "blue",
		"Referrals":   "cyan",
		"Customers":   "orange",
		"Tech debt":   "gray",
	}
	for _, component := range issue.Fields.Components {
		color, ok := colors[component.Name]
		if ok {
			tags = append(tags, favro.Tag{Name: component.Name, Color: color})
		}
	}

	postCard := favro.PostCard{
		WidgetCommonID:      "6e542518a26d446a66e327a0",
		Name:                issue.Key,
		DetailedDescription: url,
	}

	if tags != nil {
		postCard.Tags = tags
	} else {
		postCard.Tags = make([]favro.Tag, 0)
	}

	return postCard
}

func getIssuesFromJIRA(password string, issues chan Issue) {
	jiraURL := "https://jira.softwire.com/jira/rest/api/2/search?jql=status%20in%20(Backlog%2C%20%22Selected%20for%20Development%22%2C%20%22To%20Do%22%2C%20%22Ready%20for%20dev%22%2C%20Reopened%2C%20Blocked%2C%20%22In%20Progress%22%2C%20%22Ready%20for%20review%22%2C%20%22In%20Review%22%2C%20%22Ready%20for%20test%22)%20AND%20assignee%20in%20(currentUser())"
	req, err := http.NewRequest("GET", jiraURL, nil)

	if err != nil {
		panic(err)
	}

	req.SetBasicAuth("ANP", password)

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	var jiraResponse JIRASearchResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&jiraResponse); err != nil {
		panic(err)
	}

	for _, issue := range jiraResponse.Issues {
		issues <- issue
	}

	close(issues)
}

func getAllCardsFromFavro(apiKey string) map[string]favro.Card {
	url := "https://favro.com/api/v1/cards?collectionId=9dd056e03605ab662a010ac2"
	resp, _ := sendRequestToFavro("GET", url, nil, apiKey)

	defer resp.Body.Close()
	var favroResponse favro.FavroResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&favroResponse); err != nil {
		panic(err)
	}
	cardMap := make(map[string]favro.Card)
	for _, card := range favroResponse.Cards {
		cardMap[card.Name] = card
	}
	return cardMap
}

func submitCardToFavro(apiKey string, card favro.PostCard, done chan bool) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(card)

	url := "https://favro.com/api/v1/cards"
	_, err := sendRequestToFavro("POST", url, b, apiKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(card.Name)
	done <- true
}

func sendRequestToFavro(method string, url string, body io.Reader, apiKey string) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("organizationid", "77200328d88ae2549a16c2ff")
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("andrew.paine@softwire.com", apiKey)

	return http.DefaultClient.Do(req)
}

type JIRASearchResponse struct {
	Expand     string  `json:"expand"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

type Issue struct {
	Key    string `json:"key"`
	Fields struct {
		Components []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"components"`
		Description string `json:"description"`
		Status      struct {
			Name string `json:"name"`
		} `json:"status"`
	} `json:"fields"`
}
