package channelauth

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
	"gotube/internal/config"
	"gotube/pkg/model"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestUtil_GenerateAuthUrl(t *testing.T) {
	conf := getFakeConfig()

	util := New(conf, nil, nil)

	user := model.User{
		ID:        1,
		Name:      "Mohammad",
		Email:     "m.m@test.com",
		Password:  "12345678",
		CreatedAt: time.Time{},
	}

	url := util.GenerateAuthUrl(user)

	assert.NotEmpty(t, url)
}

func getFakeConfig() config.Data {
	conf := config.Data{
		"test",
		60,
		"test.com",
		"8080",
		"test",
		"test",
		"test.com",
	}
	return conf
}

func TestUtil_HandleCallback(t *testing.T) {
	conf := config.Data{
		"test",
		60,
		"test.com",
		"8080",
		"test",
		"test",
		"test.com",
	}

	util := New(conf, &FakeTokenExchanger{}, &FakeChannelLister{})

	params := url.Values{}
	// the time is added by google and I don't know why
	state := state{
		UserID: 1,
	}
	stateString, _ := json.Marshal(state)
	params.Set("state", string(stateString)+"||")
	params.Set("code", "test-code")
	req, err := http.NewRequest("GET", "http://text.com/callack?"+params.Encode(), nil)
	assert.NoError(t, err)

	channel, err := util.HandleCallback(req)

	assert.NoError(t, err)
	assert.NotNil(t, channel)
	assert.Equal(t, "test-id", channel.YoutubeID)
	assert.Equal(t, "test-title", channel.Title)
}

type FakeTokenExchanger struct {
}

func (f *FakeTokenExchanger) Exchange(ctx context.Context, config *oauth2.Config, code string) (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "test-access-token-type",
		RefreshToken: "test-refresh-token",
		Expiry:       time.Now(),
	}, nil
}

type FakeChannelLister struct {
}

func (f *FakeChannelLister) GetMineChannel(ctx context.Context, config *oauth2.Config, token oauth2.Token) (*youtube.Channel, error) {
	return &youtube.Channel{
		Id: "test-id",
		Snippet: &youtube.ChannelSnippet{
			Title: "test-title",
		},
	}, nil
}
