package permissions

type DeclaredPermission struct {
	PermissionCode string
	PermissionName string
	PermissionType string
	ResourceCode   string
	GroupCode      string
	Description    string
}

func DeclaredPermissions() []DeclaredPermission {
	return []DeclaredPermission{
		{
			PermissionCode: "captcha:read",
			PermissionName: "查看验证码日志",
			PermissionType: "READ",
			ResourceCode:   "captcha",
			GroupCode:      "SYSTEM",
			Description:    "验证码验证日志与统计查看",
		},
		{
			PermissionCode: "captcha:write",
			PermissionName: "管理验证码策略",
			PermissionType: "WRITE",
			ResourceCode:   "captcha",
			GroupCode:      "SYSTEM",
			Description:    "验证码配置与策略管理",
		},
	}
}
