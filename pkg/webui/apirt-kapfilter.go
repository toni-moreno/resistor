package webui

import (
	"github.com/go-macaron/binding"
	"github.com/influxdata/kapacitor/alert"
	//	"github.com/influxdata/kapacitor/models"
	//"github.com/toni-moreno/resistor/pkg/agent"
	//	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
	//	"time"
)

func NewAPIRtKapFilter(m *macaron.Macaron) error {

	bind := binding.Bind
	m.Group("/api/rt/kapfilter", func() {
		m.Post("/alert/", reqAlertSignedIn, bind(alert.Data{}), RTAlertHandler)
	})
	return nil
}

//RTAlertHandler xx
func RTAlertHandler(ctx *Context, al alert.Data) {
	rb := ctx.Req.Body()
	s, _ := rb.String()
	log.Debugf("REQ: %s", s)
	log.Debugf("ALERT: %#+v", al)
	log.Debugf("ALERT Data: %#+v", al.Data)
	log.Debugf("ALERT Series: %#+v", al.Data.Series)

	for _, v := range al.Data.Series {
		log.Debugf("ALERT ROW: %#+v", v)
	}
	ctx.JSON(200, "hola")
}
