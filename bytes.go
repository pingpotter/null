package null

import (
	"database/sql/driver"

	"gopkg.in/nullbio/null.v4/convert"
)

// NullBytes is a nullable byte slice.
type NullBytes struct {
	Bytes []byte
	Valid bool
}

// Bytes is a nullable []byte.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Bytes struct {
	NullBytes
}

// NewBytes creates a new Bytes
func NewBytes(b []byte, valid bool) Bytes {
	return Bytes{
		NullBytes: NullBytes{
			Bytes: b,
			Valid: valid,
		},
	}
}

// BytesFrom creates a new Bytes that will be null if len zero.
func BytesFrom(b []byte) Bytes {
	return NewBytes(b, len(b) != 0)
}

// BytesFromPtr creates a new Bytes that be null if len zero.
func BytesFromPtr(b *[]byte) Bytes {
	if b == nil || len(*b) == 0 {
		return NewBytes(nil, false)
	}
	n := NewBytes(*b, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// Bytes UnmarshalJSON is different in that it only
// unmarshals sql.NullBytes defined as JSON objects,
// It supports all JSON types.
// It also supports unmarshalling a sql.NullBytes.
func (b *Bytes) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		b.Bytes = nil
		b.Valid = false
	} else {
		b.Bytes = append(b.Bytes[0:0], data...)
		b.Valid = true
	}

	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Bytes if the input is blank.
// It will return an error if the input is not an integer, blank, or "null".
func (b *Bytes) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		b.Valid = false
	} else {
		b.Bytes = append(b.Bytes[0:0], text...)
		b.Valid = true
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if the Bytes is invalid.
func (b Bytes) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	return b.Bytes, nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode nil if the Bytes is invalid.
func (b Bytes) MarshalText() ([]byte, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bytes, nil
}

// SetValid changes this Bytes's value and also sets it to be non-null.
func (b *Bytes) SetValid(n []byte) {
	b.Bytes = n
	b.Valid = true
}

// Ptr returns a pointer to this Bytes's value, or a nil pointer if this Bytes is null.
func (b Bytes) Ptr() *[]byte {
	if !b.Valid {
		return nil
	}
	return &b.Bytes
}

// IsZero returns true for null or zero Bytes's, for future omitempty support (Go 1.4?)
func (b Bytes) IsZero() bool {
	return !b.Valid
}

// Scan implements the Scanner interface.
func (n *NullBytes) Scan(value interface{}) error {
	if value == nil {
		n.Bytes, n.Valid = []byte{}, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Bytes, value)
}

// Value implements the driver Valuer interface.
func (n NullBytes) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bytes, nil
}
