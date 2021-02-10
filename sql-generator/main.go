package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"gopkg.in/yaml.v2"
)

func loadYml() Config {
	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	c := loadYml()

	assetNum := 50
	paramNum := 200

	f, err := os.OpenFile(c.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	sql := addAssetType(c.TenantId)
	if _, err := f.WriteString(sql); err != nil {
		log.Println(err)
	}

	sql = addParamType(c.TenantId)
	if _, err := f.WriteString(sql); err != nil {
		log.Println(err)
	}

	sql = updateParamType()
	if _, err := f.WriteString(sql); err != nil {
		log.Println(err)
	}

	idx := 1
	for i := 1; i <= paramNum; i++ {
		paramId := uuid.New().String()
		paramName := "parameter" + strconv.Itoa(i)
		sql := addParameter(paramId, paramName, c.ParamGroup, i)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		var props string
		if idx == 1 {
			props = `{"type": "default", "lowerLimit": 900, "targetValue": 1020, "upperLimit": 1200}`
		} else if idx == 2 {
			props = `{"type": "default", "lowerLimit": 1000, "targetValue": 2000, "upperLimit": 3000}`
		} else if idx == 3 {
			props = `{"type": "default", "lowerLimit": 800, "targetValue": 1800, "upperLimit": 2800}`
		} else if idx == 4 {
			props = `{"type": "default", "lowerLimit": 360, "targetValue": 470, "upperLimit": 550}`
		} else if idx == 5 {
			props = `{"type": "default", "lowerLimit": 1.1, "targetValue": 2.7, "upperLimit": 4.5}`
		} else if idx == 6 {
			props = `{"type": "default", "lowerLimit": 1600, "targetValue": 1770, "upperLimit": 1970}`
		} else if idx == 7 {
			props = `{"type": "default", "lowerLimit": 0.9, "targetValue": 1.02, "upperLimit": 1.15}`
			idx = 0
		}

		idx++

		sql = addParamModelValue(uuid.New().String(), paramId, props)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}
	}

	deviceId := uuid.New().String()
	sql = addDevice(deviceId, c.TenantId, "KMC-01")
	if _, err := f.WriteString(sql); err != nil {
		log.Println(err)
	}

	for i := 1; i <= assetNum; i++ {
		assetId := uuid.New().String()
		assetName := "Pump" + strconv.Itoa(i)
		sensorName := "S" + strconv.Itoa(i)

		sql := queryForAsset(assetId, c.CatalogId, assetName, c.TemplateId, c.TypeId)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = linkAssetWithCatalog(c.CatalogId, assetId)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addNestEgg(assetId, c.NestId)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sensorId := uuid.New().String()
		sql = addSensor(sensorId, assetId, sensorName, deviceId)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = linkSensorParamGroup(assetId, sensorId, c.ParamGroup)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addSensorStatus(sensorId)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		fdc := uuid.New().String()
		sql = addFDCModel(fdc, c.TenantId, assetId)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addModelParam(uuid.New().String(), fdc, "parameter2", 0)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addModelHistory(uuid.New().String(), assetId, fdc, "FDC")
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addStatus(assetId, fdc)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		ruleId := uuid.New().String()
		sql = addModelAlarmRule(ruleId, assetId, fdc)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		ruleItemId := uuid.New().String()
		sql = addModelAlarmRuleItem(ruleItemId, ruleId)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addModelAlarmRuleProperty(uuid.New().String(), ruleItemId)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		//v1.0
		//baeProps := `[{\"seq\":\"1\",\"name\":\"modelInfo\",\"description\":\"\",\"inputType\":\"TEXTAREA\",\"dataType\":\"TRAINING\",\"referenceType\":\"\",\"defaultValue\":\"\",\"value\":\"{\\"nIn\\":5,\\"nOut\\":5,\\"nLayers\\":4,\\"hiddenLayerSizes\\":[5,3,3,5],\\"denseLayers\\":[{\\"nIn\\":5,\\"nOut\\":5,\\"b\\":[-0.02859377596374497,0.01935353298078354,-0.01121440863629936,0.007174609861919613,0.006718840769984184],\\"seed\\":-1517918040,\\"activationCode\\":\\"ReLU\\",\\"w\\":[[-0.3034549012500158,0.9203633878883828,0.7264737783510763,-0.36187772012855196,-0.5510869770620085],[-0.32990154400916877,-0.11119461410219257,0.3094185301277254,-0.7704397505928792,0.5755094931522647],[0.9610127368873562,-0.9466482138146094,0.6695501825922757,-0.0469426312278216,-0.3628913256422779],[0.41428371531432434,0.9144008379000774,-0.33237087719851777,-0.902869420943007,0.11569438161606906],[-0.8663917266246136,0.9833746722297892,-0.6724094797930992,-0.5414539637983147,0.4267680964603218]]},{\\"nIn\\":5,\\"nOut\\":3,\\"b\\":[0.005664849599676268,-0.018478581775295756,-0.04001879441250058],\\"seed\\":-1517918040,\\"activationCode\\":\\"ReLU\\",\\"w\\":[[-0.2982748371024824,0.9503858485213565,0.7450510197818194,-0.3793771625882676,-0.5797979061246302],[-0.34054498701499036,-0.12028698168230333,0.25866175088218113,-0.7517300522941412,0.5774674767090235],[0.9191219428929299,-0.9600347720482035,0.6407851463119885,-0.09114115675158391,-0.4032418543999576]]},{\\"nIn\\":3,\\"nOut\\":3,\\"b\\":[-0.04542065131143342,0.0,-0.0031379027068379074],\\"seed\\":-1517918040,\\"activationCode\\":\\"ReLU\\",\\"w\\":[[-0.4192430807739384,1.2145943831471693,0.8840944191936619],[-0.4908006778577193,-0.748051684382309,-0.4267662180298355],[-0.14689928304360564,0.3900394721061242,-0.9665645642201744]]},{\\"nIn\\":3,\\"nOut\\":5,\\"b\\":[0.004552716365142011,0.016381667140860943,0.004712876569661641,-0.021070743928128886,-0.026903360560060943],\\"seed\\":-1517918040,\\"activationCode\\":\\"Linear\\",\\"w\\":[[-0.3656227561185739,1.2246905186284998,0.9394040561359984],[-0.4799202859490072,-0.748051684382309,-0.4254887678392957],[-0.1324953494513111,0.39075553356182935,-0.9683373454711743],[0.7457834129742265,1.2588578403475645,-1.2293426686066573],[0.8486748227266485,-0.06246475936545304,-0.47424236877272036]]}],\\"mean\\":[0.9089397971666666,1020.3972950555559,1.8538348170740744,2121.146388,1787.8770603518515],\\"std\\":[0.035528861862199654,5.2898207605052745,0.012229934331297439,118.21647258590939,17.640434189614375]}\"},{\"seq\":\"2\",\"name\":\"normalizer\",\"description\":\"\",\"inputType\":\"TEXT\",\"dataType\":\"NUMBER\",\"referenceType\":\"\",\"defaultValue\":\"\",\"value\":\"3\"}]`
		//v2.0
		baeProps := `[
			{
			  \"attributeId\": 82,
			  \"value\": \"{\\"nIn\\":5,\\"nOut\\":5,\\"nLayers\\":4,\\"hiddenLayerSizes\\":[5,3,3,5],\\"denseLayers\\":[{\\"nIn\\":5,\\"nOut\\":5,\\"b\\":[-0.02859377596374497,0.01935353298078354,-0.01121440863629936,0.007174609861919613,0.006718840769984184],\\"seed\\":-1517918040,\\"activationCode\\":\\"ReLU\\",\\"w\\":[[-0.3034549012500158,0.9203633878883828,0.7264737783510763,-0.36187772012855196,-0.5510869770620085],[-0.32990154400916877,-0.11119461410219257,0.3094185301277254,-0.7704397505928792,0.5755094931522647],[0.9610127368873562,-0.9466482138146094,0.6695501825922757,-0.0469426312278216,-0.3628913256422779],[0.41428371531432434,0.9144008379000774,-0.33237087719851777,-0.902869420943007,0.11569438161606906],[-0.8663917266246136,0.9833746722297892,-0.6724094797930992,-0.5414539637983147,0.4267680964603218]]},{\\"nIn\\":5,\\"nOut\\":3,\\"b\\":[0.005664849599676268,-0.018478581775295756,-0.04001879441250058],\\"seed\\":-1517918040,\\"activationCode\\":\\"ReLU\\",\\"w\\":[[-0.2982748371024824,0.9503858485213565,0.7450510197818194,-0.3793771625882676,-0.5797979061246302],[-0.34054498701499036,-0.12028698168230333,0.25866175088218113,-0.7517300522941412,0.5774674767090235],[0.9191219428929299,-0.9600347720482035,0.6407851463119885,-0.09114115675158391,-0.4032418543999576]]},{\\"nIn\\":3,\\"nOut\\":3,\\"b\\":[-0.04542065131143342,0.0,-0.0031379027068379074],\\"seed\\":-1517918040,\\"activationCode\\":\\"ReLU\\",\\"w\\":[[-0.4192430807739384,1.2145943831471693,0.8840944191936619],[-0.4908006778577193,-0.748051684382309,-0.4267662180298355],[-0.14689928304360564,0.3900394721061242,-0.9665645642201744]]},{\\"nIn\\":3,\\"nOut\\":5,\\"b\\":[0.004552716365142011,0.016381667140860943,0.004712876569661641,-0.021070743928128886,-0.026903360560060943],\\"seed\\":-1517918040,\\"activationCode\\":\\"Linear\\",\\"w\\":[[-0.3656227561185739,1.2246905186284998,0.9394040561359984],[-0.4799202859490072,-0.748051684382309,-0.4254887678392957],[-0.1324953494513111,0.39075553356182935,-0.9683373454711743],[0.7457834129742265,1.2588578403475645,-1.2293426686066573],[0.8486748227266485,-0.06246475936545304,-0.47424236877272036]]}],\\"mean\\":[0.9089397971666666,1020.3972950555559,1.8538348170740744,2121.146388,1787.8770603518515],\\"std\\":[0.035528861862199654,5.2898207605052745,0.012229934331297439,118.21647258590939,17.640434189614375]}\"
			},
			{
			  \"attributeId\": 83,
			  \"value\": \"1\"
			}
		  ]`

		bae := uuid.New().String()
		sql = addBAEModel(bae, c.TenantId, assetId, baeProps)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addModelParam(uuid.New().String(), bae, "parameter7", 0)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addModelParam(uuid.New().String(), bae, "parameter1", 1)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addModelParam(uuid.New().String(), bae, "parameter5", 2)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addModelParam(uuid.New().String(), bae, "parameter3", 3)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addModelParam(uuid.New().String(), bae, "parameter2", 4)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addModelHistory(uuid.New().String(), assetId, bae, "UNSUPERVISED")
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addStatus(assetId, bae)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		ruleId = uuid.New().String()
		sql = addModelAlarmRule(ruleId, assetId, bae)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		ruleItemId = uuid.New().String()
		sql = addModelAlarmRuleItem(ruleItemId, ruleId)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addModelAlarmRuleProperty(uuid.New().String(), ruleItemId)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}

		sql = addParamAlarmRule(assetId, c.ParamGroup)
		if _, err := f.WriteString(sql); err != nil {
			log.Println(err)
		}
	}

	sql = addParamAlarmRuleItem()
	if _, err := f.WriteString(sql); err != nil {
		log.Println(err)
	}

	sql = addParamAlarmRuleProperty()
	if _, err := f.WriteString(sql); err != nil {
		log.Println(err)
	}

	sql = addParamAssetValue(uuid.New().String(), c.TemplateId, c.ParamGroup)
	if _, err := f.WriteString(sql); err != nil {
		log.Println(err)
	}

	log.Println("completed.")
}

