package sadmin

func (m *SAdminCreateUserRequest) Censor() {
	if m.Info != nil && m.Info.Password != "" {
		m.Info.Password = "..."
	}
}

func (m *SAdminResetPasswordRequest) Censor() {
	if m.Password != "" {
		m.Password = "..."
	}
	if m.Confirm != "" {
		m.Confirm = "..."
	}
}
