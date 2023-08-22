package config

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
)

type Config struct {
	PlantonCloudKubeAgentMachineAccountEmail string
	PlantonCloudKubeAgentClientSecret        string
	PlantonCloudKubeAgentKubernetesClusterId string
	PlantonCloudServiceApiEndpoint           string
	OpenCostApiEndpoint                      string
	OpenCostPollingIntervalSeconds           int
	TokenExpirationBufferMinutes             int
	TokenExpirationCheckIntervalSeconds      int
}

const (
	EnvVarPlantonCloudKubeAgentMachineAccountEmail = "PLANTON_CLOUD_KUBE_AGENT_MACHINE_ACCOUNT_EMAIL"
	EnvVarPlantonCloudKubeAgentClientSecret        = "PLANTON_CLOUD_KUBE_AGENT_CLIENT_SECRET"
	EnvVarPlantonCloudKubeAgentKubernetesClusterId = "PLANTON_CLOUD_KUBE_AGENT_HOSTING_CLUSTER_ID"
	EnvVarPlantonCloudServiceApiEndpoint           = "PLANTON_CLOUD_SERVICE_API_ENDPOINT"
	EnvVarOpenCostApiEndpoint                      = "OPEN_COST_API_ENDPOINT"
	EnvVarOpenCostPollingIntervalSeconds           = "OPEN_COST_POLLING_INTERVAL_SECONDS"
	EnvVarTokenExpirationBufferMinutes             = "TOKEN_EXPIRATION_BUFFER_MINUTES"
	EnvVarTokenExpirationCheckIntervalSeconds      = "TOKEN_EXPIRATION_CHECK_INTERVAL_SECONDS"
)

// Load config from environment variables
func Load() (*Config, error) {
	clientId, ok := os.LookupEnv(EnvVarPlantonCloudKubeAgentMachineAccountEmail)
	if !ok {
		return nil, errors.Errorf("%s environment variable is not set", EnvVarPlantonCloudKubeAgentMachineAccountEmail)
	}
	clientSecret, ok := os.LookupEnv(EnvVarPlantonCloudKubeAgentClientSecret)
	if !ok {
		return nil, errors.Errorf("%s environment variable is not set", EnvVarPlantonCloudKubeAgentClientSecret)
	}
	kubernetesClusterId, ok := os.LookupEnv(EnvVarPlantonCloudKubeAgentKubernetesClusterId)
	if !ok {
		return nil, errors.Errorf("%s environment variable is not set", EnvVarPlantonCloudKubeAgentKubernetesClusterId)
	}
	plantonCloudServiceApiEndpoint, ok := os.LookupEnv(EnvVarPlantonCloudServiceApiEndpoint)
	if !ok {
		return nil, errors.Errorf("%s environment variable is not set", EnvVarPlantonCloudServiceApiEndpoint)
	}
	openCostApiEndpoint, ok := os.LookupEnv(EnvVarOpenCostApiEndpoint)
	if !ok {
		return nil, errors.Errorf("%s environment variable is not set", EnvVarOpenCostApiEndpoint)
	}
	openCostPollingIntervalSecondsStr, ok := os.LookupEnv(EnvVarOpenCostPollingIntervalSeconds)
	if !ok {
		return nil, errors.Errorf("%s environment variable is not set", EnvVarOpenCostPollingIntervalSeconds)
	}
	openCostPollingIntervalSeconds, err := strconv.Atoi(openCostPollingIntervalSecondsStr)
	if err != nil {
		return nil, errors.Errorf("%s environment variable should be set to integer value", EnvVarOpenCostPollingIntervalSeconds)
	}
	tokenExpirationBufferMinutesStr, ok := os.LookupEnv(EnvVarTokenExpirationBufferMinutes)
	if !ok {
		return nil, errors.Errorf("%s environment variable is not set", EnvVarTokenExpirationBufferMinutes)
	}
	tokenExpirationBufferMinutes, err := strconv.Atoi(tokenExpirationBufferMinutesStr)
	if err != nil {
		return nil, errors.Errorf("%s environment variable should be set to integer value", EnvVarTokenExpirationBufferMinutes)
	}
	tokenExpirationCheckIntervalSecondsStr, ok := os.LookupEnv(EnvVarTokenExpirationCheckIntervalSeconds)
	if !ok {
		return nil, errors.Errorf("%s environment variable is not set", EnvVarTokenExpirationCheckIntervalSeconds)
	}
	tokenExpirationCheckIntervalSeconds, err := strconv.Atoi(tokenExpirationCheckIntervalSecondsStr)
	if err != nil {
		return nil, errors.Errorf("%s environment variable should be set to integer value", EnvVarTokenExpirationCheckIntervalSeconds)
	}
	return &Config{
		PlantonCloudKubeAgentMachineAccountEmail: clientId,
		PlantonCloudKubeAgentClientSecret:        clientSecret,
		PlantonCloudKubeAgentKubernetesClusterId: kubernetesClusterId,
		PlantonCloudServiceApiEndpoint:           plantonCloudServiceApiEndpoint,
		OpenCostApiEndpoint:                      openCostApiEndpoint,
		OpenCostPollingIntervalSeconds:           openCostPollingIntervalSeconds,
		TokenExpirationBufferMinutes:             tokenExpirationBufferMinutes,
		TokenExpirationCheckIntervalSeconds:      tokenExpirationCheckIntervalSeconds,
	}, nil
}
