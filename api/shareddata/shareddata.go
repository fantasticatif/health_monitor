package shareddata

import "os"

var Authenticated_api_route = "api/authenticated"

var Domain = os.Getenv("domain")

func ResetEnvVariables() {
	Domain = os.Getenv("domain")
}
