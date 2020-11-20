package ticker

import (
	"fmt"

	"gobunny/store/model"
)

// Key returns a valid model.Key value for a Ticker model
func Key(symbol string) (model.Key, error) {
	return model.Key(fmt.Sprintf("ticker:%s", symbol)), nil
}
