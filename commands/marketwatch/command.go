package marketwatch

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"gobunny/commands"
	"gobunny/errors"
	"gobunny/models/ticker"
	"gobunny/store"
)

var baseURL = "https://www.marketwatch.com"

type command struct {
	log    *log.Logger
	client *http.Client
	db     store.Store
}

// NewCommand returns a Command which allows searching MarketWatch
func NewCommand(logger *log.Logger, db store.Store) commands.Command {
	return &command{
		log: logger,
		client: &http.Client{
			CheckRedirect: func(request *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		db: db,
	}
}

func (c *command) Aliases() []string {
	return []string{"mw", "tick", "ticker"}
}

func (c *command) Name() string {
	return "marketwatch"
}

func (c *command) Handle(args commands.Arguments, response http.ResponseWriter, request *http.Request) error {
	if len(args) == 0 {
		http.Redirect(response, request, baseURL, http.StatusSeeOther)
		return nil
	}

	if len(args) > 1 {
		if _, err := response.Write([]byte(c.Help())); err != nil {
			return errors.NewErrResponseClosed(err)
		}
	}

	key, err := ticker.Key(args[0])
	if err != nil {
		return err
	}

	t := &ticker.Ticker{}
	if err := c.db.Get(key, t); err != nil {
		if err != store.ErrNotFound {
			return err
		}

		t, err = c.fetchTicker(key.String())
		if err != nil {
			return err
		}
	}

	http.Redirect(response, request, t.Href, http.StatusSeeOther)
	return nil
}

func (c *command) Help() string {
	return fmt.Sprintf(
		"usage: gobunny %s <ticker>",
		c.Name(),
	)
}

func (c *command) Readme() string {
	return fmt.Sprintf(
		"'gobunny %s' allows looking up securities tickers on MarketWatch\n\n"+
			"aliases: %s",
		c.Name(),
		strings.Join(c.Aliases(), ", "),
	)
}

func (c *command) fetchTicker(symbol string) (*ticker.Ticker, error) {
	params := url.Values{}
	params.Add("Lookup", symbol)
	params.Add("Type", "all")
	searchURL, err := url.Parse(
		fmt.Sprintf("%s/tools/quotes/lookup.asp?%s", baseURL, params.Encode()))

	if err != nil {
		return nil, err
	}

	response, err := c.client.Get(searchURL.String())
	if err != nil {
		return nil, err
	}

	redirectTo, err := response.Location()
	if err != nil {
		return nil, fmt.Errorf("invalid redirect: %s", err.Error())
	}

	t, err := ticker.New(ticker.WithSymbol(symbol), ticker.WithHref(redirectTo.String()))
	if err != nil {
		return nil, err
	}

	return t, nil
}
