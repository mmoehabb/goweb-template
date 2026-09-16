package pages

import (
	"github.com/a-h/templ"
	"goweb/pages/user"
)

var Registry = map[string]func() templ.Component{
	"/":             Index,
	"/user/profile": user.Profile,
}
