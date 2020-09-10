package sprintf

//
// import (
// 	"github.com/jexia/semaphore/pkg/specs"
// 	"github.com/jexia/semaphore/pkg/specs/labels"
// 	"github.com/jexia/semaphore/pkg/specs/template"
// 	"github.com/jexia/semaphore/pkg/specs/types"
// )
//
// var (
// 	schema = &specs.ParameterMap{
// 		Property: &specs.Property{
// 			Type:  types.Message,
// 			Label: labels.Optional,
// 			Nested: map[string]*specs.Property{
// 				"numeric": {
// 					Name:  "bad_label",
// 					Path:  "bad_label",
// 					Type:  types.Int32,
// 					Label: labels.Optional,
// 				},
// 				"message": {
// 					Name:  "message",
// 					Path:  "message",
// 					Type:  types.String,
// 					Label: labels.Optional,
// 					Reference: &specs.PropertyReference{
// 						Resource: template.InputResource,
// 						Path:     "message",
// 					},
// 				},
// 				"status": {
// 					Name:  "status",
// 					Path:  "status",
// 					Type:  types.Enum,
// 					Label: labels.Optional,
// 					Enum:  enum,
// 					Reference: &specs.PropertyReference{
// 						Resource: template.InputResource,
// 						Path:     "status",
// 					},
// 				},
// 				"nested": {
// 					Name:  "nested",
// 					Path:  "nested",
// 					Type:  types.Message,
// 					Label: labels.Optional,
// 					Nested: map[string]*specs.Property{
// 						"first": {
// 							Name:  "first",
// 							Path:  "nested.first",
// 							Type:  types.String,
// 							Label: labels.Optional,
// 							Reference: &specs.PropertyReference{
// 								Resource: template.InputResource,
// 								Path:     "nested.first",
// 							},
// 						},
// 						"second": {
// 							Name:  "second",
// 							Path:  "nested.second",
// 							Type:  types.String,
// 							Label: labels.Optional,
// 							Reference: &specs.PropertyReference{
// 								Resource: template.InputResource,
// 								Path:     "nested.second",
// 							},
// 						},
// 					},
// 				},
// 				"repeating_string": {
// 					Name:  "repeating_string",
// 					Path:  "repeating_string",
// 					Type:  types.String,
// 					Label: labels.Repeated,
// 					Reference: &specs.PropertyReference{
// 						Resource: template.InputResource,
// 						Path:     "repeating_string",
// 					},
// 				},
// 				"repeating_enum": {
// 					Name:  "repeating_enum",
// 					Path:  "repeating_enum",
// 					Type:  types.Enum,
// 					Label: labels.Repeated,
// 					Enum:  enum,
// 					Reference: &specs.PropertyReference{
// 						Resource: template.InputResource,
// 						Path:     "repeating_enum",
// 					},
// 				},
// 				"repeating": {
// 					Name:  "repeating",
// 					Path:  "repeating",
// 					Type:  types.Message,
// 					Label: labels.Repeated,
// 					Nested: map[string]*specs.Property{
// 						"value": {
// 							Name:  "value",
// 							Path:  "repeating.value",
// 							Type:  types.String,
// 							Label: labels.Optional,
// 							Reference: &specs.PropertyReference{
// 								Resource: template.InputResource,
// 								Path:     "repeating.value",
// 							},
// 						},
// 					},
// 					Reference: &specs.PropertyReference{
// 						Resource: template.InputResource,
// 						Path:     "repeating",
// 					},
// 				},
// 			},
// 		},
// 	}
// )
