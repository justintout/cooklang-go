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
	if q.Units == "" {
		return q.S
	}
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

// unitFamilies maps a normalized unit to its family name and its size in
// that family's base unit. Only units in the same family can be summed. The
// families cover the common cooking cases: us volume (tsp/tbsp/cup), metric
// mass (g/kg), metric volume (ml/l) and us mass (oz/lb).
var unitFamilies = map[string]struct {
	family string
	factor float64
}{
	"tsp":         {"volume-us", 1},
	"teaspoon":    {"volume-us", 1},
	"teaspoons":   {"volume-us", 1},
	"tbsp":        {"volume-us", 3},
	"tablespoon":  {"volume-us", 3},
	"tablespoons": {"volume-us", 3},
	"cup":         {"volume-us", 48},
	"cups":        {"volume-us", 48},
	"g":           {"mass-metric", 1},
	"gram":        {"mass-metric", 1},
	"grams":       {"mass-metric", 1},
	"kg":          {"mass-metric", 1000},
	"kilogram":    {"mass-metric", 1000},
	"kilograms":   {"mass-metric", 1000},
	"ml":          {"volume-metric", 1},
	"milliliter":  {"volume-metric", 1},
	"milliliters": {"volume-metric", 1},
	"l":           {"volume-metric", 1000},
	"liter":       {"volume-metric", 1000},
	"liters":      {"volume-metric", 1000},
	"oz":          {"mass-us", 1},
	"ounce":       {"mass-us", 1},
	"ounces":      {"mass-us", 1},
	"lb":          {"mass-us", 16},
	"pound":       {"mass-us", 16},
	"pounds":      {"mass-us", 16},
}

// unitKey returns a key identifying which quantities may be summed together.
// Quantities with a known unit share a key with everything in the same family;
// otherwise the (lowercased) unit string must match exactly.
func unitKey(units string) string {
	u := strings.ToLower(strings.TrimSpace(units))
	if f, ok := unitFamilies[u]; ok {
		return "family:" + f.family
	}
	return "unit:" + u
}

// sumQuantities collapses a list of quantities for a single ingredient into a
// display list. Numeric quantities whose units are identical or convertible
// within a common family are summed into one total per group, preserving first
// appearance order; non-numeric quantities (e.g. "some", "a few") are kept as
// their own entries.
func sumQuantities(qs []Quantity) []string {
	type group struct {
		total float64
		units string // display unit: the units of the first member
	}
	order := make([]string, 0, len(qs))
	groups := make(map[string]*group)
	out := make([]string, 0, len(qs))

	for _, q := range qs {
		if q.N < 0 {
			// non-numeric or "some" quantity: keep it verbatim
			out = append(out, q.String())
			continue
		}
		key := unitKey(q.Units)
		g, ok := groups[key]
		if !ok {
			g = &group{units: q.Units}
			groups[key] = g
			order = append(order, key)
		}
		if f, isFamily := unitFamilies[strings.ToLower(strings.TrimSpace(q.Units))]; isFamily {
			if df, ok := unitFamilies[strings.ToLower(strings.TrimSpace(g.units))]; ok {
				g.total += float64(q.N) * f.factor / df.factor
				continue
			}
		}
		g.total += float64(q.N)
	}

	for _, key := range order {
		g := groups[key]
		q := Quantity{N: float32(g.total), S: formatQuantity(g.total), Units: g.units}
		out = append(out, q.String())
	}
	if len(out) == 0 {
		return []string{}
	}
	return out
}

// formatQuantity renders a summed amount without trailing zeros, matching the
// compact style quantities are written in.
func formatQuantity(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
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
