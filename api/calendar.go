package api

import (
	"context"

	"github.com/labstack/gommon/log"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func TestApi(client_id, client_secret string) string {
	var ctx = context.Background()
	var client = InitConfig(client_id, client_secret)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Info(err)
	}

	log.Info(srv)

	return "test done"
}
