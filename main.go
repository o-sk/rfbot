package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/o-sk/rfbot/config"
	"github.com/slack-go/slack"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
)

func makeFilter(ngWords []string) (*regexp.Regexp, error) {
	if len(ngWords) == 0 {
		return nil, nil
	}

	ws := make([]string, len(ngWords))
	for i, w := range ngWords {
		ws[i] = regexp.QuoteMeta(w)
	}
	return regexp.Compile("(" + strings.Join(ws, "|") + ")")
}

func main() {
	app := &cli.App{
		Name: "slackbot-sample",
		Action: func(c *cli.Context) error {
			cfg := config.Load("config.yml")
			api := slack.New(cfg.Slack.APIToken)

			filter, err := makeFilter(cfg.Filter.NgWords)
			if err != nil {
				return err
			}

			okUsers := make(map[string]struct{}, len(cfg.Filter.OkUsers))
			for _, u := range cfg.Filter.OkUsers {
				okUsers[u] = struct{}{}
			}

			rtm := api.NewRTM(slack.RTMOptionConnParams(url.Values{
				"batch_presence_aware": {"1"},
			}))
			go rtm.ManageConnection()

			for msg := range rtm.IncomingEvents {
				switch ev := msg.Data.(type) {
				case *slack.MessageEvent:
					if ev.Channel == cfg.Redirect.FromChannel {
						if ev.SubType != "" || ev.ThreadTimestamp != "" {
							continue
						}

						if _, ok := okUsers[ev.User]; !ok &&
							filter != nil && filter.FindString(ev.Text) != "" {
							continue
						}

						text := fmt.Sprintf("https://%s.slack.com/archives/%s/p%s",
							cfg.Slack.Team,
							ev.Channel,
							strings.Join(strings.Split(ev.Timestamp, "."), ""),
						)
						rtm.SendMessage(rtm.NewOutgoingMessage(text, cfg.Redirect.ToChannel))
					}

				case *slack.RTMError:
					fmt.Printf("Error: %s\n", ev.Error())

				case *slack.InvalidAuthEvent:
					return xerrors.New("Invalid credentials")

				default:
				}
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
