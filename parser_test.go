package cooklang_test

import (
	"reflect"
	"testing"

	"github.com/justintout/cooklang-go"
)

func TestSections(t *testing.T) {
	src := `= Dough
Mix @flour{200%g} and @water{100%ml}.

== Filling ==
Slice @apples{4}.

=== Assembly ===
Combine dough and filling.

==
Serve warm.
`
	r := cooklang.MustParse(src)
	if len(r.Sections) != 4 {
		t.Fatalf("got %d sections, want 4", len(r.Sections))
	}
	wantNames := []string{"Dough", "Filling", "Assembly", ""}
	for i, want := range wantNames {
		if r.Sections[i].Name != want {
			t.Errorf("section[%d].Name = %q, want %q", i, r.Sections[i].Name, want)
		}
		if len(r.Sections[i].Steps) != 1 {
			t.Errorf("section[%d] step count = %d, want 1", i, len(r.Sections[i].Steps))
		}
	}
	if len(r.Steps) != 4 {
		t.Fatalf("got %d steps, want 4", len(r.Steps))
	}
	for i, want := range wantNames {
		if r.Steps[i].Section != want {
			t.Errorf("step[%d].Section = %q, want %q", i, r.Steps[i].Section, want)
		}
		if r.Steps[i].Number != i {
			t.Errorf("step[%d].Number = %d, want %d", i, r.Steps[i].Number, i)
		}
	}
}

func TestNotes(t *testing.T) {
	t.Run("single-line and multi-line notes", func(t *testing.T) {
		src := `> A single-line background note.

= Dough
> Keep the butter cold
and work quickly.

> First blockquote line.
> Second blockquote line.

Mix @flour{200%g} and @butter{100%g}.
`
		r := cooklang.MustParse(src)
		if len(r.Notes) != 3 {
			t.Fatalf("got %d notes, want 3: %+v", len(r.Notes), r.Notes)
		}
		wantNotes := []struct {
			value   string
			section string
		}{
			{"A single-line background note.", ""},
			{"Keep the butter cold and work quickly.", "Dough"},
			{"First blockquote line. Second blockquote line.", "Dough"},
		}
		for i, want := range wantNotes {
			if r.Notes[i].Value != want.value {
				t.Errorf("note[%d].Value = %q, want %q", i, r.Notes[i].Value, want.value)
			}
			if r.Notes[i].Section != want.section {
				t.Errorf("note[%d].Section = %q, want %q", i, r.Notes[i].Section, want.section)
			}
		}
		if len(r.Steps) != 1 {
			t.Fatalf("got %d steps, want 1", len(r.Steps))
		}
	})
}

func TestForcedLineBreaks(t *testing.T) {
	src := "Mix the dry ingredients:\\\n- @flour{200%g}\\\n- @sugar{50%g}\nand stir well.\n"
	r := cooklang.MustParse(src)
	if len(r.Steps) != 1 {
		t.Fatalf("got %d steps, want 1", len(r.Steps))
	}
	step := r.Steps[0]
	want := []cooklang.DirectionItem{
		{Type: "text", Value: "Mix the dry ingredients:\n- "},
		{Type: "ingredient", Name: "flour", Quantity: "200", Units: "g"},
		{Type: "text", Value: "\n- "},
		{Type: "ingredient", Name: "sugar", Quantity: "50", Units: "g"},
		{Type: "text", Value: " and stir well."},
	}
	got := make([]cooklang.DirectionItem, 0, len(step.DirectionItems))
	for _, d := range step.DirectionItems {
		got = append(got, d.DirectionItem())
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got items:\n%+v\nwant:\n%+v", got, want)
	}
}

func TestOptionalAndPreparationInSteps(t *testing.T) {
	src := "Add @onion{1}(peeled and finely chopped), @garlic{}(minced), @?chives, and @?chilli flakes{1%pinch} into a #?splatter guard{}. Do not parse @? spaced or #? spaced as items.\n"
	r := cooklang.MustParse(src)
	if len(r.Steps) != 1 {
		t.Fatalf("got %d steps, want 1", len(r.Steps))
	}
	s := r.Steps[0]
	if len(s.Ingredients) != 4 {
		t.Fatalf("got %d ingredients, want 4: %+v", len(s.Ingredients), s.Ingredients)
	}
	if s.Ingredients[0].Name != "onion" || s.Ingredients[0].Preparation != "peeled and finely chopped" || s.Ingredients[0].Optional {
		t.Errorf("unexpected first ingredient: %+v", s.Ingredients[0])
	}
	if s.Ingredients[1].Name != "garlic" || s.Ingredients[1].Preparation != "minced" || s.Ingredients[1].Optional {
		t.Errorf("unexpected second ingredient: %+v", s.Ingredients[1])
	}
	if s.Ingredients[2].Name != "chives" || !s.Ingredients[2].Optional {
		t.Errorf("unexpected third ingredient: %+v", s.Ingredients[2])
	}
	if s.Ingredients[3].Name != "chilli flakes" || !s.Ingredients[3].Optional || s.Ingredients[3].Units != "pinch" {
		t.Errorf("unexpected fourth ingredient: %+v", s.Ingredients[3])
	}
	if len(s.Cookware) != 1 {
		t.Fatalf("got %d cookware, want 1: %+v", len(s.Cookware), s.Cookware)
	}
	if s.Cookware[0].Name != "splatter guard" || !s.Cookware[0].Optional {
		t.Errorf("unexpected cookware: %+v", s.Cookware[0])
	}
	lastText := s.Text[len(s.Text)-1].Value
	if lastText != ". Do not parse @? spaced or #? spaced as items." {
		t.Errorf("unexpected trailing text: %q", lastText)
	}
}

func TestRecipeReferences(t *testing.T) {
	src := "Pour @./sauces/hollandaise{150%g} and @?./garnish/herbs{} over @./base.\n"
	r := cooklang.MustParse(src)
	if len(r.Steps) != 1 {
		t.Fatalf("got %d steps, want 1", len(r.Steps))
	}
	ings := r.Steps[0].Ingredients
	if len(ings) != 3 {
		t.Fatalf("got %d ingredients, want 3: %+v", len(ings), ings)
	}
	if ings[0].Name != "./sauces/hollandaise" || !ings[0].IsReference() || ings[0].Optional {
		t.Errorf("unexpected ref 0: %+v", ings[0])
	}
	if ings[1].Name != "./garnish/herbs" || !ings[1].IsReference() || !ings[1].Optional {
		t.Errorf("unexpected ref 1: %+v", ings[1])
	}
	if ings[2].Name != "./base" || !ings[2].IsReference() || ings[2].Optional {
		t.Errorf("unexpected ref 2: %+v", ings[2])
	}
}
