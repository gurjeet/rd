package main

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	t "github.com/shivammg/parsers/types"
)

// Epsilon represents empty string.
const Epsilon = "ε"

func TestOne(test *testing.T) {
	/*
	   	"(", "id", "*", "id", ")", "+", "id"
	   Recursive Descent parser for the following grammar:
	   	E  -> TE'
	   	E' -> +TE'|ε
	   	T  -> FT'
	   	T' -> *FT'|ε
	   	F  -> id|(E)

	   - Used in parsing addition and multiplication arithmetic expressions.
	   - ε represents empty string.
	*/
	wantJSON := `{
		"Symbol": "E",
		"Children": [
		  {
			"Symbol": "T",
			"Children": [
			  {
				"Symbol": "F",
				"Children": [
				  {
					"Symbol": "(",
					"Children": null
				  },
				  {
					"Symbol": "E",
					"Children": [
					  {
						"Symbol": "T",
						"Children": [
						  {
							"Symbol": "F",
							"Children": [
							  {
								"Symbol": "id",
								"Children": null
							  }
							]
						  },
						  {
							"Symbol": "T'",
							"Children": [
							  {
								"Symbol": "*",
								"Children": null
							  },
							  {
								"Symbol": "F",
								"Children": [
								  {
									"Symbol": "id",
									"Children": null
								  }
								]
							  },
							  {
								"Symbol": "T'",
								"Children": [
								  {
									"Symbol": "ε",
									"Children": null
								  }
								]
							  }
							]
						  }
						]
					  },
					  {
						"Symbol": "E'",
						"Children": [
						  {
							"Symbol": "ε",
							"Children": null
						  }
						]
					  }
					]
				  },
				  {
					"Symbol": ")",
					"Children": null
				  }
				]
			  },
			  {
				"Symbol": "T'",
				"Children": [
				  {
					"Symbol": "ε",
					"Children": null
				  }
				]
			  }
			]
		  },
		  {
			"Symbol": "E'",
			"Children": [
			  {
				"Symbol": "+",
				"Children": null
			  },
			  {
				"Symbol": "T",
				"Children": [
				  {
					"Symbol": "F",
					"Children": [
					  {
						"Symbol": "id",
						"Children": null
					  }
					]
				  },
				  {
					"Symbol": "T'",
					"Children": [
					  {
						"Symbol": "ε",
						"Children": null
					  }
					]
				  }
				]
			  },
			  {
				"Symbol": "E'",
				"Children": [
				  {
					"Symbol": "ε",
					"Children": null
				  }
				]
			  }
			]
		  }
		]
	  }`
	var want, got interface{}
	json.Unmarshal([]byte(wantJSON), &want)
	p := NewParser([]string{"(", "id", "*", "id", ")", "+", "id"})

	/* shit */
	p.Register("E", func() (*t.Tree, error) {
		t1, err := p.Run("T")
		if err != nil {
			return nil, err
		}
		t2, err := p.Run("E'")
		if err != nil {
			return nil, err
		}
		return t.NewTree("E", t1, t2), nil
	})

	p.Register("E'", func() (*t.Tree, error) {
		if p.Match("+") {
			t1, err := p.Run("T")
			if err != nil {
				return nil, err
			}
			t2, err := p.Run("E'")
			if err != nil {
				return nil, err
			}
			return t.NewTree("E'", t.NewTree("+"), t1, t2), nil
		}
		// epsilon exists for the rule
		p.Backtrack()
		return t.NewTree("E'", t.NewTree(Epsilon)), nil
	})

	p.Register("T", func() (*t.Tree, error) {
		t1, err := p.Run("F")
		if err != nil {
			return nil, err
		}
		t2, err := p.Run("T'")
		if err != nil {
			return nil, err
		}
		return t.NewTree("T", t1, t2), nil
	})

	p.Register("T'", func() (*t.Tree, error) {
		if p.Match("*") {
			t1, err := p.Run("F")
			if err != nil {
				return nil, err
			}
			t2, err := p.Run("T'")
			if err != nil {
				return nil, err
			}
			return t.NewTree("T'", t.NewTree("*"), t1, t2), nil
		}
		// epsilon exists for the rule
		p.Backtrack()
		return t.NewTree("T'", t.NewTree(Epsilon)), nil
	})

	p.Register("F", func() (*t.Tree, error) {
		if p.Match("id") {
			return t.NewTree("F", t.NewTree("id")), nil
		}
		p.Backtrack()
		if p.Match("(") {
			t1, err := p.Run("E")
			if err != nil {
				return nil, err
			}
			if p.Match(")") {
				return t.NewTree("F", t.NewTree("("), t1, t.NewTree(")")), nil
			}
		}
		return nil, errors.New("No match")
	})

	/* shit */

	tree, err := p.Run("E")
	if err != nil {
		test.Fatal(err)
	}
	gotJSON, _ := json.Marshal(tree)
	json.Unmarshal(gotJSON, &got)
	if !reflect.DeepEqual(want, got) {
		test.Errorf("Expected: %v\nGot: %v\n", want, got)
	}
}