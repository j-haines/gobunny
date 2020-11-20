package model

type (
	// Key is a type-enforced string value appropriate for use as an index value
	Key string

	// Value is a type-enforced string value appropriate for decoding into a Model object
	Value string

	// Marshaler encodes a Model into its string representation
	Marshaler interface {
		Marshal() (string, error)
	}

	// Unmarshaler decodes a Model from a given string
	Unmarshaler interface {
		Unmarshal(string) error
	}

	// Model types are Store-backed data structures
	Model interface {
		Marshaler
		Unmarshaler

		Key() Key
		Value() Value

		Validate() error
	}
)

func (k Key) String() string {
	return string(k)
}

func (v Value) String() string {
	return string(v)
}
