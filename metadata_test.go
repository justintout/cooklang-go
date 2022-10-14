package cooklang

import "testing"

func TestMetadata(t *testing.T) {
	t.Run("Metadata#Add", func(t *testing.T) {
		m := Metadata{}
		m.Add(">> one: thing")
		if v, ok := m["one"]; !ok || v != "thing" {
			t.Errorf(`expected key/value "one: thing", got: %+v`, m)
		}
		m.Add("  >> another:      different thing     ")

		one, ok := m["one"]
		if !ok {
			t.Errorf(`key/value "one: thing" removed by Add call: %+v`, m)
		}
		if one != "thing" {
			t.Errorf(`key/value "one: thing" changed by Add call to different key: %+v`, m)
		}

		if v, ok := m["another"]; !ok || v != "different thing" {
			t.Errorf(`expected new key/value "another: different thing", got: %+v`, m)
		}

		m.Add(">> one: actually this thing")
		if v, ok := m["one"]; !ok || v != "actually this thing" {
			t.Errorf(`expected updated ey/value "one: actually this thing", got: %+v`, m)
		}
	})
}
