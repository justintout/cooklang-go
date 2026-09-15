package cooklang

import (
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// fractions matches a fraction with an optional leading whole number, and
// tolerates spaces around the slash. The numerator and denominator may not
// start with a zero, which keeps "01/2" as written rather than reading it as
// one half.
var fractions = regexp.MustCompile(`^(?:(\d+)\s+)?([1-9]\d*)\s*/\s*([1-9]\d*)$`)

// Quantity is the representation of a quantity for ingredients and cookware,
// or a duration for a timer in Cooklang.
//
// TODO: scaling - it's kinda circular? AST needs to report scaling type then, if manual, the specific scaling portions?
type Quantity struct {
	N     float32
	S     string
	Units string
}

func (q Quantity) String() string {
	if q.Units == "" {
		return q.S
	}
	return fmt.Sprintf("%s %s", q.S, q.Units)
}

// Canonical renders the quantity the way the canonical test format does: as
// the numeric value when the quantity has one, and as written otherwise.
func (q Quantity) Canonical() string {
	if q.N < 0 {
		return q.S
	}
	return strconv.FormatFloat(float64(q.N), 'f', -1, 32)
}

func (q Quantity) MarshalJSON() ([]byte, error) {
	qq := struct {
		Quantity string `json:"quantity"`
		Units    string `json:"units"`
	}{
		Quantity: q.S,
		Units:    q.Units,
	}
	return json.Marshal(qq)
}

// TODO: add *Servings arg for scaling?
func parseQuantity(source string, defaultS string, defaultN float32) Quantity {
	q := Quantity{S: defaultS, N: defaultN}

	s := strings.SplitN(strings.Trim(source, "{}"), dividerQuantity, 2)
	if len(s) > 1 {
		q.Units = strings.TrimSpace(s[1])
	}

	amount := strings.TrimSpace(s[0])
	if amount == "" {
		return q
	}

	if v, err := strconv.ParseFloat(amount, 32); err == nil {
		q.N = float32(v)
		q.S = amount
		return q
	}

	if m := fractions.FindStringSubmatch(amount); m != nil {
		/*
				_, .---.__c--.
			(__( )_._( )_`_>  lol ratatouille
					`~~"  `~"
		*/
		r := new(big.Rat)
		if m[1] != "" {
			r.SetString(m[1])
		}
		f := new(big.Rat)
		f.SetString(m[2] + "/" + m[3])
		q.N, _ = r.Add(r, f).Float32()
		q.S = amount
		return q
	}

	// Not a number and not a fraction: keep it as written, and record that
	// there is no numeric value to fall back on.
	q.N, q.S = -1, amount
	return q
}

func strictParseQuantity(source string) (Quantity, error) {
	if source == "" || source == "{}" {
		return Quantity{}, fmt.Errorf("empty quantity not allowed in this context")
	}

	s := strings.SplitN(source[1:len(source)-1], "%", 2)
	if len(s) != 2 {
		return Quantity{}, fmt.Errorf("must have a quantity and unit in this context")
	}

	return parseQuantity(source, "", -1), nil
}
