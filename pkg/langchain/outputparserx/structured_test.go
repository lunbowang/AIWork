/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package outputparserx

import "testing"

func Test_Structured_GetFormatInstructions(t *testing.T) {
	out := NewStructured([]ResponseSchema{
		{
			Name:        "title",
			Description: "this is title ",
		}, {
			Name:        "deadlineAt",
			Description: "todo deadline",
			Type:        "int64",
		}, {
			Name:        "executeIds",
			Description: "todo execute ids",
			Type:        "[]string",
		}, {
			Name:        "record",
			Description: "record todo handler",
			Schemas: []ResponseSchema{
				{
					Name:        "FinishAt",
					Description: "todo finish time",
					Type:        "int64",
				}, {
					Name:        "Content",
					Description: "record content",
				},
			},
		},
	})

	t.Log(out.GetFormatInstructions())
}
