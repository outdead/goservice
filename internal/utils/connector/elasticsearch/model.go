package elasticsearch

// Model is an interface for Elasticsearch data structures.
// Describes methods for getting the table name and unique identifier.
type Model interface {
	// TableName returns the table name.
	TableName() string

	// CalculateID generates and returns unique identifier from structure data.
	CalculateID() string
}
