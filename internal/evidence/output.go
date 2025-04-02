package evidence

type EvidenceInfo struct {
	Content string `json:"content,omitempty"`
	Format  string `json:"format,omitempty"`
}

type ProviderInfo struct {
	RunId      string `json:"run_id,omitempty"`
	RunAttempt string `json:"run_attempt,omitempty"`
	RunNumber  string `json:"run_number,omitempty"`
	JobName    string `json:"job_name,omitempty"`
	Provider   string `json:"provider,omitempty"`
}

type Output struct {
	EvidenceInfo EvidenceInfo `json:"evidence_info,omitempty"`
	ProviderInfo ProviderInfo `json:"provider_info,omitempty"`
}
