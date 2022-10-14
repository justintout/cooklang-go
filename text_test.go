package cooklang

import "testing"

func TestText(t *testing.T) {
	t.Run("NewText", func(t *testing.T) {
		txt := NewText("asdf")
		if txt.Type != "text" || txt.Value != "asdf" {
			t.Errorf("expected new Text w/ value 'asdf', got: %v", t)
		}
	})
}
