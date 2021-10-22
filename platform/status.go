package platform

import (
	"context"
	"fmt"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"strings"
)

func (p *Platform) StatusFunc() interface{} {
	return p.status
}

func (p *Platform) status(
	ctx context.Context,
	ji *component.JobInfo,
	deploy *Deployment,
	ui terminal.UI,
) (*sdk.StatusReport, error) {
	st := ui.Status()
	defer st.Close()
	st.Update("Determining overall function health...")

	// Get our client
	client, err := p.newClient()
	if err != nil {
		return nil, err
	}

	info, err := client.GetFunctionInfo(ctx, deploy.GetFunctionName(), deploy.GetNamespace())
	if err != nil {
		return nil, err
	}

	report := &sdk.StatusReport{}

	if info.AvailableReplicas == info.Replicas {
		report.Health = sdk.StatusReport_READY
	} else {
		report.Health = sdk.StatusReport_PARTIAL
	}

	st.Update("Finished building report for OpenFaaS function")

	if report.Health == sdk.StatusReport_READY {
		st.Step(terminal.StatusOK, fmt.Sprintf("Function %q is reporting ready!", deploy.GetFunctionName()))
	} else {
		if report.Health == sdk.StatusReport_PARTIAL {
			st.Step(terminal.StatusWarn, fmt.Sprintf("Function %q is reporting partially available!", deploy.GetFunctionName()))
		} else {
			st.Step(terminal.StatusError, fmt.Sprintf("Function %q is reporting not ready!", deploy.GetFunctionName()))
		}

		// Extra advisory wording to let user know that the deployment could be still starting up
		// if the report was generated immediately after it was deployed or released.
		st.Step(terminal.StatusWarn, mixedHealthWarn)
	}

	return report, nil
}

var (
	mixedHealthWarn = strings.TrimSpace(`
Waypoint detected that the current function is not ready, however it
might be available or still starting up.
`)
)