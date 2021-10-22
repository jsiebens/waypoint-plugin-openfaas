package platform

import "github.com/hashicorp/waypoint-plugin-sdk/component"

var _ component.Platform = (*Platform)(nil)
var _ component.Generation = (*Platform)(nil)

// DeployConfig is the configuration structure for the Platform.
type DeployConfig struct {
	Gateway     string `hcl:"gateway,optional"`
	TlsInsecure bool   `hcl:"tls_insecure,optional"`
	Username    string `hcl:"username,optional"`
	Password    string `hcl:"password,optional"`
	Token       string `hcl:"token,optional"`

	Namespace              string             `hcl:"namespace,optional"`
	FProcess               string             `hcl:"f_process,optional"`
	EnvVars                map[string]string  `hcl:"env_vars,optional"`
	Constraints            []string           `hcl:"constraints,optional"`
	Secrets                []string           `hcl:"secrets,optional"`
	Labels                 map[string]string  `hcl:"labels,optional"`
	Annotations            map[string]string  `hcl:"annotations,optional"`
	Limits                 *FunctionResources `hcl:"limits,block"`
	Requests               *FunctionResources `hcl:"requests,block"`
	ReadOnlyRootFilesystem bool               `hcl:"read_only_root_filesystem,optional"`
}

type FunctionResources struct {
	Memory string `hcl:"memory,optional"`
	CPU    string `hcl:"cpu,optional"`
}

type Platform struct {
	config DeployConfig
}

func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

func (p *Platform) ConfigSet(config interface{}) error {
	return nil
}
