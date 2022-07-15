package payloads

// CreateGroupRequest represents a request to create a new group
type CreateGroupRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Owners      []string `json:"owners,omitempty"`
}
