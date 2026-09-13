package ai

// SafeAction represents an allowlisted navigation action that the reviewer assistant may suggest.
type SafeAction struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Target string `json:"target"` // internal route or anchor
}

// safeActionRegistry defines all allowlisted actions known to the server.
var safeActionRegistry = map[string]SafeAction{
	"view-project-edutrace": {
		ID:     "view-project-edutrace",
		Label:  "View EduTrace Case Study",
		Target: "/projects/edutrace",
	},
	"view-project-gdgoc-ecommerce": {
		ID:     "view-project-gdgoc-ecommerce",
		Label:  "View GDGOC E-Commerce Case Study",
		Target: "/projects/gdgoc-ecommerce",
	},
	"view-project-maritime-ai-dashboard": {
		ID:     "view-project-maritime-ai-dashboard",
		Label:  "View Maritime AI Case Study",
		Target: "/projects/maritime-ai-dashboard",
	},
	"view-project-smart-kitchen": {
		ID:     "view-project-smart-kitchen",
		Label:  "View Smart Kitchen Case Study",
		Target: "/projects/smart-kitchen",
	},
	"go-to-projects": {
		ID:     "go-to-projects",
		Label:  "Explore Projects Section",
		Target: "/#projects",
	},
	"go-to-skills": {
		ID:     "go-to-skills",
		Label:  "View Skills & Evidence",
		Target: "/#skills",
	},
	"go-to-journey": {
		ID:     "go-to-journey",
		Label:  "View Academic Journey",
		Target: "/#journey",
	},
	"go-to-contact": {
		ID:     "go-to-contact",
		Label:  "Contact Arham",
		Target: "/#contact",
	},
}

// GetSafeAction returns a copy of the safe action if registered, or false.
func GetSafeAction(id string) (SafeAction, bool) {
	action, exists := safeActionRegistry[id]
	return action, exists
}

// IsValidAction checks whether an action ID is allowlisted.
func IsValidAction(id string) bool {
	_, exists := safeActionRegistry[id]
	return exists
}

// AllSafeActions returns a list of all registered safe actions.
func AllSafeActions() []SafeAction {
	actions := make([]SafeAction, 0, len(safeActionRegistry))
	for _, a := range safeActionRegistry {
		actions = append(actions, a)
	}
	return actions
}
