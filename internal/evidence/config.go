package evidence

import "context"

type Config struct {
	context.Context
	Content           string `json:"content,omitempty"`
	Format            string `json:"format,omitempty"`
	GhaRunId          string `json:"gha-run-id,omitempty"`
	GhaRunAttempt     string `json:"gha-run-attempt,omitempty"`
	GhaRunNumber      string `json:"gha-run-number,omitempty"`
	CloudBeesApiUrl   string `json:"cloudbees-api-url,omitempty"`
	CloudBeesApiToken string `json:"cloudbees-api-token,omitempty"`
	GhaRepository     string `json:"gha-repository,omitempty"`
	GhaWorkflowRef    string `json:"gha-workflow-ref,omitempty"`
	GhaServerUrl      string `json:"gha-server-url,omitempty"`
	GhaJobName        string `json:"gha-job-name,omitempty"`
}
