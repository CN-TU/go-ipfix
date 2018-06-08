package ipfix

// StructuredSemantic represents the semantic of the structured data type according to RFC6313
type StructuredSemantic byte

const (
	NoneOfSemantic StructuredSemantic = iota
	ExactlyOneOfSemantic
	OneOrMoreOfSemantic
	AllOfSemantic
	OrderedSemantic
	UndefinedSemantic StructuredSemantic = 0xFF
)
