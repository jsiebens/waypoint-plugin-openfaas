package main

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
	"github.com/jsiebens/waypoint-plugin-openfaas/platform"
)

func main() {
	sdk.Main(sdk.WithComponents(
		&platform.Platform{},
	))
}
