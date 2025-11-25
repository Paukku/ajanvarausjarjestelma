package auth

import pbcommon "github.com/Paukku/ajanvarausjarjestelma/backend/pb/common"

type RoleString string

var RoleStringToEnum = map[RoleString]pbcommon.UserRole{
	"Admin":    pbcommon.UserRole_ADMIN,
	"Owner":    pbcommon.UserRole_OWNER,
	"Employee": pbcommon.UserRole_EMPLOYEE,
}
