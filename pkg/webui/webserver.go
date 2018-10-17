package webui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"github.com/go-macaron/toolbox"

	"crypto/md5"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"github.com/toni-moreno/resistor/pkg/kapa"
	"gopkg.in/macaron.v1"
)

var (
	logDir     string
	confDir    string
	log        *logrus.Logger
	confHTTP   *config.HTTPConfig
	instanceID string
)

// SetLogDir et dir for logs
func SetLogDir(dir string) {
	logDir = dir
}

// SetConfDir et dir for logs
func SetConfDir(dir string) {
	confDir = dir
}

// SetLogger set output log
func SetLogger(l *logrus.Logger) {
	log = l
}

//UserLogin for login purposes
type UserLogin struct {
	UserName string `form:"username" binding:"Required"`
	Password string `form:"password" binding:"Required"`
}

var cookie string

// WebServer the main process
func WebServer(publicPath string, httpPort int, cfg *config.HTTPConfig, id string) {
	confHTTP = cfg
	instanceID = id
	var port int
	if cfg.Port > 0 {
		port = cfg.Port
	} else {
		port = httpPort
	}

	bind := binding.Bind

	/*	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})*/

	f, _ := os.OpenFile(logDir+"/http_access.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	m := macaron.NewWithLogger(f)
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(toolbox.Toolboxer(m))
	// register middleware
	m.Use(GetContextHandler())
	//	m.Use(gzip.Gziper())
	log.Infof("setting HTML Static Path to %s", publicPath)
	m.Use(macaron.Static(publicPath,
		macaron.StaticOptions{
			// Prefix is the optional prefix used to serve the static directory content. Default is empty string.
			Prefix: "",
			// SkipLogging will disable [Static] log messages when a static file is served. Default is false.
			SkipLogging: false,
			// IndexFile defines which file to serve as index if it exists. Default is "index.html".
			IndexFile: "index.html",
			// Expires defines which user-defined function to use for producing a HTTP Expires Header. Default is nil.
			// https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
			Expires: func() string { return "max-age=0" },
		}))

	//Cookie should be unique for each resistor instance ,
	//if cockie_id is not set it takes the instanceID value to generate a unique array with as a md5sum

	cookie = confHTTP.CookieID

	if len(confHTTP.CookieID) == 0 {
		currentsum := md5.Sum([]byte(instanceID))
		cookie = fmt.Sprintf("%x", currentsum)
	}

	m.Use(Sessioner(session.Options{
		// Name of provider. Default is "memory".
		Provider: "memory",
		// Provider configuration, it's corresponding to provider.
		ProviderConfig: "",
		// Cookie name to save session ID. Default is "MacaronSession".
		CookieName: "resistor-sess-" + cookie,
		// Cookie path to store. Default is "/".
		CookiePath: "/",
		// GC interval time in seconds. Default is 3600.
		Gclifetime: 3600,
		// Max life time in seconds. Default is whatever GC interval time is.
		Maxlifetime: 3600,
		// Use HTTPS only. Default is false.
		Secure: false,
		// Cookie life time. Default is 0.
		CookieLifeTime: 0,
		// Cookie domain name. Default is empty.
		Domain: "",
		// Session ID length. Default is 16.
		IDLength: 16,
		// Configuration section name. Default is "session".
		Section: "session",
	}))

	m.Use(macaron.Renderer(macaron.RenderOptions{
		// Directory to load templates. Default is "templates".
		Directory: "pkg/templates",
		// Extensions to parse template files from. Defaults are [".tmpl", ".html"].
		Extensions: []string{".tmpl", ".html"},
		// Funcs is a slice of FuncMaps to apply to the template upon compilation. Default is [].
		/*Funcs: []template.FuncMap{map[string]interface{}{
			"AppName": func() string {
				return "resistor"
			},
			"AppVer": func() string {
				return "0.5.1"
			},
		}},*/
		// Delims sets the action delimiters to the specified strings. Defaults are ["{{", "}}"].
		Delims: macaron.Delims{"{{", "}}"},
		// Appends the given charset to the Content-Type header. Default is "UTF-8".
		Charset: "UTF-8",
		// Outputs human readable JSON. Default is false.
		IndentJSON: true,
		// Outputs human readable XML. Default is false.
		IndentXML: true,
		// Prefixes the JSON output with the given bytes. Default is no prefix.
		// PrefixJSON: []byte("macaron"),
		// Prefixes the XML output with the given bytes. Default is no prefix.
		// PrefixXML: []byte("macaron"),
		// Allows changing of output to XHTML instead of HTML. Default is "text/html".
		HTMLContentType: "text/html",
	}))
	/*
		m.Use(cache.Cacher(cache.Options{
			// Name of adapter. Default is "memory".
			Adapter: "memory",
			// Adapter configuration, it's corresponding to adapter.
			AdapterConfig: "",
			// GC interval time in seconds. Default is 60.
			Interval: 60,
			// Configuration section name. Default is "cache".
			Section: "cache",
		}))*/

	m.Post("/login", bind(UserLogin{}), myLoginHandler)
	m.Post("/logout", myLogoutHandler)

	NewAPICfgImportExport(m)

	NewAPIRtAgent(m)
	NewAPIRtKapFilter(m)      //Webservice for alert filtering
	NewAPIRtKapProxy(m)       //Kapacitor proxy
	NewAPIRtKapacitor(m)      //Kapacitor tasks
	NewAPIRtAlertEventHist(m) //Alert Events History
	NewAPIRtAlertEvent(m)     //Alert Events
	//Catalog
	NewAPICfgIfxServer(m) //Influx Servers
	NewAPICfgIfxDB(m)     //Influx Databases
	//Servers
	NewAPICfgKapacitor(m)      //Kapacitor URL's
	NewAPICfgIfxMeasurement(m) //Influx Measurments
	NewAPICfgEndpoint(m)       //Endpoint list
	//Config

	NewAPICfgAlertID(m)      //Alert Admin
	NewAPICfgProduct(m)      //Product Admin
	NewAPICfgProductGroup(m) //Product Group Admin
	NewAPICfgTemplate(m)     //Alert Template
	NewAPICfgRangeTime(m)    //RangeTime
	NewAPICfgDeviceStat(m)   //Device Stats
	NewAPICfgOperation(m)    //Operation

	log.Printf("Server is running on localhost:%d...", port)
	go startCleanAlertsProc()
	httpServer := fmt.Sprintf("0.0.0.0:%d", port)
	err := http.ListenAndServe(httpServer, m)
	if err != nil {
		log.Errorf("Error starting HTTP server: %s", err)
	}
}

/****************/
/*LOGIN
/****************/

func myLoginHandler(ctx *Context, user UserLogin) {
	//fmt.Printf("USER LOGIN: USER: +%#v (Config: %#v)", user, confHTTP)
	if user.UserName == confHTTP.AdminUser && user.Password == confHTTP.AdminPassword {
		ctx.SignedInUser = user.UserName
		ctx.IsSignedIn = true
		ctx.Session.Set(SessKeyUserID, user.UserName)
		kapa.SetLogger(log)
		kapa.SetPort(strconv.Itoa(confHTTP.Port))
		log.Println("Admin login OK")
		ctx.JSON(200, cookie)
	} else {
		log.Println("Admin login ERROR")
		ctx.JSON(400, "ERROR user or password not match")
	}
}

func myLogoutHandler(ctx *Context) {
	log.Printf("USER LOGOUT: USER: +%#v ", ctx.SignedInUser)
	ctx.Session.Destory(ctx)
	//ctx.Redirect("/login")
}

func startCleanAlertsProc() {
	log.Debugf("startCleanAlertsProc. Init...")
	d, err := time.ParseDuration(agent.MainConfig.Alerting.CleanPeriod)
	if err != nil {
		log.Warnf("startCleanAlertsProc. Error on ParseDuration with %s, using 5m. Error: %s", agent.MainConfig.Alerting.CleanPeriod, err)
		d, err = time.ParseDuration("5m")
	}
	t := time.NewTicker(d)
	for {
		log.Debugf("startCleanAlertsProc. Starting clean process again...")
		cleanAlertEvents()

	LOOP:
		for {
			select {
			case <-t.C:
				log.Debugf("startCleanAlertsProc. tick received...")
				break LOOP
			}
		}
	}
}

func cleanAlertEvents() {
	//Move previous alert events with level OK from alert_event to alert_event_hist
	filter := "level = 'OK'"
	MoveAlertEvents(filter)
}
