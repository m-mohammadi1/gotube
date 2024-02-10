package youtubeutil

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Util struct {
}

func New() *Util {
	return &Util{}
}

func (u *Util) Exchange(ctx context.Context, config *oauth2.Config, code string) (*oauth2.Token, error) {
	token, err := config.Exchange(ctx, code)

	if err != nil {
		return nil, errors.New("invalid google code")
	}

	return token, nil
}

func (u *Util) GetMineChannel(ctx context.Context, config *oauth2.Config, token oauth2.Token) (*youtube.Channel, error) {
	youtubeService, err := youtube.NewService(ctx, option.WithHTTPClient(config.Client(ctx, &token)))
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

	return chosenChannel, nil
}
