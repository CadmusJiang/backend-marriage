package authz

import "unicode/utf8"

// ShouldMask 非公司/组管理员需要对敏感字段做脱敏
func ShouldMask(scope AccessScope) bool {
	return !(scope.RoleType == "company_manager" || scope.RoleType == "group_manager")
}

// MaskPhone 行业惯例：保留前三后四
func MaskPhone(phone string) string {
	if phone == "" {
		return phone
	}
	// 仅对长度>=7 的数字串做常规掩码
	if utf8.RuneCountInString(phone) < 7 {
		return "***" // 过短的不暴露
	}
	// 简单按字节切分，适用于常见数字手机号
	b := []rune(phone)
	prefix := string(b[:3])
	suffixLen := 4
	if len(b) < 7+(0) { // already guarded above
		suffixLen = 2
	}
	if len(b) < 3+suffixLen {
		return "***"
	}
	suffix := string(b[len(b)-suffixLen:])
	return prefix + "****" + suffix
}

// MaskWechat 行业惯例：保留前2后2
func MaskWechat(wechat string) string {
	if wechat == "" {
		return wechat
	}
	r := []rune(wechat)
	if len(r) <= 4 {
		return string(r[:1]) + "***"
	}
	prefix := string(r[:2])
	suffix := string(r[len(r)-2:])
	return prefix + "****" + suffix
}
