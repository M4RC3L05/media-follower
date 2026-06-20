package views

import "github.com/m4rc3l05/media-follower/internal/common/middlewares"

type ContextViewKey string

var GlobalViewVarsContentViewKey ContextViewKey = "GLOBAL_VIEW_VARS_CONTENT_VIEW_KEY"

type GlobalViewVars struct {
	CSRFToken     *string
	FlashMessages *middlewares.FlashMessageData
}
