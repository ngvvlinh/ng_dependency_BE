package convert

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"o.o/api/main/authorization"
	"o.o/api/main/invitation"
	"o.o/backend/com/main/authorization/convert"
	"o.o/backend/com/main/invitation/model"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

// +gen:convert: o.o/backend/com/main/invitation/model  -> o.o/api/main/invitation
// +gen:convert: o.o/api/main/invitation

const ExpiresIn = 24 * time.Hour

func createInvitation(args *invitation.CreateInvitationArgs, out *invitation.Invitation) {
	apply_invitation_CreateInvitationArgs_invitation_Invitation(args, out)
	out.ID = cm.NewID()
}

func ConvertStringsToRoles(args []string) (roles []authorization.Role) {
	for _, arg := range args {
		roles = append(roles, authorization.Role(arg))
	}
	return
}

func ConvertRolesToStrings(roles []authorization.Role) (outs []string) {
	for _, role := range roles {
		outs = append(outs, role.String())
	}
	return
}

func GenerateToken(secretKey, email string, accountID int64, roles []authorization.Role, expiresAt int64) (string, error) {
	if expiresAt == 0 {
		expiresAt = time.Now().Add(ExpiresIn).Unix()
	}
	claims := &invitation.Claims{
		Email:     email,
		AccountID: dot.ID(accountID),
		Roles:     roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetExpiresAt(secretKey string, token string) (time.Time, error) {
	claims := &invitation.Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return time.Now(), err
	}

	return time.Unix(claims.StandardClaims.ExpiresAt, 0), nil
}

func InvitationDB(args *invitation.Invitation, out *model.Invitation) {
	convert_invitation_Invitation_invitationmodel_Invitation(args, out)
	out.Roles = convert.ConvertRolesToStrings(args.Roles)
}

func Invitation(args *model.Invitation, out *invitation.Invitation) {
	convert_invitationmodel_Invitation_invitation_Invitation(args, out)
	out.Roles = convert.ConvertStringsToRoles(args.Roles)
}
