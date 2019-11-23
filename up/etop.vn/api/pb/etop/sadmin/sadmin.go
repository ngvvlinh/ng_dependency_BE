package sadmin

func (m *SAdminCreateUserRequest) Censor() {
	if m.GetInfo().Password != "" {
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
