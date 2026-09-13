package dto

// DataEnvelope wraps a single resource response.
type DataEnvelope[T any] struct {
	Data T `json:"data"`
}

// Meta provides metadata for collection responses.
type Meta struct {
	Count int `json:"count"`
}

// ListEnvelope wraps a collection response.
type ListEnvelope[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

// NewDataEnvelope creates a DataEnvelope for a single item.
func NewDataEnvelope[T any](data T) DataEnvelope[T] {
	return DataEnvelope[T]{Data: data}
}

// NewListEnvelope creates a ListEnvelope for a collection.
func NewListEnvelope[T any](items []T) ListEnvelope[T] {
	if items == nil {
		items = []T{}
	}
	return ListEnvelope[T]{
		Data: items,
		Meta: Meta{Count: len(items)},
	}
}
