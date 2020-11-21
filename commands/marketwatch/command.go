package marketwatch

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"gobunny/commands"
	"gobunny/handlers"
	"gobunny/models/ticker"
	"gobunny/store"

	"github.com/go-chi/chi"
)

const baseURL = "https://www.marketwatch.com"

var possibleSecurities = []string{
	"stock",
	"cryptocurrency",
	"index",
	"fund",
	"currency",
}

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

func (c *command) Handle(response http.ResponseWriter, request *http.Request) error {
	var tick string
	if tick = chi.URLParam(request, "ticker"); len(tick) == 0 {
		return handlers.Redirect(response, request, baseURL)
	}

	key, err := ticker.Key(tick)
	if err != nil {
		return handlers.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	t := &ticker.Ticker{}
	if err := c.db.Get(key, t); err != nil {
		if err != store.ErrNotFound {
			fmt.Printf("unexpected error getting ticker from store: %s", err.Error())
			return handlers.NewHTTPError("", http.StatusInternalServerError)
		}

		t, err = c.fetchTicker(tick)
		if err != nil {
			return err
		}

		if err := c.db.Create(t.Key(), t); err != nil {
			fmt.Printf("unexpected error creating ticker in store: %s", err.Error())
			return handlers.NewHTTPError("", http.StatusInternalServerError)
		}
	}

	return handlers.Redirect(response, request, t.Href)
}

func (c *command) Routes() []commands.Route {
	return []commands.Route{
		{
			Method: http.MethodGet,
			Patterns: []string{
				`marketwatch`,
				`marketwatch {ticker}`,
				`ticker {ticker}`,
				`tick {ticker}`,
			},
		},
	}
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

	var redirectTo *url.URL
	if response.StatusCode == http.StatusFound {
		redirectTo, err = response.Location()
		if err != nil {
			return nil, fmt.Errorf("invalid redirect: %s", err.Error())
		}
	} else {
		for _, security := range possibleSecurities {
			securityURL, found := c.tryFetchSecurity(symbol, security)
			if found {
				redirectTo = securityURL
				break
			}
		}
	}

	if redirectTo == nil {
		return nil, handlers.NewHTTPError("ticker not found", http.StatusNotFound)
	}

	t, err := ticker.New(ticker.WithSymbol(symbol), ticker.WithHref(redirectTo.String()))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (c *command) tryFetchSecurity(symbol string, security string) (*url.URL, bool) {
	s := fmt.Sprintf("%s/investing/%s/%s", baseURL, security, symbol)
	securityURL, err := url.Parse(s)
	if err != nil {
		c.log.Printf("malformed MarketWatch URL '%s': %s", s, err.Error())
		return nil, false
	}

	response, err := c.client.Get(securityURL.String())
	if err != nil {
		c.log.Printf("error during HTTP client GET '%s': %s", securityURL.String(), err.Error())
		return nil, false
	}

	if response.StatusCode == http.StatusOK {
		return securityURL, true
	}

	return nil, false
}
