package webui

import (
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	"github.com/go-macaron/binding"
	kapacitorClient "github.com/influxdata/kapacitor/client/v1"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
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
	_, elapsed, version, err := GetKapaClient(dev)
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
	devcfgarray, err := GetKapaServers("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddKapacitor Insert new snmpdevice to de internal BBDD --pending--
func AddKapacitor(ctx *Context, dev config.KapacitorCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddKapacitorCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		// Deploy the templates from resistor database into this kapacitor server
		tplcfgarray, err := GetTemplates("")
		if err == nil {
			kapasrvsarray := []*config.KapacitorCfg{&dev}
			for _, tplcfg := range tplcfgarray {
				_, _, _ = SetKapaTemplate(*tplcfg, kapasrvsarray)
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
		log.Warningf("Error on update for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
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
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}

// GetKapaClient Gets Kapacitor Go Cient
func GetKapaClient(dev config.KapacitorCfg) (*kapacitorClient.Client, time.Duration, string, error) {
	config := kapacitorClient.Config{
		URL: dev.URL,
		//Timeout:            (time.Second * time.Duration(data.Get("timeout_seconds").(int))),
		//InsecureSkipVerify: data.Get("insecure_skip_verify").(bool),
	}

	/*
		Code to be used for authentication
		if data.Get("auth_username").(string) != "" {
			credentials := kapacitorClient.Credentials{
				Method:   kapacitorClient.UserAuthentication,
				Username: data.Get("auth_username").(string),
				Password: data.Get("auth_password").(string),
				Token:    data.Get("auth_token").(string),
			}

			err := credentials.Validate()

			if err != nil {
				return nil, fmt.Errorf("error validating credentials: %s", err)
			}

			config.Credentials = &credentials
		}
	*/

	kapaClient, err := kapacitorClient.New(config)
	if err != nil {
		log.Warningf("Error getting kapacitor go client for kapacitor server %s: Err: %s", dev.ID, err)
		return kapaClient, 0, "", err
	}
	elapsed, version, err := kapaClient.Ping()

	return kapaClient, elapsed, version, err
}

// GetKapaServers Gets Kapacitor servers array
func GetKapaServers(kapacitorid string) ([]*config.KapacitorCfg, error) {
	filter := ""
	if len(kapacitorid) > 0 {
		filter = fmt.Sprintf("id = '%s'", kapacitorid)
	}
	log.Debugf("Getting Kapacitor Servers with filter: %s.", filter)
	devcfgarray, err := agent.MainConfig.Database.GetKapacitorCfgArray(filter)
	if err != nil {
		log.Errorf("Error on getting Kapacitor Servers :%+s", err)
	}
	return devcfgarray, err
}

// GetKapaServersFromArray Gets Kapacitor servers array
func GetKapaServersFromArray(kapacitoridarray []string) ([]*config.KapacitorCfg, error) {
	filter := ""
	if len(kapacitoridarray) > 0 {
		filter = "id IN ("
		for _, kapacitorid := range kapacitoridarray {
			filter += fmt.Sprintf("'%s',", kapacitorid)
		}
		filter = filter[:len(filter)-1] + ")"
	}
	log.Debugf("Getting Kapacitor Servers with filter: %s.", filter)
	devcfgarray, err := agent.MainConfig.Database.GetKapacitorCfgArray(filter)
	if err != nil {
		log.Errorf("Error on getting Kapacitor Servers :%+s", err)
	}
	return devcfgarray, err
}

// ListKapaTemplate lists a kapacitor template
func ListKapaTemplate(cli *kapacitorClient.Client, id string) (kapacitorClient.Template, error) {
	template, err := cli.ListTemplates(&kapacitorClient.ListTemplatesOptions{Pattern: id})
	if err != nil {
		log.Errorf("Failed to list template with id %s. Error: %s", id, err)
		return kapacitorClient.Template{}, err
	}

	return template[0], nil
}

// ListKapaTemplates lists all kapacitor templates
func ListKapaTemplates(cli *kapacitorClient.Client) ([]kapacitorClient.Template, error) {
	templates, err := cli.ListTemplates(nil)
	if err != nil {
		log.Errorf("Failed to list templates. Error: %s", err)
		return nil, err
	}

	return templates, nil
}

// GetKapaTemplates Gets templates from the Kapacitor Servers. Returns true if all actions have been done without error, false elsewhere.
func GetKapaTemplates(tplcfgarray []*config.TemplateCfg, devcfgarray []*config.KapacitorCfg) bool {
	log.Debugf("GetKapaTemplates. INIT.")
	allGetsOK := true
	for _, dev := range tplcfgarray {
		_, _, sKapaSrvsNotOK := GetKapaTemplate(dev, devcfgarray)
		dev.ServersWOLastDeployment = sKapaSrvsNotOK
		if len(sKapaSrvsNotOK) > 0 {
			allGetsOK = false
		}
	}
	log.Debugf("GetKapaTemplates. END.")
	return allGetsOK
}

// GetKapaTemplate Gets template from the Kapacitor Servers.
// Returns:
//     - the number of kapacitor servers
//     - the number of kapacitor servers where the template is     deployed with the last version
//     - the list   of kapacitor servers where the template is NOT deployed with the last version
func GetKapaTemplate(dev *config.TemplateCfg, devcfgarray []*config.KapacitorCfg) (int, int, []string) {
	log.Debugf("GetKapaTemplate. Trying to get template with id: %s. Modified (UTC): %+v", dev.ID, dev.Modified.UTC())
	iNumKapaServers := len(devcfgarray)
	iNumLastDeployed := 0
	sKapaSrvsNotOK := make([]string, 0)

	// For each Kapacitor server
	// Get Kapacitor Server Config by Kapacitor Server ID
	// Get Kapacitor Go Client by Kapacitor Server Config
	// Get link to kapacitor template
	// Get template from Kapacitor server
	for i := 0; i < iNumKapaServers; i++ {
		kapaServerCfg := devcfgarray[i]
		log.Debugf("Kapacitor Server ID, URL: %+s, %s", kapaServerCfg.ID, kapaServerCfg.URL)
		kapaClient, _, _, err := GetKapaClient(*kapaServerCfg)
		if err != nil {
			log.Errorf("Error creating Kapacitor Go client for kapacitor server %s. Error: %+s", kapaServerCfg.ID, err)
			sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
		} else {
			l := kapaClient.TemplateLink(dev.ID)
			t, err := kapaClient.Template(l, nil)
			if err != nil {
				log.Errorf("Error getting Kapacitor Template %s for kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
			} else {
				log.Debugf("Kapacitor template %s found into kapacitor server %s.", dev.ID, kapaServerCfg.ID)
			}
			if t.ID == "" {
				log.Debugf("Kapacitor template %s not found into kapacitor server %s.", dev.ID, kapaServerCfg.ID)
				sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
			} else {
				log.Debugf("Kapacitor template %s found into kapacitor server %s. Modified (UTC): %+v.", dev.ID, kapaServerCfg.ID, t.Modified.UTC())
				d, _ := time.ParseDuration("10s")
				diff := math.Abs(dev.Modified.UTC().Sub(t.Modified.UTC()).Seconds())
				if diff > d.Seconds() {
					log.Debugf("GetKapaTemplate. Difference between update moments %.3f s greater than 10 s.", diff)
					sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
				} else {
					log.Debugf("GetKapaTemplate. Difference between update moments %.3f s lower than 10 s.", diff)
					iNumLastDeployed++
				}
			}
		}
	}
	log.Debugf("GetKapaTemplate. END.")
	return iNumKapaServers, iNumLastDeployed, sKapaSrvsNotOK
}

// SetKapaTemplate Creates or updates template into the Kapacitor Servers.
// Returns:
//     - the number of kapacitor servers where the template will be deployed with the last version
//     - the number of kapacitor servers where the template is      deployed with the last version
//     - the list   of kapacitor servers where the template is NOT  deployed with the last version
func SetKapaTemplate(dev config.TemplateCfg, devcfgarray []*config.KapacitorCfg) (int, int, []string) {
	log.Debugf("SetKapaTemplate. Trying to create or update template with id: %s", dev.ID)
	iNumKapaServers := len(devcfgarray)
	iNumLastDeployed := 0
	sKapaSrvsNotOK := make([]string, 0)

	taskType := kapacitorClient.StreamTask

	// For each Kapacitor server
	// Get Kapacitor Server Config by Kapacitor Server ID
	// Get Kapacitor Go Client by Kapacitor Server Config
	// Get link to kapacitor template
	// Create or update template into Kapacitor server
	for i := 0; i < iNumKapaServers; i++ {
		kapaServerCfg := devcfgarray[i]
		log.Debugf("Kapacitor Server ID, URL: %+s, %s", kapaServerCfg.ID, kapaServerCfg.URL)
		kapaClient, _, _, err := GetKapaClient(*kapaServerCfg)
		if err != nil {
			log.Errorf("Error creating Kapacitor Go client for kapacitor server %s. Error: %+s", kapaServerCfg.ID, err)
			sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
		} else {
			l := kapaClient.TemplateLink(dev.ID)
			t, err := kapaClient.Template(l, nil)
			if err != nil {
				log.Errorf("Error getting Kapacitor Template %s for kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
			} else {
				log.Debugf("Kapacitor template %s found into kapacitor server %s.", dev.ID, kapaServerCfg.ID)
			}
			if t.ID == "" {
				_, err := kapaClient.CreateTemplate(kapacitorClient.CreateTemplateOptions{
					ID:         dev.ID,
					Type:       taskType,
					TICKscript: dev.TplData,
				})
				if err != nil {
					log.Errorf("Error creating Kapacitor Template %s for kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
					sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
				} else {
					log.Debugf("Kapacitor template %s created into kapacitor server %s.", dev.ID, kapaServerCfg.ID)
					iNumLastDeployed++
				}
			} else {
				_, err := kapaClient.UpdateTemplate(l, kapacitorClient.UpdateTemplateOptions{
					ID:         dev.ID,
					Type:       taskType,
					TICKscript: dev.TplData,
				})
				if err != nil {
					log.Errorf("Error updating Kapacitor Template %s for kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
					sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
				} else {
					log.Debugf("Kapacitor template %s updated into kapacitor server %s.", dev.ID, kapaServerCfg.ID)
					iNumLastDeployed++
				}
			}
		}
	}
	log.Debugf("SetKapaTemplate. END.")
	return iNumKapaServers, iNumLastDeployed, sKapaSrvsNotOK
}

// DeleteKapaTemplate Deletes template from the Kapacitor Servers.
// Returns:
//     - the number of kapacitor servers
//     - the number of kapacitor servers where the template is     deleted
//     - the list   of kapacitor servers where the template is NOT deleted
func DeleteKapaTemplate(id string, devcfgarray []*config.KapacitorCfg) (int, int, []string) {
	log.Debugf("DeleteKapaTemplate. Trying to delete template with id: %s", id)
	iNumKapaServers := len(devcfgarray)
	iNumDeleted := 0
	sKapaSrvsNotOK := make([]string, 0)

	// For each Kapacitor server
	// Get Kapacitor Server Config by Kapacitor Server ID
	// Get Kapacitor Go Client by Kapacitor Server Config
	// Get link to kapacitor template
	// Delete template from Kapacitor server
	for i := 0; i < iNumKapaServers; i++ {
		kapaServerCfg := devcfgarray[i]
		log.Debugf("Kapacitor Server ID, URL: %+s, %s", kapaServerCfg.ID, kapaServerCfg.URL)
		kapaClient, _, _, err := GetKapaClient(*kapaServerCfg)
		if err != nil {
			log.Errorf("Error creating Kapacitor Go client for kapacitor server %s. Error: %+s", kapaServerCfg.ID, err)
			sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
		} else {
			l := kapaClient.TemplateLink(id)
			err = kapaClient.DeleteTemplate(l)
			if err != nil {
				log.Errorf("Error deleting Kapacitor Template %s from kapacitor server %s. Error: %+s", id, kapaServerCfg.ID, err)
				sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
			} else {
				log.Debugf("Kapacitor template %s deleted from kapacitor server %s.", id, kapaServerCfg.ID)
				iNumDeleted++
			}
		}
	}
	log.Debugf("DeleteKapaTemplate. END.")
	return iNumKapaServers, iNumDeleted, sKapaSrvsNotOK
}

// ListKapaTask lists a kapacitor task
func ListKapaTask(cli *kapacitorClient.Client, id string) (kapacitorClient.Task, error) {
	task, err := cli.ListTasks(&kapacitorClient.ListTasksOptions{Pattern: id})
	if err != nil {
		log.Errorf("Failed to list task with id %s. Error: %s", id, err)
		return kapacitorClient.Task{}, err
	}

	return task[0], nil
}

// ListKapaTasks lists all kapacitor tasks
func ListKapaTasks(cli *kapacitorClient.Client) ([]kapacitorClient.Task, error) {
	tasks, err := cli.ListTasks(nil)
	if err != nil {
		log.Errorf("Failed to list tasks. Error: %s", err)
		return nil, err
	}

	return tasks, nil
}

// GetKapaTasks Gets tasks from the Kapacitor Server related. Returns true if all actions have been done without error, false elsewhere.
func GetKapaTasks(tplcfgarray []*config.AlertIDCfg) bool {
	log.Debugf("GetKapaTasks. INIT.")
	allGetsOK := true
	for _, dev := range tplcfgarray {
		_, _, sKapaSrvsNotOK := GetKapaTask(dev)
		dev.ServersWOLastDeployment = sKapaSrvsNotOK
		if len(sKapaSrvsNotOK) > 0 {
			allGetsOK = false
		}
	}
	log.Debugf("GetKapaTasks. END.")
	return allGetsOK
}

// GetKapaTask Gets task from the Kapacitor Server related to this alert.
// Returns:
//     - the number of kapacitor servers
//     - the number of kapacitor servers where the task is     deployed with the last version
//     - the list   of kapacitor servers where the task is NOT deployed with the last version
func GetKapaTask(dev *config.AlertIDCfg) (int, int, []string) {
	log.Debugf("GetKapaTask. Trying to get task with id: %s. Modified (UTC): %+v", dev.ID, dev.Modified.UTC())
	iNumKapaServers := 0
	iNumLastDeployed := 0
	sKapaSrvsNotOK := make([]string, 0)
	devcfgarray, err := GetKapaServers(dev.KapacitorID)
	if err != nil {
		log.Warningf("Error getting kapacitor server for alert: %s. Error: %+s", dev.ID, err)
	} else {
		iNumKapaServers = len(devcfgarray)
		// For each Kapacitor server
		// Get Kapacitor Server Config by Kapacitor Server ID
		// Get Kapacitor Go Client by Kapacitor Server Config
		// Get link to kapacitor task
		// Get task from Kapacitor server
		for i := 0; i < len(devcfgarray); i++ {
			kapaServerCfg := devcfgarray[i]
			log.Debugf("Kapacitor Server ID, URL: %+s, %s", kapaServerCfg.ID, kapaServerCfg.URL)
			kapaClient, _, _, err := GetKapaClient(*kapaServerCfg)
			if err != nil {
				log.Errorf("Error creating Kapacitor Go client for kapacitor server %s. Error: %+s", kapaServerCfg.ID, err)
				sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
			} else {
				l := kapaClient.TaskLink(dev.ID)
				t, err := kapaClient.Task(l, nil)
				if err != nil {
					log.Errorf("Error getting Kapacitor Task %s for kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
				} else {
					log.Debugf("Kapacitor task %s found into kapacitor server %s.", dev.ID, kapaServerCfg.ID)
				}
				if t.ID == "" {
					log.Debugf("Kapacitor task %s not found into kapacitor server %s.", dev.ID, kapaServerCfg.ID)
					sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
				} else {
					log.Debugf("Kapacitor task %s found into kapacitor server %s. Modified (UTC): %+v.", dev.ID, kapaServerCfg.ID, t.Modified.UTC())
					d, _ := time.ParseDuration("10s")
					diff := math.Abs(t.Modified.UTC().Sub(dev.Modified.UTC()).Seconds())
					if diff > d.Seconds() {
						log.Debugf("GetKapaTask. Difference between update moments %.3f s greater than 10 s.", diff)
						sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
					} else {
						log.Debugf("GetKapaTask. Difference between update moments %.3f s lower than 10 s.", diff)
						iNumLastDeployed++
					}
				}
			}
		}
	}

	log.Debugf("GetKapaTask. END.")
	return iNumKapaServers, iNumLastDeployed, sKapaSrvsNotOK
}

// SetKapaTask Creates or updates task into the Kapacitor Servers
// For this version, the array of Kapacitor Servers only contains 1 server
// First of all, a checking is done to ensure the template used for the task has the last deployment into all kapacitor servers
// Returns:
//     - the number of kapacitor servers where the task will be deployed
//     - the number of kapacitor servers where the task is      deployed with the last version
//     - the list   of kapacitor servers where the task is NOT  deployed with the last version
func SetKapaTask(dev config.AlertIDCfg, devcfgarray []*config.KapacitorCfg) (int, int, []string) {
	log.Debugf("SetKapaTask. Trying to create or update task with id: %s and info: %+v into kapacitor servers: %+v", dev.ID, dev, devcfgarray)
	iNumKapaServers := len(devcfgarray)
	iNumLastDeployed := 0
	sKapaSrvsNotOK := make([]string, 0)

	// Ensure the template used for the task has the last deployment into all kapacitor servers
	sTemplateID := getTemplateID(dev)
	devTpl, err := GetResTemplateCfgByID(sTemplateID)
	if err != nil {
		log.Warningf("Error getting template %s from resistor database. Error: %+s", sTemplateID, err)
		sKapaSrvsNotOK = getKapaCfgIDArray(devcfgarray)
	} else {
		if len(devTpl.ServersWOLastDeployment) > 0 {
			log.Warningf("Template %s has not the last deployment for kapacitor servers %+v.", sTemplateID, devTpl.ServersWOLastDeployment)
			sKapaSrvsNotOK = getKapaCfgIDArray(devcfgarray)
		} else {
			//Getting DBRPs
			DBRPs := make([]kapacitorClient.DBRP, 1)
			DBRPs[0].Database = getIfxDBNameByID(dev.InfluxDB)
			DBRPs[0].RetentionPolicy = dev.InfluxRP

			taskType := kapacitorClient.StreamTask

			//Getting JSON vars from user input
			vars := setKapaTaskVars(dev)

			// For each Kapacitor server
			// Get Kapacitor Server Config by Kapacitor Server ID
			// Get Kapacitor Go Client by Kapacitor Server Config
			// Get link to kapacitor task
			// Create or update task into Kapacitor server
			for i := 0; i < len(devcfgarray); i++ {
				kapaServerCfg := devcfgarray[i]
				log.Debugf("Kapacitor Server ID, URL: %+s, %s", kapaServerCfg.ID, kapaServerCfg.URL)
				kapaClient, _, _, err := GetKapaClient(*kapaServerCfg)
				if err != nil {
					log.Errorf("Error creating Kapacitor Go client for kapacitor server %s. Error: %+s", kapaServerCfg.ID, err)
					sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
				} else {
					l := kapaClient.TaskLink(dev.ID)
					t, err := kapaClient.Task(l, nil)
					if err != nil {
						log.Errorf("Error getting Kapacitor Task %s for kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
					} else {
						log.Debugf("Kapacitor task %s found into kapacitor server %s.", dev.ID, kapaServerCfg.ID)
					}
					if t.ID == "" {
						_, err := kapaClient.CreateTask(kapacitorClient.CreateTaskOptions{
							ID:         dev.ID,
							TemplateID: sTemplateID,
							Type:       taskType,
							DBRPs:      DBRPs,
							Vars:       vars,
							Status:     kapacitorClient.Enabled,
							//TICKscript: dev.TplData,
						})
						if err != nil {
							log.Errorf("Error creating Kapacitor Task %s for kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
							sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
						} else {
							log.Debugf("Kapacitor task %s created into kapacitor server %s.", dev.ID, kapaServerCfg.ID)
							iNumLastDeployed++
						}
					} else {
						_, err := kapaClient.UpdateTask(l, kapacitorClient.UpdateTaskOptions{
							ID:         dev.ID,
							TemplateID: sTemplateID,
							Type:       taskType,
							DBRPs:      DBRPs,
							Vars:       vars,
							Status:     kapacitorClient.Enabled,
							//TICKscript: dev.TplData,
						})
						if err != nil {
							log.Errorf("Error updating Kapacitor Task %s for kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
							sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
						} else {
							log.Debugf("Kapacitor task %s updated into kapacitor server %s.", dev.ID, kapaServerCfg.ID)
							iNumLastDeployed++
						}
					}
				}
			}
		}
	}

	log.Debugf("SetKapaTask. END.")
	return iNumKapaServers, iNumLastDeployed, sKapaSrvsNotOK
}

func getKapaCfgIDArray(devcfgarray []*config.KapacitorCfg) []string {
	idarray := make([]string, 0)
	for i := 0; i < len(devcfgarray); i++ {
		cfg := devcfgarray[i]
		idarray = append(idarray, cfg.ID)
	}
	return idarray
}

// getTemplateID Gets TemplateID from AlertIDCfg
// example: "THRESHOLD_2EX_AC_TH_FMOVAVG"
// TrigerType + _2EX_ + CritDirection + _ + ThresholdTypeTranslated + _F + StatFunc
func getTemplateID(dev config.AlertIDCfg) string {
	sRet := "DEADMAN"
	if dev.TrigerType != "DEADMAN" {
		sTriggerType := dev.TrigerType
		sThresholdType := translateThresholdType(dev.TrigerType, dev.ThresholdType, dev.TrendSign)
		sRet = fmt.Sprintf("%s_2EX_%s_%s_F%s", sTriggerType, dev.CritDirection, sThresholdType, dev.StatFunc)
	}
	log.Debugf("getTemplateID. %s.", sRet)
	return sRet
}

// translateThresholdType Translates ThresholdType
func translateThresholdType(sTriggerType string, sThresholdType string, sTrendSign string) string {
	sRet := sThresholdType
	if sThresholdType == "relative" {
		// only for TREND
		sRet = "TRP"
		if sTrendSign == "negative" {
			sRet = "TRN"
		}
	} else { // absolute
		if sTriggerType == "TREND" {
			sRet = "TAP"
			if sTrendSign == "negative" {
				sRet = "TAN"
			}
		} else {
			sRet = "TH"
		}
	}
	log.Debugf("translateThresholdType. TriggerType: %s, ThresholdType: %s, TrendSign: %s. Returns: %s.", sTriggerType, sThresholdType, sTrendSign, sRet)
	return sRet
}

func getOwnIP() string {
	sRet := "7.116.100.107"
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Errorf("getOwnIP. Error dialing udp: %+v.", err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	log.Debugf("getOwnIP. localAddr: %+v.", localAddr)
	sRet = localAddr.IP.String()
	log.Debugf("getOwnIP. %s", sRet)
	return sRet
}

// setKapaTaskVars Creates Vars for the Kapacitor task
func setKapaTaskVars(dev config.AlertIDCfg) kapacitorClient.Vars {
	//Getting JSON vars from user input
	vars := make(kapacitorClient.Vars)

	vars["RESISTOR_IP"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: getOwnIP()}
	vars["RESISTOR_PORT"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: "8090"}
	//Core Settings
	vars["ID"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.ID}
	vars["ID_LINE"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.BaselineID}
	vars["ID_PRODUCT"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.ProductID}
	vars["ID_GROUP"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.GroupID}
	vars["ID_NUMALERT"] = kapacitorClient.Var{Type: kapacitorClient.VarInt, Value: dev.NumAlertID}
	vars["ID_INSTRUCTION"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.OperationID}
	//External Services Settings
	//OutHTTP
	vars["OUT_HTTP"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: strings.Join(dev.OutHTTP, ",")}
	//KapacitorID
	//Data Origin Settings
	vars["INFLUX_BD"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: getIfxDBNameByID(dev.InfluxDB)}
	vars["INFLUX_RP"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.InfluxRP}
	vars["INFLUX_MEAS"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.InfluxMeasurement}
	vars["FIELD"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.Field}
	//TagDescription
	//Don't use InfluxFilter
	//vars["INFLUX_FILTER"] = kapacitorClient.Var{Type: kapacitorClient.VarLambda, Value: dev.InfluxFilter}
	dIntervalCheck, err := time.ParseDuration(dev.IntervalCheck)
	if err != nil {
		log.Warningf("Error parsing duration from interval check %s. 0 will be assigned. Error: %s", dev.IntervalCheck, err)
	}
	vars["INTERVAL_CHECK"] = kapacitorClient.Var{Type: kapacitorClient.VarDuration, Value: dIntervalCheck}
	//Alert Settings
	//vars["TIPO_TRIGUER"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.TrigerType}
	if dev.TrigerType != "DEADMAN" {
		vars["STAT_FUN"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.StatFunc}
		vars["CRIT_DIRECTION"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.CritDirection}
		if dev.TrigerType == "TREND" {
			dShift, err := time.ParseDuration(dev.Shift)
			if err != nil {
				log.Warningf("Error parsing duration from shift field value %s. 0 will be assigned. Error: %s", dev.Shift, err)
			}
			vars["SHIFT"] = kapacitorClient.Var{Type: kapacitorClient.VarDuration, Value: dShift}
		}
		//ThresholdType NOT USED !!!
		vars["TH_CRIT_DEF"] = kapacitorClient.Var{Type: kapacitorClient.VarFloat, Value: dev.ThCritDef}
		vars["TH_CRIT_EX1"] = kapacitorClient.Var{Type: kapacitorClient.VarFloat, Value: dev.ThCritEx1}
		vars["TH_CRIT_EX2"] = kapacitorClient.Var{Type: kapacitorClient.VarFloat, Value: dev.ThCritEx2}
		min, max, weekdays := 0, 23, "0123456"
		if len(dev.ThCritRangeTimeID) > 0 {
			min, max, weekdays = getRangeTimeInfo(dev.ThCritRangeTimeID)
			vars["TH_CRIT_MIN_HOUR"] = kapacitorClient.Var{Type: kapacitorClient.VarInt, Value: min}
			vars["TH_CRIT_MAX_HOUR"] = kapacitorClient.Var{Type: kapacitorClient.VarInt, Value: max}
			vars["DAY_WEEK_CRIT"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: weekdays}
		}
		vars["TH_WARN_DEF"] = kapacitorClient.Var{Type: kapacitorClient.VarFloat, Value: dev.ThWarnDef}
		vars["TH_WARN_EX1"] = kapacitorClient.Var{Type: kapacitorClient.VarFloat, Value: dev.ThWarnEx1}
		vars["TH_WARN_EX2"] = kapacitorClient.Var{Type: kapacitorClient.VarFloat, Value: dev.ThWarnEx2}
		if len(dev.ThWarnRangeTimeID) > 0 {
			min, max, weekdays = getRangeTimeInfo(dev.ThWarnRangeTimeID)
			vars["TH_WARN_MIN_HOUR"] = kapacitorClient.Var{Type: kapacitorClient.VarInt, Value: min}
			vars["TH_WARN_MAX_HOUR"] = kapacitorClient.Var{Type: kapacitorClient.VarInt, Value: max}
			vars["DAY_WEEK_WARN"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: weekdays}
		}
		vars["TH_INFO_DEF"] = kapacitorClient.Var{Type: kapacitorClient.VarFloat, Value: dev.ThInfoDef}
		vars["TH_INFO_EX1"] = kapacitorClient.Var{Type: kapacitorClient.VarFloat, Value: dev.ThInfoEx1}
		vars["TH_INFO_EX2"] = kapacitorClient.Var{Type: kapacitorClient.VarFloat, Value: dev.ThInfoEx2}
		if len(dev.ThInfoRangeTimeID) > 0 {
			min, max, weekdays = getRangeTimeInfo(dev.ThInfoRangeTimeID)
			vars["TH_INFO_MIN_HOUR"] = kapacitorClient.Var{Type: kapacitorClient.VarInt, Value: min}
			vars["TH_INFO_MAX_HOUR"] = kapacitorClient.Var{Type: kapacitorClient.VarInt, Value: max}
			vars["DAY_WEEK_INFO"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: weekdays}
		}
	}
	//Extra Settings
	vars["GRAFANA_SERVER"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.GrafanaServer}
	vars["GRAFANA_DASH_LABEL"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.GrafanaDashLabel}
	vars["GRAFANA_DASH_PANELID"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.GrafanaDashPanelID}
	vars["DEVICEID_TAG"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.DeviceIDTag}
	vars["DEVICEID_LABEL"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.DeviceIDLabel}
	vars["EXTRA_TAG"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.ExtraTag}
	vars["EXTRA_LABEL"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.ExtraLabel}

	//ALERT_EXTRA_TEXT corresponds to Description field on form???
	vars["ALERT_EXTRA_TEXT"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: ""}
	//vars["FIELD_DEFAULT"] = kapacitorClient.Var{Type: kapacitorClient.VarFloat, Value: ""}
	/*
		vars["MOV_AVG_POINTS"] = kapacitorClient.Var{Type: kapacitorClient.VarInt, Value: ""}
	*/

	return vars
}

// getIfxDBNameByID
func getIfxDBNameByID(id int64) string {
	name := ""
	dev, err := agent.MainConfig.Database.GetIfxDBCfgByID(id)
	if err != nil {
		log.Warningf("Error getting influx db name for id %d. Empty string will be returned. Error: %s", id, err)
	} else {
		name = dev.Name
	}
	return name
}

// getRangeTimeInfo
func getRangeTimeInfo(id string) (int, int, string) {
	min := 0
	max := 23
	weekdays := "0123456"
	dev, err := agent.MainConfig.Database.GetRangeTimeCfgByID(id)
	if err != nil {
		log.Warningf("Error getting range time info for id %s. 0, 23, 0123456 will be returned. Error: %s", id, err)
	} else {
		min = dev.MinHour
		max = dev.MaxHour
		weekdays = dev.WeekDays
	}
	return min, max, weekdays
}

// DeleteKapaTask Deletes task from the Kapacitor Servers
// Returns:
//     - the number of kapacitor servers
//     - the number of kapacitor servers where the task is     deleted
//     - the list   of kapacitor servers where the task is NOT deleted
func DeleteKapaTask(id string, devcfgarray []*config.KapacitorCfg) (int, int, []string) {
	log.Debugf("DeleteKapaTask. Trying to delete task with id: %s", id)
	iNumKapaServers := len(devcfgarray)
	iNumDeleted := 0
	sKapaSrvsNotOK := make([]string, 0)

	// For each Kapacitor server
	// Get Kapacitor Server Config by Kapacitor Server ID
	// Get Kapacitor Go Client by Kapacitor Server Config
	// Get link to kapacitor task
	// Delete task from Kapacitor server
	for i := 0; i < len(devcfgarray); i++ {
		kapaServerCfg := devcfgarray[i]
		log.Debugf("Kapacitor Server ID, URL: %+s, %s", kapaServerCfg.ID, kapaServerCfg.URL)
		kapaClient, _, _, err := GetKapaClient(*kapaServerCfg)
		if err != nil {
			log.Errorf("Error creating Kapacitor Go client for kapacitor server %s. Error: %+s", kapaServerCfg.ID, err)
			sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
		} else {
			l := kapaClient.TaskLink(id)
			t, err := kapaClient.Task(l, nil)
			if err != nil {
				log.Errorf("Error getting Kapacitor Task %s for kapacitor server %s. Error: %+s", id, kapaServerCfg.ID, err)
			} else {
				log.Debugf("Kapacitor task %s found into kapacitor server %s.", id, kapaServerCfg.ID)
			}
			if t.ID == "" {
				log.Debugf("Kapacitor task %s does not exist on kapacitor server %s.", id, kapaServerCfg.ID)
			} else {
				err = kapaClient.DeleteTask(l)
				if err != nil {
					log.Errorf("Error deleting Kapacitor Task %s from kapacitor server %s. Error: %+s", id, kapaServerCfg.ID, err)
					sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
				} else {
					log.Debugf("Kapacitor task %s deleted from kapacitor server %s.", id, kapaServerCfg.ID)
					iNumDeleted++
				}
			}
		}
	}
	log.Debugf("DeleteKapaTask. END. iNumKapaServers:%d, iNumDeleted:%d, sKapaSrvsNotOK:%+v.", iNumKapaServers, iNumDeleted, sKapaSrvsNotOK)
	return iNumKapaServers, iNumDeleted, sKapaSrvsNotOK
}
