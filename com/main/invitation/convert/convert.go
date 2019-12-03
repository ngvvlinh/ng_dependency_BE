package convert

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"etop.vn/api/main/invitation"
	"etop.vn/backend/com/main/invitation/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/capi/dot"
)

// +gen:convert: etop.vn/backend/com/main/invitation/model -> etop.vn/api/main/invitation
// +gen:convert: etop.vn/api/main/invitation

const ExpiresIn = 24 * time.Hour

func createInvitation(args *invitation.CreateInvitationArgs, out *invitation.Invitation) {
	apply_invitation_CreateInvitationArgs_invitation_Invitation(args, out)
	out.ID = cm.NewID()
}

func ConvertStringsToRoles(args []string) (roles []invitation.Role) {
	for _, arg := range args {
		roles = append(roles, invitation.Role(arg))
	}
	return
}

func ConvertRolesToStrings(roles []invitation.Role) (outs []string) {
	for _, role := range roles {
		outs = append(outs, string(role))
	}
	return
}

func GenerateToken(secretKey, email string, accountID dot.ID, roles []invitation.Role, expiresAt int64) (string, error) {
	if expiresAt == 0 {
		expiresAt = time.Now().Add(ExpiresIn).Unix()
	}
	claims := &invitation.Claims{
		Email:     email,
		AccountID: accountID,
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
	out.Roles = ConvertRolesToStrings(args.Roles)
}

func Invitation(args *model.Invitation, out *invitation.Invitation) {
	convert_invitationmodel_Invitation_invitation_Invitation(args, out)
	out.Roles = ConvertStringsToRoles(args.Roles)
}
