package etop

func (x AccountType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}

func (m *CreateUserRequest) Censor() {
	if m.Password != "" {
		m.Password = "..."
	}
	if m.RegisterToken != "" {
		m.RegisterToken = "..."
	}
}

func (m *LoginRequest) Censor() {
	if m.Password != "" {
		m.Password = "..."
	}
}

func (m *ChangePasswordRequest) Censor() {
	if m.CurrentPassword != "" {
		m.CurrentPassword = "..."
	}
	if m.NewPassword != "" {
		m.NewPassword = "..."
	}
	if m.ConfirmPassword != "" {
		m.ConfirmPassword = "..."
	}
}

func (m *ChangePasswordUsingTokenRequest) Censor() {
	if m.ResetPasswordToken != "" {
		m.ResetPasswordToken = "..."
	}
	if m.NewPassword != "" {
		m.NewPassword = "..."
	}
	if m.ConfirmPassword != "" {
		m.ConfirmPassword = "..."
	}
}
