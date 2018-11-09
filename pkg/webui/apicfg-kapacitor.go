package webui

import (
	"time"

	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"github.com/toni-moreno/resistor/pkg/kapa"
	"gopkg.in/macaron.v1"
)

// NewAPICfgKapacitor Kapacitor ouput
func NewAPICfgKapacitor(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/kapacitor", func() {
		m.Get("/", reqSignedIn, GetKapacitor)
		m.Post("/", reqSignedIn, bind(config.KapacitorCfg{}), AddKapacitor)
		m.Put("/:id", reqSignedIn, bind(config.KapacitorCfg{}), UpdateKapacitor)
		m.Delete("/:id", reqSignedIn, DeleteKapacitor)
		m.Get("/:id", reqSignedIn, GetKapacitorCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetKapacitorAffectOnDel)
		m.Post("/ping", reqSignedIn, bind(config.KapacitorCfg{}), PingKapacitor)
	})

	return nil
}

// PingKapacitor Pings kapacitor server and returns time elapsed and kapacitor server version
func PingKapacitor(ctx *Context, dev config.KapacitorCfg) {
	_, elapsed, version, err := kapa.GetKapaClient(dev)
	if err != nil {
		log.Warningf("Error pinging Kapacitor Server %s: Err: %s", dev.ID, err)
		ctx.JSON(404, err.Error())
		return
	}
	ctx.JSON(200, &struct {
		Message string
		Elapsed time.Duration
	}{
		Message: version,
		Elapsed: elapsed,
	})
}

// GetKapacitor Return kapacitor servers list to frontend
func GetKapacitor(ctx *Context) {
	devcfgarray, err := kapa.GetKapaServers("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Kapacitor :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting Kapacitor %+v", &devcfgarray)
}

// AddKapacitor Insert new Kapacitor to the internal BBDD --pending--
func AddKapacitor(ctx *Context, dev config.KapacitorCfg) {
	log.Printf("ADDING Kapacitor %+v", dev)
	affected, err := agent.MainConfig.Database.AddKapacitorCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for Kapacitor %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		// Deploy the templates from resistor database into this kapacitor server
		tplcfgarray, err := GetTemplates("")
		if err == nil {
			kapasrvsarray := []*config.KapacitorCfg{&dev}
			for _, tplcfg := range tplcfgarray {
				_, _, _, _ = kapa.SetKapaTemplate(*tplcfg, kapasrvsarray)
			}
		}
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateKapacitor --pending--
func UpdateKapacitor(ctx *Context, dev config.KapacitorCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateKapacitorCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for Kapacitor %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteKapacitor delete a backend config
func DeleteKapacitor(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelKapacitorCfg(id)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetKapacitorCfgByID --pending--
func GetKapacitorCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetKapacitorCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetKapacitorAffectOnDel --pending--
func GetKapacitorAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetKapacitorCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for Kapacitors %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
