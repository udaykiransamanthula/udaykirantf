// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mocking

import (
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/internal/configs/configschema"
)

var (
	normalAttributes = map[string]*configschema.Attribute{
		"id": {
			Type: cty.String,
		},
		"value": {
			Type: cty.String,
		},
	}

	computedAttributes = map[string]*configschema.Attribute{
		"id": {
			Type:     cty.String,
			Computed: true,
		},
		"value": {
			Type: cty.String,
		},
	}

	normalBlock = configschema.Block{
		Attributes: normalAttributes,
	}

	computedBlock = configschema.Block{
		Attributes: computedAttributes,
	}
)

func TestComputedValuesForDataSource(t *testing.T) {
	tcs := map[string]struct {
		target           cty.Value
		with             cty.Value
		schema           *configschema.Block
		expected         cty.Value
		expectedFailures []string
	}{
		"nil_target_no_unknowns": {
			target: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("kj87eb9"),
				"value": cty.StringVal("Hello, world!"),
			}),
			with:   cty.NilVal,
			schema: &normalBlock,
			expected: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("kj87eb9"),
				"value": cty.StringVal("Hello, world!"),
			}),
		},
		"empty_target_no_unknowns": {
			target: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("kj87eb9"),
				"value": cty.StringVal("Hello, world!"),
			}),
			with:   cty.EmptyObjectVal,
			schema: &normalBlock,
			expected: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("kj87eb9"),
				"value": cty.StringVal("Hello, world!"),
			}),
		},
		"basic_computed_attribute_preset": {
			target: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("kj87eb9"),
				"value": cty.StringVal("Hello, world!"),
			}),
			with:   cty.NilVal,
			schema: &computedBlock,
			expected: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("kj87eb9"),
				"value": cty.StringVal("Hello, world!"),
			}),
		},
		"basic_computed_attribute_random": {
			target: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.NullVal(cty.String),
				"value": cty.StringVal("Hello, world!"),
			}),
			with:   cty.NilVal,
			schema: &computedBlock,
			expected: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("ssnk9qhr"),
				"value": cty.StringVal("Hello, world!"),
			}),
		},
		"basic_computed_attribute_supplied": {
			target: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.NullVal(cty.String),
				"value": cty.StringVal("Hello, world!"),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"id": cty.StringVal("myvalue"),
			}),
			schema: &computedBlock,
			expected: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("myvalue"),
				"value": cty.StringVal("Hello, world!"),
			}),
		},
		"nested_single_block_preset": {
			target: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ObjectVal(map[string]cty.Value{
					"id":    cty.NullVal(cty.String),
					"value": cty.StringVal("Hello, world!"),
				}),
			}),
			with: cty.NilVal,
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"block": {
						Block:   computedBlock,
						Nesting: configschema.NestingSingle,
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ObjectVal(map[string]cty.Value{
					"id":    cty.StringVal("ssnk9qhr"),
					"value": cty.StringVal("Hello, world!"),
				}),
			}),
		},
		"nested_single_block_supplied": {
			target: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ObjectVal(map[string]cty.Value{
					"id":    cty.NullVal(cty.String),
					"value": cty.StringVal("Hello, world!"),
				}),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ObjectVal(map[string]cty.Value{
					"id": cty.StringVal("myvalue"),
				}),
			}),
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"block": {
						Block:   computedBlock,
						Nesting: configschema.NestingSingle,
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ObjectVal(map[string]cty.Value{
					"id":    cty.StringVal("myvalue"),
					"value": cty.StringVal("Hello, world!"),
				}),
			}),
		},
		"nested_list_block_preset": {
			target: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("one"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
			with: cty.NilVal,
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"block": {
						Block:   computedBlock,
						Nesting: configschema.NestingList,
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("ssnk9qhr"),
						"value": cty.StringVal("one"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("amyllmyg"),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
		},
		"nested_list_block_supplied": {
			target: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("one"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ObjectVal(map[string]cty.Value{
					"id": cty.StringVal("myvalue"),
				}),
			}),
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"block": {
						Block:   computedBlock,
						Nesting: configschema.NestingList,
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("myvalue"),
						"value": cty.StringVal("one"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("myvalue"),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
		},
		"nested_set_block_preset": {
			target: cty.ObjectVal(map[string]cty.Value{
				"block": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("one"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
			with: cty.NilVal,
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"block": {
						Block:   computedBlock,
						Nesting: configschema.NestingSet,
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"block": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("ssnk9qhr"),
						"value": cty.StringVal("one"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("amyllmyg"),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
		},
		"nested_set_block_supplied": {
			target: cty.ObjectVal(map[string]cty.Value{
				"block": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("one"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ObjectVal(map[string]cty.Value{
					"id": cty.StringVal("myvalue"),
				}),
			}),
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"block": {
						Block:   computedBlock,
						Nesting: configschema.NestingSet,
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"block": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("myvalue"),
						"value": cty.StringVal("one"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("myvalue"),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
		},
		"nested_map_block_preset": {
			target: cty.ObjectVal(map[string]cty.Value{
				"block": cty.MapVal(map[string]cty.Value{
					"one": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("one"),
					}),
					"two": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
			with: cty.NilVal,
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"block": {
						Block:   computedBlock,
						Nesting: configschema.NestingMap,
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"block": cty.MapVal(map[string]cty.Value{
					"one": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("ssnk9qhr"),
						"value": cty.StringVal("one"),
					}),
					"two": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("amyllmyg"),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
		},
		"nested_map_block_supplied": {
			target: cty.ObjectVal(map[string]cty.Value{
				"block": cty.MapVal(map[string]cty.Value{
					"one": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("one"),
					}),
					"two": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ObjectVal(map[string]cty.Value{
					"id": cty.StringVal("myvalue"),
				}),
			}),
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"block": {
						Block:   computedBlock,
						Nesting: configschema.NestingMap,
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"block": cty.MapVal(map[string]cty.Value{
					"one": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("myvalue"),
						"value": cty.StringVal("one"),
					}),
					"two": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("myvalue"),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
		},
		"nested_single_attribute": {
			target: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.ObjectVal(map[string]cty.Value{
					"id":    cty.NullVal(cty.String),
					"value": cty.StringVal("Hello, world!"),
				}),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.ObjectVal(map[string]cty.Value{
					"id": cty.StringVal("myvalue"),
				}),
			}),
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"nested": {
						NestedType: &configschema.Object{
							Attributes: computedAttributes,
							Nesting:    configschema.NestingSingle,
						},
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.ObjectVal(map[string]cty.Value{
					"id":    cty.StringVal("myvalue"),
					"value": cty.StringVal("Hello, world!"),
				}),
			}),
		},
		"nested_list_attribute": {
			target: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("one"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.ObjectVal(map[string]cty.Value{
					"id": cty.StringVal("myvalue"),
				}),
			}),
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"nested": {
						NestedType: &configschema.Object{
							Attributes: computedAttributes,
							Nesting:    configschema.NestingList,
						},
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("myvalue"),
						"value": cty.StringVal("one"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("myvalue"),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
		},
		"nested_set_attribute": {
			target: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("one"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.ObjectVal(map[string]cty.Value{
					"id": cty.StringVal("myvalue"),
				}),
			}),
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"nested": {
						NestedType: &configschema.Object{
							Attributes: computedAttributes,
							Nesting:    configschema.NestingSet,
						},
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("myvalue"),
						"value": cty.StringVal("one"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("myvalue"),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
		},
		"nested_map_attribute": {
			target: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.MapVal(map[string]cty.Value{
					"one": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("one"),
					}),
					"two": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.ObjectVal(map[string]cty.Value{
					"id": cty.StringVal("myvalue"),
				}),
			}),
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"nested": {
						NestedType: &configschema.Object{
							Attributes: computedAttributes,
							Nesting:    configschema.NestingMap,
						},
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.MapVal(map[string]cty.Value{
					"one": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("myvalue"),
						"value": cty.StringVal("one"),
					}),
					"two": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("myvalue"),
						"value": cty.StringVal("two"),
					}),
				}),
			}),
		},
		"invalid_replacement_path": {
			target: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("kj87eb9"),
				"value": cty.StringVal("Hello, world!"),
			}),
			with:   cty.StringVal("Hello, world!"),
			schema: &normalBlock,
			expected: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("kj87eb9"),
				"value": cty.StringVal("Hello, world!"),
			}),
			expectedFailures: []string{
				"The requested replacement value must be an object type, but was string.",
			},
		},
		"invalid_replacement_path_nested": {
			target: cty.ObjectVal(map[string]cty.Value{
				"nested_object": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id": cty.NullVal(cty.String),
					}),
				}),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"nested_object": cty.StringVal("Hello, world!"),
			}),
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"nested_object": {
						NestedType: &configschema.Object{
							Attributes: map[string]*configschema.Attribute{
								"id": {
									Type:     cty.String,
									Computed: true,
								},
							},
							Nesting: configschema.NestingSet,
						},
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"nested_object": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id": cty.StringVal("ssnk9qhr"),
					}),
				}),
			}),
			expectedFailures: []string{
				"Terraform expected an object type at nested_object within the replacement value defined at :0,0-0, but found string.",
			},
		},
		"invalid_replacement_path_nested_block": {
			target: cty.ObjectVal(map[string]cty.Value{
				"nested_object": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id": cty.NullVal(cty.String),
					}),
				}),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"nested_object": cty.StringVal("Hello, world!"),
			}),
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"nested_object": {
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"id": {
									Type:     cty.String,
									Computed: true,
								},
							},
						},
						Nesting: configschema.NestingSet,
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"nested_object": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id": cty.StringVal("ssnk9qhr"),
					}),
				}),
			}),
			expectedFailures: []string{
				"Terraform expected an object type at nested_object within the replacement value defined at :0,0-0, but found string.",
			},
		},
		"invalid_replacement_type": {
			target: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.NullVal(cty.String),
				"value": cty.StringVal("Hello, world!"),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"id": cty.ListValEmpty(cty.String),
			}),
			schema: &computedBlock,
			expected: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("ssnk9qhr"),
				"value": cty.StringVal("Hello, world!"),
			}),
			expectedFailures: []string{
				"Terraform could not replace the target type string with the replacement value defined at id within :0,0-0: string required.",
			},
		},
		"invalid_replacement_type_nested": {
			target: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.MapVal(map[string]cty.Value{
					"one": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("one"),
					}),
				}),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.ObjectVal(map[string]cty.Value{
					"id": cty.EmptyObjectVal,
				}),
			}),
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"nested": {
						NestedType: &configschema.Object{
							Attributes: computedAttributes,
							Nesting:    configschema.NestingMap,
						},
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"nested": cty.MapVal(map[string]cty.Value{
					"one": cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("ssnk9qhr"),
						"value": cty.StringVal("one"),
					}),
				}),
			}),
			expectedFailures: []string{
				"Terraform could not replace the target type string with the replacement value defined at nested.id within :0,0-0: string required.",
			},
		},
		"invalid_replacement_type_nested_block": {
			target: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.NullVal(cty.String),
						"value": cty.StringVal("one"),
					}),
				}),
			}),
			with: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ObjectVal(map[string]cty.Value{
					"id": cty.EmptyObjectVal,
				}),
			}),
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"block": {
						Block:   computedBlock,
						Nesting: configschema.NestingList,
					},
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"block": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"id":    cty.StringVal("ssnk9qhr"),
						"value": cty.StringVal("one"),
					}),
				}),
			}),
			expectedFailures: []string{
				"Terraform could not replace the target type string with the replacement value defined at block.id within :0,0-0: string required.",
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {

			// We'll just make sure that any random strings are deterministic.
			testRand = rand.New(rand.NewSource(0))
			defer func() {
				testRand = nil
			}()

			actual, diags := ComputedValuesForDataSource(tc.target, ReplacementValue{
				Value: tc.with,
			}, tc.schema)

			var actualFailures []string
			for _, diag := range diags {
				actualFailures = append(actualFailures, diag.Description().Detail)
			}
			if diff := cmp.Diff(tc.expectedFailures, actualFailures); len(diff) > 0 {
				t.Errorf("unexpected failures\nexpected:\n%s\nactual:\n%s\ndiff:\n%s", tc.expectedFailures, actualFailures, diff)
			}

			if actual.Equals(tc.expected).False() {
				t.Errorf("\nexpected: (%s)\nactual:   (%s)", tc.expected.GoString(), actual.GoString())
			}
		})
	}
}
