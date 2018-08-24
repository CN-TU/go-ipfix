package ipfix

// StructuredSemantic represents the semantic of the structured data type according to RFC6313
type StructuredSemantic byte

const (
	// NoneOfSemantic as defined by RFC6313
	NoneOfSemantic StructuredSemantic = iota
	// ExactlyOneOfSemantic as defined by RFC6313
	ExactlyOneOfSemantic
	// OneOrMoreOfSemantic as defined by RFC6313
	OneOrMoreOfSemantic
	// AllOfSemantic as defined by RFC6313
	AllOfSemantic
	// OrderedSemantic as defined by RFC6313
	OrderedSemantic
	// UndefinedSemantic as defined by RFC6313
	UndefinedSemantic StructuredSemantic = 0xFF
)
