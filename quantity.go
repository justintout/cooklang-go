package cooklang

import (
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

var isFraction = regexp.MustCompile(`^([0-9]+\ )?[0-9]+/[0-9]$`)

type scalingType string

const (
	scalingNone   scalingType = "none"
	scalingLinear             = "linear"
	scalingManual             = "manual"
)

// Quantity is the representation of a quantity for ingredients and cookware,
// or a duration for a timer in Cooklang.
//
// TODO: scaling - it's kinda circular? AST needs to report scaling type then, if manual, the specific scaling portions?
type Quantity struct {
	N     float32
	S     string
	Units string
	// Scaling scalingType
	// ScaledQuantities []
	raw string
}

func (q Quantity) String() string {
	return fmt.Sprintf("%s %s", q.S, q.Units)
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
	q := Quantity{raw: source, S: defaultS, N: defaultN}
	if source == "" || source == "{}" {
		return q
	}

	s := strings.SplitN(strings.Trim(source, "{}"), dividerQuantity, 2)
	if len(s) > 1 {
		q.Units = s[1]
	}

	if s[0] == "" {
		return q
	}

	if v, err := strconv.ParseFloat(s[0], 32); err == nil {
		q.N = float32(v)
		q.S = s[0]
		return q
	}

	if isFraction.MatchString(s[0]) {
		/*
				_, .---.__c--.
			(__( )_._( )_`_>  lol ratatouille
					`~~"  `~"
		*/
		r := new(big.Rat)
		for _, ss := range strings.Split(s[0], " ") {
			rr := new(big.Rat)
			rr.SetString(ss)
			r.Add(r, rr)
		}
		q.N, _ = r.Float32()
		q.S = s[0]
		return q
	}

	q.S = s[0]
	return q
}

func strictParseQuantity(source string) (Quantity, error) {
	if source == "" || source == "{}" {
		return Quantity{}, fmt.Errorf("empty quantity not allowed in this context")
	}

	s := strings.SplitN(source[1:len(source)-2], "%", 2)
	if len(s) != 2 {
		return Quantity{}, fmt.Errorf("must have a quantity and unit in this context")
	}

	return parseQuantity(source, "", -1), nil
}
