package service

import (
	"context"
	"os/exec"
	"strings"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/docker/docker/api/types/system"
)

const gpuEnvironmentCacheKey = "gpu:environment"

func findGPUContainerCommand(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return path
}

func findNVIDIAContainerCommand(name string, _ map[string]interface{}) string {
	return findGPUContainerCommand(name)
}

func checkNVIDIAToolkit(ctx context.Context, findCommand func(string) string, runVersion func(context.Context, string) (string, error)) dto.GPUEnvironmentCheck {
	check := dto.GPUEnvironmentCheck{Vendor: "nvidia", Item: "container-toolkit"}
	executable := findCommand("nvidia-container-toolkit")
	if executable == "" {
		executable = findCommand("nvidia-container-runtime")
	}
	if executable == "" {
		check.Status = "failed"
		check.Reason = "NVIDIA Container Toolkit is not installed"
		return check
	}
	version, err := runVersion(ctx, executable)
	if err != nil {
		check.Status = "failed"
		check.Reason = "NVIDIA Container Toolkit is not executable"
		return check
	}
	check.Status = "passed"
	check.Version = strings.TrimSpace(version)
	return check
}

func checkNVIDIADockerConfig(info system.Info, _ func() (map[string]interface{}, error), findCommand func(string) string) dto.GPUEnvironmentCheck {
	check := dto.GPUEnvironmentCheck{Vendor: "nvidia", Item: "docker-runtime"}
	if _, ok := info.Runtimes["nvidia"]; !ok {
		check.Status = "failed"
		check.Reason = "Docker does not register the NVIDIA runtime"
		return check
	}
	if findCommand("nvidia-container-runtime") == "" {
		check.Status = "failed"
		check.Reason = "nvidia-container-runtime is not executable"
		return check
	}
	check.Status = "passed"
	return check
}
