package platform

import (
	"context"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

func (p *Platform) DestroyFunc() interface{} {
	return p.destroy
}

func (p *Platform) destroy(ctx context.Context, ui terminal.UI, deployment *Deployment) error {
	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	// Get our client
	client, err := p.newClient()
	if err != nil {
		return err
	}

	err = client.DeleteFunction(ctx, deployment.GetFunctionName(), deployment.GetNamespace())

	if err == nil || isFunctionNotFound(err) {
		return nil
	}

	return err
}
