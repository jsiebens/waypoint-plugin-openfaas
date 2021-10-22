package platform

import (
	"context"
	"fmt"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/openfaas/faas-cli/proxy"
	"github.com/openfaas/faas-cli/stack"
)

func (p *Platform) DeployFunc() interface{} {
	return p.deploy
}

func (p *Platform) deploy(ctx context.Context,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI) (*Deployment, error) {

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	// Get our client
	client, err := p.newClient()
	if err != nil {
		return nil, err
	}

	// Create our deployment and set an initial ID
	var result Deployment
	result.FunctionName = src.App
	result.Namespace = p.config.Namespace

	// Build our env vars
	env := map[string]string{
		"PORT": "8080",
	}

	for k, v := range p.config.EnvVars {
		env[k] = v
	}

	for k, v := range deployConfig.Env() {
		env[k] = v
	}

	resourceRequest := proxy.FunctionResourceRequest{}
	if p.config.Requests != nil {
		resourceRequest.Requests = &stack.FunctionResources{
			Memory: p.config.Requests.Memory,
			CPU:    p.config.Requests.CPU,
		}
	}
	if p.config.Limits != nil {
		resourceRequest.Limits = &stack.FunctionResources{
			Memory: p.config.Limits.Memory,
			CPU:    p.config.Limits.CPU,
		}
	}

	req := &proxy.DeployFunctionSpec{
		Update:                  true,
		FunctionName:            result.FunctionName,
		Image:                   img.Name(),
		Namespace:               p.config.Namespace,
		FProcess:                p.config.FProcess,
		EnvVars:                 env,
		Constraints:             p.config.Constraints,
		Secrets:                 p.config.Secrets,
		Labels:                  p.config.Labels,
		Annotations:             p.config.Annotations,
		FunctionResourceRequest: resourceRequest,
		ReadOnlyRootFilesystem:  p.config.ReadOnlyRootFilesystem,
	}

	st.Update("Deploying function...")

	statusCode := client.DeployFunction(ctx, req)
	if statusCode >= 300 {
		return nil, fmt.Errorf("error deploying function: status code %d", statusCode)
	}

	st.Step(terminal.StatusOK, "Function deployment successful")

	return &result, nil
}
