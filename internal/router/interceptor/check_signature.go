package interceptor

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/pkg/env"
)

var whiteListPath = map[string]bool{}

func (i *interceptor) CheckSignature() core.HandlerFunc {
	return func(c core.Context) {
		if !env.Active().IsPro() {
			return
		}

		// 已移除 authorized 校验逻辑，生产环境直接放行
	}
}
