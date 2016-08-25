package strategy

import validator "gopkg.in/bluesuncorp/validator.v8"

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// Field is the information related with the transformation for a specific field
// inside the entity.
// Foreign converts into a Local through the specific transformation.
// Type is the data type
type Field struct {
	Foreign        string `bson:"foreign" json:"foreign"`
	Local          string `bson:"local" json:"local"`
	Transformation string `bson:"transformation" json:"transformation"`
	Type           string `bson:"type" json:"type"`
	DateTimeFormat string `jbson:"datetimeformat" json:"datetimeformat"`
}

// Validate validates a field value with the validator.
func (field *Field) Validate() error {
	if err := validate.Struct(field); err != nil {
		return err
	}

	return nil
}

// Entity is the information related to transform one concept from the external source
// into a type of Item in the Coral systems
type Entity struct {
	Foreign        string            `json:"foreign"`
	Local          string            `json:"local"`
	OrderBy        string            `json:"orderby"`
	Fields         []Field           `json:"fields"`
	DateTimeFormat string            `json:"datetimeformat"`
	Status         map[string]string `json:"status"`
}

// Validate validates an Entity value with the validator.
func (entity *Entity) Validate() error {
	if err := validate.Struct(entity); err != nil {
		return err
	}

	return nil
}

// Strategy explains which entities or data we are getting from the source
// and which transformation nees to happen.
type Strategy struct {
	Name           string            `bson:"name" json:"name"`
	DateTimeFormat string            `bson:"datetimeformat" json:"datetimeformat"`
	Entities       map[string]Entity `bson:"entities" json:"entities"`
}

// Validate validates a Strategy value with the validator.
func (strategy *Strategy) Validate() error {
	if err := validate.Struct(strategy); err != nil {
		return err
	}

	return nil
}
