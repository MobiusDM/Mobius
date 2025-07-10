package webhooks

import (
	"time"

	"github.com/notawar/mobius/backend/server/mobius"
)

type VulnArgs struct {
	Vulnerablities []mobius.SoftwareVulnerability
	Meta           map[string]mobius.CVEMeta
	AppConfig      *mobius.AppConfig
	Time           time.Time
}
