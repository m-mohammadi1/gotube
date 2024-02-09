package model

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"strings"
	"time"
)

type Channel struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	YoutubeID string    `json:"youtube_id"`
	Title     string    `json:"title"`
	AddedAt   time.Time `json:"added_at"`
	Token     oauth2.Token
}

type State struct {
	UserID int64 `json:"user_id"`
}

func CreateAuthUrl() {
	conf := &oauth2.Config{
		ClientID:     "91390813615-elqduhqvcle1q4mtidg405sutug8376u.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://127.0.0.1:8282/api/google/authorize/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/yt-analytics.readonly",
			"https://www.googleapis.com/auth/yt-analytics-monetary.readonly",
			"https://www.googleapis.com/auth/youtube.readonly",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/youtubepartner-channel-audit",
			"https://www.googleapis.com/auth/youtube.upload",
			"https://www.googleapis.com/auth/youtube.force-ssl",
			"https://www.googleapis.com/auth/youtubepartner",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.

	state := State{
		UserID: 1,
	}
	stateString, _ := json.Marshal(state)
	url := conf.AuthCodeURL(string(stateString) + "0")

	fmt.Printf("Visit the URL for the auth dialog: %v", url)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	log.Print("Start handling")
	conf := &oauth2.Config{
		ClientID:     "91390813615-elqduhqvcle1q4mtidg405sutug8376u.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://127.0.0.1:8282/api/google/authorize/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/yt-analytics.readonly",
			"https://www.googleapis.com/auth/yt-analytics-monetary.readonly",
			"https://www.googleapis.com/auth/youtube.readonly",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/youtubepartner-channel-audit",
			"https://www.googleapis.com/auth/youtube.upload",
			"https://www.googleapis.com/auth/youtube.force-ssl",
			"https://www.googleapis.com/auth/youtubepartner",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
	code := r.URL.Query().Get("code")
	stateString := r.URL.Query().Get("state")

	var state State
	err := json.Unmarshal([]byte(strings.Split(stateString, "0")[0]), &state)
	if err != nil {
		log.Println("error here")
		log.Println(err)
		return
	}
	log.Printf("User ID: %d", state.UserID)

	if code == "" {
		log.Print("Error getting code from request")
		http.Error(w, "failed to get get code", 400)
		return
	}

	token, err := conf.Exchange(r.Context(), code)
	if err != nil {
		log.Print("Second")
		log.Print(err)
		http.Error(w, "failed to get get code", 400)
	}

	log.Print(token)

	ctx := r.Context()

	youtubeService, err := youtube.NewService(ctx, option.WithHTTPClient(conf.Client(ctx, token)))
	if err != nil {
		log.Fatal("failed to retrieve channel data")
	}
	call := youtubeService.Channels.List([]string{"id", "snippet"}).Mine(true)
	channels, err := call.Do()
	if err != nil {
		log.Fatal("failed to retrieve channel list")
	}

	for _, channel := range channels.Items {
		log.Printf("Channel ID: %s, Title: %s", channel.Id, channel.Snippet.Title)
	}
}
