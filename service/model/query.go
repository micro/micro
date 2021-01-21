package model

type Query struct {
	Index
	Order  Order
	Value  interface{}
	Offset int64
	Limit  int64
}

// Equals is an equality query by `fieldName`
// It filters records where `fieldName` equals to a value.
func Equals(fieldName string, value interface{}) Query {
	return Query{
		Index: Index{
			Type:      queryTypeEq,
			FieldName: fieldName,
			Order: Order{
				FieldName: fieldName,
				Type:      OrderTypeAsc,
			},
		},
		Value: value,
		Order: Order{
			FieldName: fieldName,
			Type:      OrderTypeAsc,
		},
	}
}

// QueryByID is short hand for querying by the primary index
func QueryByID(id string) Query {
	return Equals("ID", id)
}