func addStatus(assetId, modelId string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO status(id, version, alarm_receiving, asset_id, description, master_id, model_id, name, role, value, ")
	b.WriteString("created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("(uuid(),0,1,'" + assetId + "','Unacceptable Description','status_001','" + modelId + "','Unacceptable','ASSET',1.0,")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")

	b.WriteString("INSERT INTO status(id, version, alarm_receiving, asset_id, description, master_id, model_id, name, role, value, ")
	b.WriteString("created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("(uuid(),0,1,'" + assetId + "','Unsatisfactory Description','status_002','" + modelId + "','Unsatisfactory','ASSET',0.8,")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")

	b.WriteString("INSERT INTO status(id, version, alarm_receiving, asset_id, description, master_id, model_id, name, role, value, ")
	b.WriteString("created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("(uuid(),0,0,'" + assetId + "','Satisfactory Description','status_003','" + modelId + "','Satisfactory','ASSET',0.5,")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")

	b.WriteString("INSERT INTO status(id, version, alarm_receiving, asset_id, description, master_id, model_id, name, role, ")
	b.WriteString("created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("(uuid(),0,0,'" + assetId + "','Good Description','status_004','" + modelId + "','Good','ASSET',")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")

	b.WriteString("commit;\n")

	return b.String()
}

func addModelAlarmRuleProperty(id, alarmRuleItemId string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO alarm_rule_property(id, version, name, value, alarm_rule_item_id, ")
	b.WriteString("created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + id + "',0,'status','STATUS1','" + alarmRuleItemId + "',")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")
	b.WriteString("commit;\n")
	return b.String()
}

func addModelAlarmRuleItem(id, alarmRuleId string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO alarm_rule_item(id, version, activating, alarm_rule_id, severity, ")
	b.WriteString("created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + id + "',1,1,'" + alarmRuleId + "','STATUS1',")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")
	b.WriteString("commit;\n")
	return b.String()
}

func addModelAlarmRule(id, assetId, modelId string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO alarm_rule (id, version, asset_id, master_id, model_id, role, ")
	b.WriteString("created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + id + "',0,'" + assetId + "','alarm_rule_master_001','" + modelId + "','ASSET',")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")
	b.WriteString("commit;\n")
	return b.String()
}

// v1.0
// func addParamAlarmRuleProperty() string {
// 	var b bytes.Buffer
// 	b.WriteString("INSERT INTO alarm_rule_property(id, version, name, value, alarm_rule_item_id, ")
// 	b.WriteString("created_by, created_time, modified_by, modified_time) ")
// 	b.WriteString("SELECT uuid() as id, 1, 'type', 'LIMIT', i.id, 'admin',NOW(),'admin',NOW() ")
// 	b.WriteString("FROM alarm_rule a, alarm_rule_item i ")
// 	b.WriteString("WHERE a.id=i.alarm_rule_id AND a.role='PARAMETER'; \n ")
// 	b.WriteString("commit;\n")
// 	return b.String()
// }

func addParamAlarmRuleProperty() string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO alarm_rule_property(id, version, name, value, alarm_rule_item_id, ")
	b.WriteString("created_by, created_time, modified_by, modified_time) ")
	b.WriteString("SELECT uuid() as id, 1, 'type', 'upperLimit', i.id, 'admin',NOW(),'admin',NOW() ")
	b.WriteString("FROM alarm_rule a, alarm_rule_item i ")
	b.WriteString("WHERE a.id=i.alarm_rule_id AND a.role='PARAMETER'; \n ")
	b.WriteString("commit;\n")
	return b.String()
}

func addParamAlarmRuleItem() string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO alarm_rule_item(id, version, activating, alarm_rule_id, severity, ")
	b.WriteString("created_by, created_time, modified_by, modified_time) ")
	b.WriteString("SELECT uuid() as id, 1, 1, id as alarm_rule_id, 'STATUS1', 'admin', NOW(), 'admin', NOW() ")
	b.WriteString("FROM alarm_rule WHERE role='PARAMETER'; \n")
	b.WriteString("commit;\n")
	return b.String()
}

func addParamAlarmRule(assetId, paramGroupId string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO alarm_rule (id, version, asset_id, master_id, param_id, role, ")
	b.WriteString("created_by, created_time, modified_by, modified_time) ")
	b.WriteString("SELECT uuid() as id,0,'" + assetId + "','alarm_rule_master_001',p.id as param_id,'PARAMETER',")
	b.WriteString("'admin', NOW(), 'admin', NOW() ")
	b.WriteString("FROM parameter p WHERE p.param_group_id='" + paramGroupId + "';\n")
	b.WriteString("commit;\n")
	return b.String()
}

func addModelHistory(id, assetId, modelId, modelName string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO model_history (id, asset_id, model_id, model_name, type) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + id + "','" + assetId + "','" + modelId + "','" + modelName + "','CREATE');\n")
	b.WriteString("commit;\n")
	return b.String()
}

func addModelParam(id, modelId, paramName string, sequence int) string {
	s := strconv.Itoa(sequence)

	var b bytes.Buffer
	b.WriteString("INSERT INTO model_param (id, version, model_id, param_id, sequence, created_by, created_time, modified_by, modified_time) ")
	b.WriteString("SELECT '" + id + "',0,'" + modelId + "', p.id as param_id," + s + ",'admin', NOW(), 'admin', NOW() ")
	b.WriteString("FROM parameter p WHERE p.name='" + paramName + "';\n")
	b.WriteString("commit;\n")
	return b.String()
}

//v1.0
// func addBAEModel(id, tenantId, assetId, modelProps string) string {
// 	var b bytes.Buffer
// 	b.WriteString("INSERT INTO model (id, version, tenant_id, activating, apply_start_time, asset_id, delegating, ")
// 	b.WriteString("model_props, analysis_model_version_id, created_by, created_time, modified_by, modified_time) ")
// 	b.WriteString("VALUES ")
// 	b.WriteString("('" + id + "',1,'" + tenantId + "',1,NOW(),'" + assetId + "',1,'" + modelProps + "','analysis_model_ver_unsuper_001',")
// 	b.WriteString("'admin',NOW(),'admin',NOW());\n")
// 	b.WriteString("commit;\n")
// 	return b.String()
// }

func addBAEModel(id, tenantId, assetId, modelProps string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO model (id, version, tenant_id, activating, apply_start_time, asset_id, delegating, ")
	b.WriteString("props, analysis_model_version_id, attribute_schema_id, created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + id + "',1,'" + tenantId + "',1,NOW(),'" + assetId + "',1,'" + modelProps + "','analysis_model_ver_unsuper_001','attribute_schema_model_unsuper',")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")
	b.WriteString("commit;\n")
	return b.String()
}

//v1.0
// func addFDCModel(id, tenantId, assetId string) string {
// 	var b bytes.Buffer
// 	b.WriteString("INSERT INTO model (id, version, tenant_id, activating, apply_start_time, asset_id, delegating, ")
// 	b.WriteString("model_props, analysis_model_version_id, created_by, created_time, modified_by, modified_time) ")
// 	b.WriteString("VALUES ")
// 	b.WriteString("('" + id + "',1,'" + tenantId + "',1,NOW(),'" + assetId + "',0,'[]','analysis_model_ver_fdc_001',")
// 	b.WriteString("'admin',NOW(),'admin',NOW());\n")
// 	b.WriteString("commit;\n")
// 	return b.String()
// }

func addFDCModel(id, tenantId, assetId string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO model (id, version, tenant_id, activating, apply_start_time, asset_id, delegating, ")
	b.WriteString("props, analysis_model_version_id, created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + id + "',1,'" + tenantId + "',1,NOW(),'" + assetId + "',0,'[]','analysis_model_ver_fdc_001',")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")
	b.WriteString("commit;\n")
	return b.String()
}

func addSensorStatus(sensorId string) string {
	id := uuid.New().String()

	var b bytes.Buffer
	b.WriteString("INSERT INTO sensor_status (id, collecting_rate, sensor_id, status) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + id + "',100,'" + sensorId + "',1);\n")
	b.WriteString("commit;\n")
	return b.String()
}

func linkSensorParamGroup(assetId, sensorId, paramGroupId string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO sensor_param_group_join (asset_id, param_group_id, sensor_id) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + assetId + "','" + paramGroupId + "','" + sensorId + "');\n")
	b.WriteString("commit;\n")
	return b.String()
}

func addSensor(id, assetId, sensorName, deviceId string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO sensor (id, version, asset_id, collecting, device_id, duration, name, physical_name, url, created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + id + "',0,'" + assetId + "',1,'" + deviceId + "','PT1S','" + sensorName + "', '" + sensorName + "','modbus://10.0.0.2:502',")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")
	b.WriteString("commit;\n")
	return b.String()
}

// v1.0
// func addDevice(id, tenantId, deviceName string) string {
// 	props := `{"mainPanel":"LOCAL CONTROL PANEL-21", "secondaryPanel":"WS-2 PANEL"}`

// 	var b bytes.Buffer
// 	b.WriteString("INSERT INTO device (id, version,	tenant_id, ip_addr, mac_addr, model_num, name, props, serial_num, created_by, created_time, modified_by, modified_time) ")
// 	b.WriteString("VALUES ")
// 	b.WriteString("('" + id + "',0,'" + tenantId + "','127.0.0.1','D3-E3-24-12-A3-F5','MODEL-X','" + deviceName + "','" + props + "','345-234234-122',")
// 	b.WriteString("'admin',NOW(),'admin',NOW());\n")
// 	b.WriteString("commit;\n")
// 	return b.String()
// }

func addDevice(id, tenantId, deviceName string) string {
	//props := `{"mainPanel":"LOCAL CONTROL PANEL-21", "secondaryPanel":"WS-2 PANEL"}`

	var b bytes.Buffer
	b.WriteString("INSERT INTO device (id, version,	tenant_id, name, physical_name, type_id, description, created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + id + "',0,'" + tenantId + "','" + deviceName + "','" + strings.ToUpper(deviceName) + "','type_device_02','Gateway Device',")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")
	b.WriteString("commit;\n")
	return b.String()
}

func addParamAssetValue(id, templateId, paramGroupId string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO parameter_value (id, asset_id, param_id, props) ")
	b.WriteString("SELECT uuid() as id, a.id as asset_id, v.param_id, v.props as props ")
	b.WriteString("FROM asset a, parameter p, parameter_value v ")
	b.WriteString("WHERE a.template_id='" + templateId + "' AND a.type_id='type_asset_02' ")
	b.WriteString("AND p.id=v.param_id AND v.asset_id is null ")
	b.WriteString("AND p.param_group_id='" + paramGroupId + "'; ")
	b.WriteString("commit;\n")
	return b.String()
}

func addParamModelValue(id, paramId, props string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO parameter_value (id, param_id, props) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + id + "','" + paramId + "','" + props + "');\n")
	b.WriteString("commit;\n")
	return b.String()
}

func addParameter(id, paramName, paramGroupId string, seq int) string {
	s := strconv.Itoa(seq)
	var b bytes.Buffer
	b.WriteString("INSERT INTO parameter (id, version, data_type, logical_type, name, physical_name, sequence, param_group_id, created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + id + "',0,'DOUBLE','DEFAULT','" + paramName + "','" + paramName + "'," + s + ",'" + paramGroupId + "',")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")
	b.WriteString("commit;\n")
	return b.String()
}

func addNestEgg(assetId, nestId string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO nest_egg (asset_id, nest_id) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + assetId + "','" + nestId + "');\n")
	b.WriteString("commit;\n")
	return b.String()
}

func linkAssetWithCatalog(catalogId, assetId string) string {
	var b bytes.Buffer
	b.WriteString("INSERT INTO asset_catalog_join (catalog_id, asset_id) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + catalogId + "','" + assetId + "');\n")
	b.WriteString("commit;\n")
	return b.String()
}

func queryForAsset(assetId, catalogId, assetName, templateId, typeId string) string {
	props := `[
			{
			  "dataType": "String",
			  "defaultValue": "",
			  "description": "",
			  "inputType": "TEXT",
			  "isChange": 0,
			  "isError": 0,
			  "isNew": 0,
			  "name": "Series Number",
			  "referenceType": "",
			  "seq": "1",
			  "value": "5478421570102"
			},
			{
			  "dataType": "String",
			  "defaultValue": "",
			  "description": "",
			  "inputType": "TEXT",
			  "isChange": 0,
			  "isError": 0,
			  "isNew": 0,
			  "name": "Manufacturer",
			  "referenceType": "",
			  "seq": "2",
			  "value": "Innolux"
			},
			{
			  "dataType": "String",
			  "defaultValue": "",
			  "description": "",
			  "inputType": "TEXT",
			  "isChange": 0,
			  "isError": 0,
			  "isNew": 0,
			  "name": "Frame",
			  "referenceType": "",
			  "seq": "3",
			  "value": "Frame01"
			}
		  ]`

	var b bytes.Buffer
	b.WriteString("INSERT INTO asset (id, version, catalog_id, image_url, name, physical_name, props, template_id, type_id, created_by, created_time) ")
	b.WriteString("VALUES ")
	b.WriteString("('" + assetId + "',0,'" + catalogId + "', 'apm://images/asset/poc_asset_01', '" + assetName + "','" + assetName + "','")
	b.WriteString(props + "','" + templateId + "','" + typeId + "','admin',NOW());\n")
	b.WriteString("commit;\n")
	return b.String()
}

func addParamType(tenantId string) string {
	var b bytes.Buffer

	b.WriteString("INSERT INTO type (id, version, tenant_id, name, role, sequence, created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("('type_parameter_02',0,'" + tenantId + "', 'Vibration', 'PARAM',2,")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")
	b.WriteString("commit;\n")
	return b.String()
}

func addAssetType(tenantId string) string {
	var b bytes.Buffer

	b.WriteString("INSERT INTO type (id, version, tenant_id, name, role, sequence, created_by, created_time, modified_by, modified_time) ")
	b.WriteString("VALUES ")
	b.WriteString("('type_asset_02',0,'" + tenantId + "', 'Pump', 'ASSET',2,")
	b.WriteString("'admin',NOW(),'admin',NOW());\n")
	b.WriteString("commit;\n")
	return b.String()
}

func updateParamType() string {
	var b bytes.Buffer

	b.WriteString("DELETE FROM parameter_value WHERE param_id in('poc_param_current_001','poc_param_current_002','poc_param_current_003');")
	b.WriteString("DELETE FROM parameter WHERE id in('poc_param_current_001','poc_param_current_002','poc_param_current_003');")

	b.WriteString("UPDATE param_group SET type_id='type_parameter_02', name='Vibration Param Group' WHERE id='poc_param_group_param_01';\n")
	b.WriteString("UPDATE type SET tenant_id='DEFAULT_ORG' WHERE role='DEVICE_STATUS'; \n")
	b.WriteString("UPDATE asset_template SET type_id='type_asset_02' WHERE id='poc_asset_model_01'; \n")

	stream_spec := `{
		"id": "default_nest_01",
		"topic": {
		  "source": "apm-trace-default-nest-01",
		  "score": "apm-score-default-nest-01",
		  "assetSpec": "apm-asset-spec-default-nest-01",
		  "alarmSpec": "apm-alarm-spec-default-nest-01",
		  "modelSpec": "apm-model-spec-default-nest-01"
		},
		"storage": {
		  "trace": "apm_trace_default_nest_01",
		  "score": "apm_score_default_nest_01"
		},
		"modelJobs": [],
		"parameters": []
	  }`
	b.WriteString("UPDATE nest set stream_spec='" + stream_spec + "' where id='default_nest_01'; \n")
	b.WriteString("commit;\n")
	return b.String()
}
