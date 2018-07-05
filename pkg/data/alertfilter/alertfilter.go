package alertfilter

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/influxdata/kapacitor/alert"
	"github.com/toni-moreno/resistor/pkg/config"
)

/* DeviceStatCfg current stats by device
type DeviceStatCfg struct {
	ID             int64  `xorm:"'id' pk autoincr"`
	OrderID          int64  `xorm:"orderid"`
	DeviceID       string `xorm:"deviceid" binding:"Required"`
	AlertID        string `xorm:"alertid" binding:"Required"`
	ProductID      string `xorm:"productid" binding:"Required"`
	ExceptionID      int64  `xorm:"exceptionid"`
	Active         bool   `xorm:"active"`
	BaseLine       string `xorm:"baseline"`
	FilterTagKey   string `xorm:"filterTagKey"`
	FilterTagValue string `xorm:"filterTagValue"`
	Description    string `xorm:"description"`
}

[Device] - > [Product] -> [BaseLine] ->[AlertID ] -> [ id, orderid - Excid - Active FilterTagKey monFilterTagValue ]

type MonStat struct {
	monExc            int64
	monActive         bool
	monLinia          string
	monFilterTagKey   string
	monFilterTagValue string
}

*/

var (
	log        *logrus.Logger
	confDir    string              //Needed to get File Filters data
	dbc        *config.DatabaseCfg //Needed to get Custom Filter  data
	DevStatsDB map[string]map[string]map[string]map[string][]*AlertStat
	mutex      sync.RWMutex
)

// SetConfDir  enable load File Filters from anywhere in the our FS.
func SetConfDir(dir string) {
	confDir = dir
}

// SetDB load database config to load data if needed (used in filters)
func SetDB(db *config.DatabaseCfg) {
	dbc = db
}

// SetLogger set output log
func SetLogger(l *logrus.Logger) {
	log = l
}

// AlertStat kk
type AlertStat struct {
	ID             int64
	OrderID        int64
	Active         bool
	ExceptionID    int64
	FilterTagKey   string
	FilterTagValue string
}

func getOrderedArray(ar []*AlertStat) []*AlertStat {
	return ar
}

func ReloadStats() error {
	defer mutex.Unlock()
	mutex.Lock()
	DevStatsDB = make(map[string]map[string]map[string]map[string][]*AlertStat)
	dbmap, err := dbc.GetDeviceStatCfgMap("")
	if err != nil {
		log.Errorf("Error on getting devicestats config : %s", err)
		return err
	}

	//[Device] - > [Product] -> [BaseLine] ->[AlertID ] -> [ id, orderid - Excid - Active FilterTagKey monFilterTagValue ]

	//Generating the MAP
	for _, v := range dbmap {
		//Check if Device Exist
		if _, ok := DevStatsDB[v.DeviceID]; ok {
			//Check if ProductID exist
			if _, ok := DevStatsDB[v.DeviceID][v.ProductID]; ok {
				//Check if Baseline exist
				if _, ok := DevStatsDB[v.DeviceID][v.ProductID][v.BaseLine]; ok {
					if _, ok := DevStatsDB[v.DeviceID][v.ProductID][v.BaseLine][v.AlertID]; ok {
						newalert := &AlertStat{
							ID:             v.ID,
							OrderID:        v.OrderID,
							Active:         v.Active,
							ExceptionID:    v.ExceptionID,
							FilterTagKey:   v.FilterTagKey,
							FilterTagValue: v.FilterTagValue,
						}
						d := DevStatsDB[v.DeviceID][v.ProductID][v.BaseLine][v.AlertID]
						d = append(d, newalert)
						DevStatsDB[v.DeviceID][v.ProductID][v.BaseLine][v.AlertID] = d
					}
				}
			}

		}
	}
	//Once all rules has been loaded we should Reorder the rules Array MAP
	for k1, val1 := range DevStatsDB {
		//k1 = deviceID
		for k2, val2 := range val1 {
			//k2 = ProductID
			for k3, val3 := range val2 {
				//k3 = BaseLines
				for k4, val4 := range val3 {
					// k4 = AlertID array
					//val4 rules array
					reordered := getOrderedArray(val4)
					DevStatsDB[k1][k2][k3][k4] = reordered

				}
			}
		}
	}
	return nil
}

// GetDevIDFromMeasurement get the deviceID from each measurmenet
func GetDevIDFromMeasurement(name string) string {
	//
	return "host"
}

//
func ApplyProductRules(al alert.Data, DevID string, ProductRules map[string]map[string]map[string][]*AlertStat) {

}

// ProcessAlert From Kapacitors
func ProcessAlert(al alert.Data) {
	log.Debugf("Process Alert Init")
	for _, v := range al.Data.Series {
		log.Debugf("ALERT ROW: %#+v", v)
		DevIDTagName := GetDevIDFromMeasurement(v.Name)

		if AlertDevID, ok := v.Tags[DevIDTagName]; ok {
			//DevID exist trying to compare with all rules
			//Generic rules first
			if prules, ok := DevStatsDB["*"]; ok {
				ApplyProductRules(al, AlertDevID, prules)
			}
			//specific rules after

			if prules, ok := DevStatsDB[AlertDevID]; ok {
				ApplyProductRules(al, AlertDevID, prules)
				//do something here
			} else {
				log.Infof("there is not info related to the %s device", AlertDevID)
			}

		} else {
			log.Errorf("DevID not found: There is some error trying to get the deviceID from serie %+v", v)
		}

	}
	log.Debugf("Process Alert End")
}
