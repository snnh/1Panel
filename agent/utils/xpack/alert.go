package xpack

import (
	"net/http"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/app/model"
	alertUtil "github.com/1Panel-dev/1Panel/agent/utils/alert"
)

// DeliveryResult reports whether a delivery was queued for asynchronous processing.
type DeliveryResult struct {
	Queued bool
	LogID  uint
}

// GetNodeErrorAlert returns the number of nodes reporting errors. Community
// edition manages a single local agent, so it is always zero.
func GetNodeErrorAlert() (uint, error) {
	return 0, nil
}

// DeliverCustomWebhookAlertLog sends a custom webhook alert synchronously.
func DeliverCustomWebhookAlertLog(
	alertType string,
	info dto.AlertDTO,
	create dto.AlertLogCreate,
	project string,
	params []dto.Param,
	config model.AlertConfig,
	transport *http.Transport,
	agentInfo *dto.AgentInfo,
	_ dto.AlertTaskMetadata,
) (DeliveryResult, error) {
	err := alertUtil.CreateCustomWebhookAlertLog(alertType, info, create, project, params, config, transport, agentInfo)
	return DeliveryResult{}, err
}

// DeliverTaskScanCustomWebhookAlertLog sends a task scan custom webhook alert synchronously.
func DeliverTaskScanCustomWebhookAlertLog(
	alert dto.AlertDTO,
	alertType string,
	create dto.AlertLogCreate,
	pushAlert dto.PushAlert,
	config model.AlertConfig,
	transport *http.Transport,
	agentInfo *dto.AgentInfo,
	_ dto.AlertTaskMetadata,
) (DeliveryResult, error) {
	err := alertUtil.CreateTaskScanCustomWebhookAlertLog(alert, alertType, create, pushAlert, config, transport, agentInfo)
	return DeliveryResult{}, err
}

// TestCustomWebhook sends a test request to the resolved custom webhook config.
func TestCustomWebhook(config dto.AlertCustomWebhookResolvedConfig) (dto.AlertConfigTestResult, error) {
	agentInfo, _ := GetAgentInfo()
	return alertUtil.TestCustomWebhook(config, LoadRequestTransport(), agentInfo)
}
