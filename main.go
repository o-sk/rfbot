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

func main() {
	app := &cli.App{
		Name: "slackbot-sample",
		Action: func(c *cli.Context) error {
			cfg := config.Load("config.yml")
			api := slack.New(cfg.Slack.APIToken)

			var filter *regexp.Regexp
			if len(cfg.Filter.NgWords) > 0 {
				ws := make([]string, len(cfg.Filter.NgWords))
				for i, w := range cfg.Filter.NgWords {
					ws[i] = regexp.QuoteMeta(w)
				}
				var err error
				filter, err = regexp.Compile("(" + strings.Join(ws, "|") + ")")
				if err != nil {
					return err
				}
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

						if filter != nil && filter.FindString(ev.Text) != "" {
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
