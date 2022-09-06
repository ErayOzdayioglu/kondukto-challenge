package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type ScanResults struct {
	Id                 primitive.ObjectID `json:"_id,omitempty"`
	GithubUrl          string             `json:"githubUrl"`
	IssuesByConfidence IssuesByConfidence `json:"issuesByConfidence"`
	IssuesBySeverity   IssuesBySeverity   `json:"issuesBySeverity"`
}
