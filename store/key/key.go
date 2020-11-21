package key

import "gobunny/store/model"

// Must panics if given a non-nil error, else the given model.Key
func Must(key model.Key, err error) model.Key {
	if err != nil {
		panic(err)
	}

	return key
}
