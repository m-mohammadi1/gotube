package channelauth

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"gotube/internal/config"
	"gotube/pkg/model"
	"net/http"
	"strings"
)

type Util struct {
	config config.Data
}

func NewUtil(config config.Data) Util {
	return Util{config: config}
}

type state struct {
	UserID int64 `json:"user_id"`
}

// GenerateAuthUrl create authorization link
func (u *Util) GenerateAuthUrl(user model.User) string {
	conf := u.getOauthConf()

	state := state{
		UserID: user.ID,
	}
	stateString, _ := json.Marshal(state)
	url := conf.AuthCodeURL(string(stateString) + "0")

	return url
}

// HandleCallback handle redirect back
func (u *Util) HandleCallback(r *http.Request) (*model.Channel, error) {
	conf := u.getOauthConf()

	code := r.URL.Query().Get("code")
	stateString := r.URL.Query().Get("state")

	var state state
	err := json.Unmarshal([]byte(strings.Split(stateString, "0")[0]), &state)
	if err != nil {
		return nil, errors.New("invalid callback params: state")
	}

	if code == "" {
		return nil, errors.New("invalid callback params: code")
	}

	token, err := conf.Exchange(r.Context(), code)
	if err != nil {
		return nil, errors.New("invalid google code")
	}

	// retrieve channel data
	ctx := context.Background()

	youtubeService, err := youtube.NewService(ctx, option.WithHTTPClient(conf.Client(ctx, token)))
	if err != nil {
		return nil, errors.New("failed to retrieve channel data")
	}
	call := youtubeService.Channels.List([]string{"id", "snippet"}).Mine(true)
	channels, err := call.Do()
	if err != nil {
		return nil, errors.New("failed to retrieve channel list")
	}

	chosenChannel := channels.Items[0]

	if chosenChannel == nil {
		return nil, errors.New("failed to retrieve chosen channel")
	}

	// create channel in database or update it

	return &model.Channel{
		Token:     *token,
		UserID:    state.UserID,
		YoutubeID: chosenChannel.Id,
		Title:     chosenChannel.Snippet.Title,
	}, nil
}

func (u *Util) getOauthConf() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     u.config.GoogleClientID,
		ClientSecret: u.config.GoogleClientSECRET,
		RedirectURL:  u.config.GoogleCallback,
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
	return conf
}
