package conf

import (
	"fmt"
	"os"
)

var ENV string = os.Getenv("STAGE")

func IsRunningInDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	if os.Getenv("OS_ENV") == "docker" {
		return true
	}

	return false
}

// LocalParameterStore is required for setting environment variable locally to fetch ssm keys from dev
var LocalParameterStore string = os.Getenv("LOCAL_PS")

const (
	devSSM = "dev_ssm" // unexported because it is not required outside of this package scope
)

// Server ENV constants
const (
	ENV_PROD    = "prod"
	ENV_DEV     = "dev"
	ENV_DEV2    = "dev2"
	ENV_DEV3    = "dev3"
	ENV_DEV4    = "dev4"
	ENV_DEV5    = "dev5"
	ENV_DEV6    = "dev6"
	ENV_DEV7    = "dev7"
	ENV_DEV8    = "dev8"
	ENV_DEV9    = "dev9"
	ENV_LOCAL   = "local"
	ENV_UAT     = "uat"
	ENV_UAT2    = "uat2"
	ENV_UAT3    = "uat3"
	ENV_UAT4    = "uat4"
	ENV_UAT5    = "uat5"
	ENV_UAT6    = "uat6"
	ENV_UAT7    = "uat7"
	ENV_UAT8    = "uat8"
	ENV_UAT9    = "uat9"
	ENV_UAT10   = "uat10"
	ENV_UAT11   = "uat11"
	ENV_UAT12   = "uat12"
	ENV_DEV_TDL = "dev-fb-tdl"
	ENV_UAT_TDL = "uat-fb-tdl"
)

const REGION = "ap-south-1"

const (
	LendingServiceName = "lending-middleware-apis"
	LendingWorkers     = "lending-middleware-workers"
)

/*
Base URLs
*/
func getBaseURL() string {
	switch ENV {
	case ENV_PROD:
		return "wss://lendingapis.finbox.in"
	case ENV_UAT:
		return "wss://lendinguat.finbox.in"
	case ENV_UAT_TDL:
		return "wss://lendingtest.finbox.in"
	case ENV_DEV_TDL:
		return "wss://lendingtestdev.finbox.in"
	case ENV_LOCAL:
		return "ws://localhost:3332"
	case ENV_DEV + "2":
		return "wss://lendingtestdev2.finbox.in"
	}
	return fmt.Sprintf("wss://lending%s.finbox.in", ENV)
}

var BaseURL = getBaseURL()

/* Platform URLs
 */
func getLenderDashboardURL() string {
	switch ENV {
	case ENV_PROD:
		return "wss://lenders.finbox.in"
	case ENV_UAT10:
		return "ws://finbox-lender-dashboard-uat-10.s3-website.ap-south-1.amazonaws.com"
	default:
		return "wss://lendersuat.finbox.in"
	}
}

var LenderDashboardURL = getLenderDashboardURL()

func GetPlatformDashboardURL(organizationID string) string {
	switch ENV {
	case ENV_PROD:
		if organizationID == "eabdf414-d3f3-4f1a-ad71-f01d59d9c05b" {
			return "wss://adityabirla-platform.finbox.in"
		}
		return "wss://platform.finbox.in"
	case ENV_UAT2:
		return "ws://finbox-lending-platform-uat-2.s3-website.ap-south-1.amazonaws.com"
	case ENV_UAT3:
		return "ws://finbox-lending-platform-uat-3.s3-website.ap-south-1.amazonaws.com"
	case ENV_UAT4:
		return "ws://finbox-lending-platform-uat-4.s3-website.ap-south-1.amazonaws.com"
	case ENV_UAT5:
		return "ws://finbox-lending-platform-uat-5.s3-website.ap-south-1.amazonaws.com"
	case ENV_UAT6:
		return "ws://finbox-lending-platform-uat-6.s3-website.ap-south-1.amazonaws.com"
	case ENV_UAT7:
		return "ws://finbox-lending-platform-uat-7.s3-website.ap-south-1.amazonaws.com"
	case ENV_UAT8:
		return "ws://finbox-lending-platform-uat-8.s3-website.ap-south-1.amazonaws.com"
	case ENV_UAT9:
		return "ws://finbox-lending-platform-uat-9.s3-website.ap-south-1.amazonaws.com"
	case ENV_UAT10:
		return "ws://finbox-platform-dashboard-uat-10.s3-website.ap-south-1.amazonaws.com"
	case ENV_LOCAL:
		return "ws://localhost:3332"
	default:
		return "wss://platformuat.finbox.in"
	}
}
