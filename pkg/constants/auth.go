package constants

import (
	"time"

	"github.com/google/uuid"
)

const AuthSession = "AUTH_SESSION"

const AuthAccessTokenExpiry = time.Minute * 15 // 15 minute

const AuthRefreshTokenExpiry = time.Hour * 24 * 7

const XAccessToken = "X-Access-Token"
const XRefreshToken = "X-Refresh-Token"
const XLocale = "X-Locale"
const XDevice = "X-Device"
const XLocate = "X-Locate"
const XPlatform = "X-Platform"
const XLongitude = "X-Longitude"

var (
	AuthAdminID        = uuid.MustParse("00000000-0000-7000-8000-000000000001")
	AuthAdminFirstName = "Admin"
	AuthAdminLastName  = "Last"
	AuthAdminEmail     = "admin@domain.com"
	AuthAdminPassword  = "!Password123"
)
