package model

type Organization struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Issue struct {
	ID                        string            `json:"id"`
	Name                      string            `json:"name"`
	ObjectivesAndResearchArea string            `json:"objectivesAndResearchArea"`
	Outcomes                  []IssueOutcome    `json:"outcomes"`
	Elements                  []IssueElement    `json:"elements"`
	Roles                     []IssueRole       `json:"roles"`
	Attachments               []IssueAttachment `json:"attachments"`
}

type IssueOutcome struct {
	ID      string `json:"id"`
	Outcome string `json:"outcome"`
	IssueID string `json:"-"`
}

type IssueElement struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrganizationID string `json:"organizationId"`
	ParentID       string `json:"parentId"`
	IssueID        string `json:"-"`
}

type IssueRole struct {
	ID      string   `json:"id"`
	UserID  string   `json:"userId"`
	Roles   []string `json:"roles"`
	IssueID string   `json:"-"`
}

type IssueAttachment struct {
	ID         string `json:"id"`
	Attachment string `json:"attachment"`
	IssueID    string `json:"-"`
}

type IssueRoleEntity struct {
	ID      string `json:"id"`
	UserID  string `json:"userId"`
	Roles   string `json:"roles"`
	IssueID string `json:"-"`
}

type Incident struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TODO: refactor with composition
type ExtendedIncident struct {
	State         string `json:"state"`
	ClosedAt      string `json:"closedAt"`
	FirstObserved string `json:"firstObserved"`
	ID            string `json:"id"`
	Type          string `json:"type"`
	Group         string `json:"group"`
	AssignedTo    string `json:"assignedTo"`
	Where         string `json:"where"`
	Name          string `json:"name"`
	Description   string `json:"description"`
}

type Policy struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Story struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
