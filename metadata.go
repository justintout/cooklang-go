package cooklang

import (
	"strings"

	"go.yaml.in/yaml/v3"
)

// Metadata represents recipe metadata
type Metadata map[string]string

// Source represents structured recipe source metadata defined by canonical
// Cooklang metadata conventions.
type Source struct {
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Author string `json:"author,omitempty" yaml:"author,omitempty"`
	URL    string `json:"url,omitempty" yaml:"url,omitempty"`
}

// Add inserts one or more metadata key/value entries from a front matter block
// or a version 5 metadata line. It supports YAML scalars, sequences (stored as
// comma-separated values and accessible via List/Tags) and nested mappings
// (flattened with dot-separated keys such as "source.name").
func (m Metadata) Add(input string) {
	if strings.TrimSpace(input) == "" {
		return
	}
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(input), &doc); err == nil && len(doc.Content) > 0 && doc.Content[0].Kind == yaml.MappingNode {
		addMappingNode(m, "", doc.Content[0])
		return
	}
	for _, line := range strings.Split(input, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		addLine(m, line)
	}
}

func addMappingNode(m Metadata, prefix string, node *yaml.Node) {
	for i := 0; i+1 < len(node.Content); i += 2 {
		kNode := node.Content[i]
		vNode := node.Content[i+1]
		key := strings.TrimSpace(kNode.Value)
		if prefix != "" {
			key = prefix + "." + key
		}
		switch vNode.Kind {
		case yaml.MappingNode:
			addMappingNode(m, key, vNode)
		case yaml.SequenceNode:
			items := make([]string, 0, len(vNode.Content))
			for _, item := range vNode.Content {
				if item.Kind == yaml.ScalarNode {
					items = append(items, strings.TrimSpace(item.Value))
				}
			}
			m[key] = strings.Join(items, ", ")
		case yaml.ScalarNode:
			if vNode.Tag == "!!null" && vNode.Value == "" && strings.Contains(kNode.Value, ":") {
				addLineWithPrefix(m, prefix, kNode.Value)
				continue
			}
			m[key] = strings.TrimSpace(vNode.Value)
		}
	}
}

func addLine(m Metadata, line string) {
	addLineWithPrefix(m, "", line)
}

func addLineWithPrefix(m Metadata, prefix, line string) {
	line = strings.TrimSpace(line)
	s := strings.SplitN(line, ":", 2)
	if len(s) == 1 {
		s = append(s, "")
	}
	key, val := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
	if prefix != "" {
		key = prefix + "." + key
	}
	m[key] = val
}

// Get returns the string value for key.
func (m Metadata) Get(key string) string {
	return m[key]
}

// List returns the slice of values stored under key, splitting comma-separated
// entries produced by YAML sequences or inline comma lists.
func (m Metadata) List(key string) []string {
	raw, ok := m[key]
	if !ok || strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

// Tags returns the recipe's tags as a slice of strings.
func (m Metadata) Tags() []string {
	return m.List("tags")
}

// Map returns all nested entries under prefix (such as "source") with the
// prefix and dot removed.
func (m Metadata) Map(prefix string) map[string]string {
	out := make(map[string]string)
	p := prefix + "."
	for k, v := range m {
		if sub, ok := strings.CutPrefix(k, p); ok {
			out[sub] = v
		}
	}
	return out
}

// Source returns the structured source metadata from "source.name",
// "source.author", "source.url", or a scalar "source" entry.
func (m Metadata) Source() Source {
	src := Source{
		Name:   m["source.name"],
		Author: m["source.author"],
		URL:    m["source.url"],
	}
	if src.Name == "" && src.URL == "" {
		if s, ok := m["source"]; ok {
			if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
				src.URL = s
			} else {
				src.Name = s
			}
		}
	}
	return src
}
