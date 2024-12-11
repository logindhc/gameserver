package controller

import (
	cherryGin "gameserver/cherry/components/gin"
	"gameserver/internal/code"
)

type MailController struct {
	cherryGin.BaseController
}

func (p *MailController) Init() {
	group := p.Group("/")
	group.GET("/mail/add", p.addMail)
}

// http://127.0.0.1/mail/add?sign=123&time=123
func (p *MailController) addMail(c *cherryGin.Context) {
	code.RenderResult(c, code.OK)
}
