package cooklang

import "strings"

type Metadata map[string]string

func (m Metadata) Add(input string) {
	input = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(input), ">>"))
	s := strings.SplitN(input, ":", 2)
	if len(s) == 1 {
		s = append(s, "")
	}
	s[0], s[1] = strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
	m[s[0]] = s[1]
}
