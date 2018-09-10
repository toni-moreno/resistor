package kapa

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	kapacitorClient "github.com/influxdata/kapacitor/client/v1"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
)

var (
	log  *logrus.Logger
	port string
)

// SetLogger set output log
func SetLogger(l *logrus.Logger) {
	log = l
}

// SetPort sets resistor port
func SetPort(p string) {
	port = p
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
			taskStatus := kapacitorClient.Disabled

			//Getting JSON vars from user input
			vars := setKapaTaskVars(dev)

			// For each Kapacitor server
			// Get Kapacitor Server Config by Kapacitor Server ID
			// Get Kapacitor Go Client by Kapacitor Server Config
			// Get link to kapacitor task
			// Create or update task into Kapacitor server
			for i := 0; i < len(devcfgarray); i++ {
				kapaServerCfg := devcfgarray[i]
				log.Debugf("Kapacitor Server ID: %+s, URL: %+s", kapaServerCfg.ID, kapaServerCfg.URL)
				kapaClient, _, _, err := GetKapaClient(*kapaServerCfg)
				if err != nil {
					log.Errorf("Error creating Kapacitor Go client for kapacitor server %s. Error: %+s", kapaServerCfg.ID, err)
					sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
				} else {
					l := kapaClient.TaskLink(dev.ID)
					t, err := kapaClient.Task(l, nil)
					if err != nil {
						log.Debugf("Kapacitor Task %s NOT found into kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
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
							Status:     taskStatus,
							//TICKscript: dev.TplData,
						})
						if err != nil {
							log.Errorf("Error creating Kapacitor Task %s for kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
							sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
						} else {
							log.Debugf("Kapacitor task %s created into kapacitor server %s.", dev.ID, kapaServerCfg.ID)
							if dev.Active == false {
								iNumLastDeployed++
							} else {
								//Kapacitor task has been created or updated Disabled
								//Enable Kapacitor task if Active=true has been selected on form
								//This is done in order to Kapacitor applies new values to task
								log.Debugf("Enabling kapacitor task")
								taskStatus = kapacitorClient.Enabled
								l := kapaClient.TaskLink(dev.ID)
								_, err := kapaClient.UpdateTask(l, kapacitorClient.UpdateTaskOptions{
									ID:         dev.ID,
									TemplateID: sTemplateID,
									Type:       taskType,
									DBRPs:      DBRPs,
									Vars:       vars,
									Status:     taskStatus,
									//TICKscript: dev.TplData,
								})
								if err != nil {
									log.Errorf("Error enabling Kapacitor Task %s for kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
									sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
								} else {
									iNumLastDeployed++
								}
							}
						}
					} else {
						_, err := kapaClient.UpdateTask(l, kapacitorClient.UpdateTaskOptions{
							ID:         dev.ID,
							TemplateID: sTemplateID,
							Type:       taskType,
							DBRPs:      DBRPs,
							Vars:       vars,
							Status:     taskStatus,
							//TICKscript: dev.TplData,
						})
						if err != nil {
							log.Errorf("Error updating Kapacitor Task %s for kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
							sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
						} else {
							log.Debugf("Kapacitor task %s updated into kapacitor server %s.", dev.ID, kapaServerCfg.ID)
							if dev.Active == false {
								iNumLastDeployed++
							} else {
								//Kapacitor task has been created or updated Disabled
								//Enable Kapacitor task if Active=true has been selected on form
								//This is done in order to Kapacitor applies new values to task
								log.Debugf("Enabling kapacitor task")
								taskStatus = kapacitorClient.Enabled
								l := kapaClient.TaskLink(dev.ID)
								_, err := kapaClient.UpdateTask(l, kapacitorClient.UpdateTaskOptions{
									ID:         dev.ID,
									TemplateID: sTemplateID,
									Type:       taskType,
									DBRPs:      DBRPs,
									Vars:       vars,
									Status:     taskStatus,
									//TICKscript: dev.TplData,
								})
								if err != nil {
									log.Errorf("Error enabling Kapacitor Task %s for kapacitor server %s. Error: %+s", dev.ID, kapaServerCfg.ID, err)
									sKapaSrvsNotOK = append(sKapaSrvsNotOK, kapaServerCfg.ID)
								} else {
									iNumLastDeployed++
								}
							}
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
// TriggerType + _2EX_ + CritDirection + _ + TrendTypeTranslated + _F + StatFunc
func getTemplateID(dev config.AlertIDCfg) string {
	sRet := "DEADMAN"
	if dev.TriggerType != "DEADMAN" {
		sTriggerType := dev.TriggerType
		sTrendType := translateTrendType(dev.TriggerType, dev.TrendType, dev.TrendSign)
		sRet = fmt.Sprintf("%s_2EX_%s_%s_F%s", sTriggerType, dev.CritDirection, sTrendType, dev.StatFunc)
	}
	log.Debugf("getTemplateID. %s.", sRet)
	return sRet
}

// translateTrendType Translates TrendType
func translateTrendType(sTriggerType string, sTrendType string, sTrendSign string) string {
	sRet := sTrendType
	if sTrendType == "relative" {
		// only for TREND
		sRet = "RTP"
		if sTrendSign == "negative" {
			sRet = "RTN"
		}
	} else { // absolute
		if sTriggerType == "TREND" {
			sRet = "ATP"
			if sTrendSign == "negative" {
				sRet = "ATN"
			}
		} else {
			sRet = "TH"
		}
	}
	log.Debugf("translateTrendType. TriggerType: %s, TrendType: %s, TrendSign: %s. Returns: %s.", sTriggerType, sTrendType, sTrendSign, sRet)
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
	vars["RESISTOR_PORT"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: port}
	//Core Settings
	vars["ID"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.ID}
	vars["ID_LINE"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.BaselineID}
	vars["ID_PRODUCT"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.ProductID}
	vars["ID_GROUP"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.AlertGroup}
	vars["ID_NUMALERT"] = kapacitorClient.Var{Type: kapacitorClient.VarInt, Value: dev.NumAlertID}
	vars["ID_INSTRUCTION"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.OperationID}
	//External Services Settings
	//Endpoint
	vars["OUT_HTTP"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: strings.Join(dev.Endpoint, ",")}
	//KapacitorID
	//Data Origin Settings
	vars["INFLUX_BD"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: getIfxDBNameByID(dev.InfluxDB)}
	vars["INFLUX_RP"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.InfluxRP}
	vars["INFLUX_MEAS"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.InfluxMeasurement}
	if dev.IsCustomExpression == true {
		vars["FIELD"] = kapacitorClient.Var{Type: kapacitorClient.VarLambda, Value: dev.Field}
	} else {
		vars["FIELD"] = kapacitorClient.Var{Type: kapacitorClient.VarLambda, Value: strconv.Quote(dev.Field)}
	}
	vars["FIELD_DESC"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.FieldDesc}
	//TagDescription
	if len(dev.InfluxFilter) > 0 {
		vars["INFLUX_FILTER"] = kapacitorClient.Var{Type: kapacitorClient.VarLambda, Value: dev.InfluxFilter}
	}
	dIntervalCheck, err := time.ParseDuration(dev.IntervalCheck)
	if err != nil {
		log.Warningf("Error parsing duration from interval check %s. 0 will be assigned. Error: %s", dev.IntervalCheck, err)
	}
	vars["INTERVAL_CHECK"] = kapacitorClient.Var{Type: kapacitorClient.VarDuration, Value: dIntervalCheck}
	//Alert Settings
	if dev.TriggerType != "DEADMAN" {
		vars["STAT_FUN"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.StatFunc}
		if dev.StatFunc == "MOVINGAVERAGE" {
			vars["EXTRA_DATA"] = kapacitorClient.Var{Type: kapacitorClient.VarInt, Value: dev.ExtraData}
		} else if dev.StatFunc == "PERCENTILE" {
			vars["EXTRA_DATA"] = kapacitorClient.Var{Type: kapacitorClient.VarFloat, Value: dev.ExtraData}
		}
		vars["CRIT_DIRECTION"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.CritDirection}
		if dev.TriggerType == "TREND" {
			dShift, err := time.ParseDuration(dev.Shift)
			if err != nil {
				log.Warningf("Error parsing duration from shift field value %s. 0 will be assigned. Error: %s", dev.Shift, err)
			}
			vars["SHIFT"] = kapacitorClient.Var{Type: kapacitorClient.VarDuration, Value: dShift}
		}
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
	vars["DEVICEID_TAG"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.ProductTag}
	vars["DEVICEID_LABEL"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.DeviceIDLabel}
	vars["EXTRA_TAG"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.ExtraTag}
	vars["EXTRA_LABEL"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.ExtraLabel}

	vars["ALERT_EXTRA_TEXT"] = kapacitorClient.Var{Type: kapacitorClient.VarString, Value: dev.AlertExtraText}
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

//GetResTemplateCfgByID Gets the TemplateCfg information stored on resistor database
//including the kapacitor servers without last deployment
func GetResTemplateCfgByID(id string) (config.TemplateCfg, error) {
	log.Debugf("GetResTemplateCfgByID. Trying to get template with id: %s.", id)
	dev, err := agent.MainConfig.Database.GetTemplateCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
	} else {
		kapaserversarray, err := GetKapaServers("")
		if err != nil {
			log.Warningf("Error getting kapacitor servers: %+s", err)
		} else {
			_, _, sKapaSrvsNotOK := GetKapaTemplate(&dev, kapaserversarray)
			dev.ServersWOLastDeployment = sKapaSrvsNotOK
			log.Debugf("GetResTemplateCfgByID. Template with id: %s has not the last version deployed on: %+v.", id, sKapaSrvsNotOK)
		}
	}
	return dev, err
}

// DeployKapaTask Deploys the task related to this alert into the kapacitor server
func DeployKapaTask(dev config.AlertIDCfg) ([]string, error) {
	log.Debugf("Trying to deploy the task: %s.", dev.ID)
	sKapaSrvsNotOK := make([]string, 0)
	var err error
	kapaserversarray, err := GetKapaServers(dev.KapacitorID)

	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_, _, sKapaSrvsNotOK = SetKapaTask(dev, kapaserversarray)
		if len(sKapaSrvsNotOK) > 0 {
			log.Warningf("Error deploying task %s. Not deployed for kapacitor server %s.", dev.ID, dev.KapacitorID)
		} else {
			log.Infof("Task %s deployed for kapacitor server %s.", dev.ID, dev.KapacitorID)
		}
	}
	return sKapaSrvsNotOK, err
}

// DeployKapaTemplate Deploys template into the kapacitor servers
func DeployKapaTemplate(dev config.TemplateCfg) ([]string, error) {
	log.Debugf("Trying to deploy the template: %s.", dev.ID)
	dev.Modified = time.Now().UTC()
	sKapaSrvsNotOK := make([]string, 0)
	kapaserversarray, err := GetKapaServersFromArray(dev.ServersWOLastDeployment)
	if err != nil {
		log.Warningf("Error getting kapacitor servers from array: %+v. Error: %+s.", dev.ServersWOLastDeployment, err)
	} else {
		_, _, sKapaSrvsNotOK = SetKapaTemplate(dev, kapaserversarray)
		if len(sKapaSrvsNotOK) > 0 {
			log.Warningf("Error deploying template %s on kapacitor servers: %+v. Not updated for kapacitor servers: %+v.", dev.ID, dev.ServersWOLastDeployment, sKapaSrvsNotOK)
		} else {
			log.Infof("Template %s succesfully deployed on kapacitor servers: %+v.", dev.ID, dev.ServersWOLastDeployment)
		}
	}
	return sKapaSrvsNotOK, err
}

// GetTemplateIDParts Gets TemplateID parts from TemplateID
// example: from: "TREND_2EX_AC_ATP_FMOVAVG" --> result: "TREND", "AC", "absolute", "positive", "MOVAVG"
// TriggerType + CritDirection + TrendTypeTranslated + StatFunc
func GetTemplateIDParts(sTemplateID string) (string, string, string, string, string) {
	sTriggerType, sCritDirection, sTrendType, sTrendSign, sStatFunc := "DEADMAN", "", "", "", ""
	if sTemplateID != "DEADMAN" {
		partsarray := strings.Split(sTemplateID, "_")
		if len(partsarray) == 5 {
			sTriggerType = partsarray[0]
			sCritDirection = partsarray[2]
			sTrendType, sTrendSign = getTrendDetails(partsarray[3])
			sStatFunc = partsarray[4][1:]
		}
	}
	log.Debugf("GetTemplateIDParts. %s, %s, %s, %s, %s, %s.", sTemplateID, sTriggerType, sCritDirection, sTrendType, sTrendSign, sStatFunc)
	return sTriggerType, sCritDirection, sTrendType, sTrendSign, sStatFunc
}

// getTrendDetails Gets trend details
// A to absolute
// R to relative
// P to positive
// N to negative
func getTrendDetails(sInput string) (string, string) {
	sAbsRel := ""
	sTrendSign := ""
	if sInput == "TH" || strings.Index(sInput, "A") == 0 {
		sAbsRel = "absolute"
	} else if strings.Index(sInput, "R") == 0 {
		sAbsRel = "relative"
	} else {
		//Threshold and Trend templates must have absolute or relative
		//Set notfound in order to make possible the deletion
		sAbsRel = "notfound"
	}
	if strings.Index(sInput, "P") == 2 {
		sTrendSign = "positive"
	} else if strings.Index(sInput, "N") == 2 {
		sTrendSign = "negative"
	}
	return sAbsRel, sTrendSign
}
