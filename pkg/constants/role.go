package constants

import "github.com/google/uuid"

var (
	RoleIDAdmin   = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	RoleNameAdmin = "admin"
	RoleIDUser    = uuid.MustParse("00000000-0000-0000-0000-000000000002")
	RoleNameUser  = "user"
)
