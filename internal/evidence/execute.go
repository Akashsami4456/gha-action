package evidence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details []any  `json:"details"`
}

func (config *Config) Run(_ context.Context) (err error) {
	validationError := setEnvVars(config)
	if validationError != nil {
		return validationError
	}
	cloudEventData := prepareCloudEventData(config)

	cloudEvent, err := prepareCloudEvent(config, cloudEventData)
	if err != nil {
		return err
	}

	err = sendCloudEvent(cloudEvent, config)
	if err != nil {
		return err
	}
	return nil
}

func setEnvVars(cfg *Config) error {
	ghaRunId := os.Getenv(GithubRunId)
	if ghaRunId == "" {
		return fmt.Errorf(GithubRunId + " is not set in the environment")
	}
	cfg.GhaRunId = ghaRunId

	ghaRunAttempt := os.Getenv(GithubRunAttempt)
	if ghaRunAttempt == "" {
		return fmt.Errorf(GithubRunAttempt + " is not set in the environment")
	}
	cfg.GhaRunAttempt = ghaRunAttempt

	cloudBeesApiUrl := os.Getenv(CloudbeesApiUrl)
	if cloudBeesApiUrl == "" {
		return fmt.Errorf(CloudbeesApiUrl + " is not set in the environment")
	}
	cfg.CloudBeesApiUrl = cloudBeesApiUrl

	cloudBeesApiToken := os.Getenv(CloudbeesApiToken)
	if cloudBeesApiToken == "" {
		return fmt.Errorf(CloudbeesApiToken + " is not set in the environment")
	}
	cfg.CloudBeesApiToken = cloudBeesApiToken

	ghaRunNumber := os.Getenv(GithubRunNumber)
	if ghaRunNumber == "" {
		return fmt.Errorf(GithubRunNumber + " is not set in the environment")
	}

	cfg.GhaRunNumber = ghaRunNumber

	ghaRepository := os.Getenv(GithubRepository)
	if ghaRepository == "" {
		return fmt.Errorf(GithubRepository + " is not set in the environment")
	}

	cfg.GhaRepository = ghaRepository

	ghaWorkflowRef := os.Getenv(GithubWorkflowRef)
	if ghaWorkflowRef == "" {
		return fmt.Errorf(GithubWorkflowRef + " is not set in the environment")
	}

	cfg.GhaWorkflowRef = ghaWorkflowRef

	ghaJobName := os.Getenv(GithubJobName)
	if ghaJobName == "" {
		return fmt.Errorf(GithubJobName + " is not set in the environment")
	}

	cfg.GhaJobName = ghaJobName

	cfg.GhaServerUrl = os.Getenv(GithubServerUrl)

	content := os.Getenv(Content)
	if content == "" {
		return fmt.Errorf(Content + " is not provided in the environment")
	}
	cfg.Content = content
	fmt.Println("Content: ", cfg.Content)
	return nil
}

func prepareCloudEventData(cfg *Config) Output {

	evidenceInfo := &PublishEvidence{
		Content: cfg.Content,
		Format:  cfg.Format,
	}
	providerInfo := &ProviderInfo{
		RunId:      cfg.GhaRunId,
		RunAttempt: cfg.GhaRunAttempt,
		RunNumber:  cfg.GhaRunNumber,
		JobName:    cfg.GhaJobName,
		Provider:   GithubProvider,
	}
	output := Output{
		ProviderInfo:    *providerInfo,
		PublishEvidence: *evidenceInfo,
	}
	fmt.Println("Output set data")
	fmt.Println(PrettyPrint(output))
	return output
}

func prepareCloudEvent(config *Config, output Output) (cloudevents.Event, error) {
	cloudEvent := cloudevents.NewEvent()
	cloudEvent.SetID(uuid.NewString())
	cloudEvent.SetSubject(getSubject(config))
	cloudEvent.SetType(PublishEvidenceType)
	cloudEvent.SetSource(getSource(config))
	cloudEvent.SetSpecVersion(SpecVersion)
	cloudEvent.SetTime(time.Now())
	err := cloudEvent.SetData(ContentTypeJson, output)
	fmt.Println("CloudEvent set data")
	fmt.Println(PrettyPrint(cloudEvent))
	if err != nil {
		return cloudevents.Event{}, fmt.Errorf("failed to set data: %v", err)
	}
	return cloudEvent, nil
}

func getSubject(config *Config) string {
	return config.GhaWorkflowRef + "|" + config.GhaRunId + "|" + config.GhaRunAttempt + "|" + config.GhaRunNumber
}

func getSource(config *Config) string {
	sourcePrefix := GithubProvider
	if config.GhaServerUrl != "" {
		sourcePrefix = config.GhaServerUrl + "/"
	}
	return sourcePrefix + config.GhaRepository
}

func sendCloudEvent(cloudEvent cloudevents.Event, config *Config) error {
	eventJSON, err := json.Marshal(cloudEvent)
	if err != nil {
		return fmt.Errorf("error encoding CloudEvent JSON %s", err)
	}
	req, _ := http.NewRequest(PostMethod, getCloudbeesFullUrl(config), bytes.NewBuffer(eventJSON))
	fmt.Println(PrettyPrint(cloudEvent))
	// For Local Testing
	//req, _ := http.NewRequest(PostMethod, "http://localhost:8080/events", bytes.NewBuffer(eventJSON))

	req.Header.Set(ContentTypeHeaderKey, ContentTypeCloudEventsJson)
	req.Header.Set(AuthorizationHeaderKey, Bearer+config.CloudBeesApiToken)
	client := &http.Client{}
	resp, err := client.Do(req) // Fire and forget

	if err != nil {
		return fmt.Errorf("error sending CloudEvent to platform %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
		}
		var errorResponse ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			fmt.Println("Error unmarshaling response body:", err)
		}
		return fmt.Errorf("error sending CloudEvent to platform: %s", errorResponse.Message)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	fmt.Println("CloudEvent sent successfully!")
	return nil
}

// PrettyPrint converts the input to JSON string
func PrettyPrint(in interface{}) string {
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		fmt.Println("error marshalling response", err)
	}
	return string(data)
}

func getCloudbeesFullUrl(config *Config) string {
	if !strings.HasSuffix(config.CloudBeesApiUrl, "/") {
		config.CloudBeesApiUrl += "/"
	}
	return config.CloudBeesApiUrl + "v3/external-events"
}
