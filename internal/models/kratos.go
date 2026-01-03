package models

// KratosSession represents a Kratos session
type KratosSession struct {
	ID       string                 `json:"id"`
	Active   bool                   `json:"active"`
	Identity KratosIdentity         `json:"identity"`
	Traits   map[string]interface{} `json:"traits"`
}

// KratosIdentity represents a Kratos identity
type KratosIdentity struct {
	ID     string                 `json:"id"`
	Traits map[string]interface{} `json:"traits"`
}
