package statuscake

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type jsonNumberString string

func (v *jsonNumberString) UnmarshalJSON(b []byte) error {
	if err := v.unmarshalAsInt(b); err == nil {
		return nil
	}
	if err := v.unmarshalAsString(b); err == nil {
		return nil
	}
	return fmt.Errorf("cannot unmarshal value that is neither a number nor a string: %s", truncate(b, 30))
}

func (v *jsonNumberString) unmarshalAsInt(b []byte) error {
	if bytes.Equal(b, []byte(`null`)) {
		return errors.New("cannot unmarshal JSON null to Go int")
	}

	var vv int
	if err := json.Unmarshal(b, &vv); err != nil {
		return err
	}
	*v = jsonNumberString(strconv.Itoa(vv))
	return nil
}

func (v *jsonNumberString) unmarshalAsString(b []byte) error {
	var vv string
	if err := json.Unmarshal(b, &vv); err != nil {
		return err
	}
	*v = jsonNumberString(vv)
	return nil
}

const truncateEllipses = "..."

func truncate(b []byte, max int) []byte {
	lte := len(truncateEllipses)
	min := lte + 3
	if max < min {
		max = min
	}
	if len(b) > max {
		t := make([]byte, max)
		n := max - lte
		copy(t, b[0:n])
		copy(t[n:], truncateEllipses)
		return t
	}
	return b
}
