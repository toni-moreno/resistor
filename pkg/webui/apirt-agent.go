package webui

import (
	//	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	//	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
	//	"time"
)

func NewApiRtAgent(m *macaron.Macaron) error {

	//	bind := binding.Bind

	m.Group("/api/rt/agent", func() {
		m.Get("/info/version/", reqSignedIn, RTGetVersion)
	})

	return nil
}

//RTGetVersion xx
func RTGetVersion(ctx *Context) {
	info := agent.GetRInfo()
	ctx.JSON(200, &info)
}
