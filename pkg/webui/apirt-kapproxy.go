package webui

import (
	"io/ioutil"

	"net/http"

	//	"github.com/influxdata/kapacitor/models"
	//"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/agent"
	"gopkg.in/macaron.v1"
	//	"time"
)

// NewAPIRtKapProxy set the runtime Kapacitor Proxy
func NewAPIRtKapProxy(m *macaron.Macaron) error {

	//bind := binding.Bind
	m.Group("/api/rt/kapproxy", func() {
		m.Get("/:id/*.*", reqAlertSignedIn, RTKapProxyGETHandler)
		m.Post("/:id/*.*", reqAlertSignedIn, RTKapProxyPOSTHandler)
		m.Delete("/:id/*.*", reqAlertSignedIn, RTKapProxyOTHERHandler)
		m.Patch("/:id/*.*", reqAlertSignedIn, RTKapProxyOTHERHandler)
		m.Get("/test/:id/*.*", reqAlertSignedIn, RTKapProxyTest)
	})
	return nil
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// RTKapProxyTest GET handler
func RTKapProxyTest(ctx *Context) {
	id := ctx.Params(":id")
	//uri := ctx.Params(":uri")
	path := ctx.Params(":path")
	//log.Debugf("Get Proxy TEST request to: %s Kapacitor : URI : %s", id, uri)
	log.Debugf("Get Proxy TEST request to: %s Kapacitor : URI : %s ", id, path)
	ctx.JSON(200, &struct {
		ID   string
		Path string
	}{
		ID:   id,
		Path: path})

}

// RTKapProxyOTHERHandler for handing other than GET/POST methots ( rigth now DELETE/PATCH)
func RTKapProxyOTHERHandler(ctx *Context) {

	id := ctx.Params(":id")
	path := ctx.Params(":path")
	dev, err := agent.MainConfig.Database.GetKapacitorCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
		return
	}

	Method := ctx.Req.Request.Method

	log.Debugf("Get Proxy METHOD (%s) request to: %s Kapacitor : PATH : %s : KAP : %+v", Method, id, path, dev)

	r := ctx.Req.Request

	log.Debugf("REQUEST: %+v ", ctx.Req.Request)
	log.Debugf("RESPONSEWRITTER %+v", ctx.Resp)
	KURL := dev.URL + "/" + path
	req, err := http.NewRequest(Method, KURL, r.Body)
	if err != nil {
		log.Errorf("ERROR on Request handler %s URL %s", Method, KURL)
		ctx.JSON(404, err.Error())
		return
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("ERROR on Response ERR: %s ", err)
		ctx.JSON(404, err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	copyHeader(ctx.Resp.Header(), resp.Header)
	ctx.RawData(resp.StatusCode, body)

}

// RTKapProxyGETHandler GET handler
func RTKapProxyGETHandler(ctx *Context) {
	id := ctx.Params(":id")
	path := ctx.Params(":path")
	dev, err := agent.MainConfig.Database.GetKapacitorCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
		return
	}

	log.Debugf("Get Proxy GET request to: %s Kapacitor : PATH : %s : KAP : %+v", id, path, dev)

	log.Debugf("REQUEST: %+v ", ctx.Req.Request)
	log.Debugf("RESPONSEWRITTER %+v", ctx.Resp)
	KURL := dev.URL + "/" + path
	resp, err := http.Get(KURL)
	if err != nil {
		log.Errorf("ERROR on kapacitor hit %s", KURL)
		ctx.JSON(404, err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	copyHeader(ctx.Resp.Header(), resp.Header)
	ctx.RawData(resp.StatusCode, body)
}

// RTKapProxyPOSTHandler POST handler
func RTKapProxyPOSTHandler(ctx *Context) {
	id := ctx.Params(":id")
	path := ctx.Params(":path")
	dev, err := agent.MainConfig.Database.GetKapacitorCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
		return
	}

	log.Debugf("Get ProxyPOST request to: %s Kapacitor : PATH : %s : KAP : %+v", id, path, dev)

	r := ctx.Req.Request

	log.Debugf("REQUEST: %+v ", r)
	log.Debugf("RESPONSEWRITTER %+v", ctx.Resp)
	log.Debugf("BODY: %+v", r.Body)
	KURL := dev.URL + "/" + path
	contentType := r.Header.Get("Content-Type")
	resp, err := http.Post(KURL, contentType, r.Body)
	if err != nil {
		log.Errorf("ERROR on kapacitor hit %s", KURL)
		ctx.JSON(404, err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	copyHeader(ctx.Resp.Header(), resp.Header)
	ctx.RawData(resp.StatusCode, body)

}
