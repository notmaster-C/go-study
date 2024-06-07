package utils

import (
	"database/sql"
	"fmt"
	"go-study/db"
	"os"
	"strings"
)

type NetLog_OnlineDtail struct {
	UserID            string `gorm:"column:user_id" json:"userID"`
	UserMac           string `gorm:"column:user_mac" json:"userMac"`
	UserGroupID       string `gorm:"column:user_group_id" json:"userGroupID"`
	userTemplateId    string `gorm:"column:user_template_id" json:"userTemplateID"`
	suVersion         string `gorm:"column:su_version" json:"suVersion"`
	userIpv4          string `gorm:"column:user_ipv4" json:"userIPv4"`
	userIpv6          string `gorm:"column:user_ipv6" json:"userIPv6"`
	nasIp             string `gorm:"column:nas_ip" json:"nasIP"`
	nasIpv6           string `gorm:"column:nas_ipv6" json:"nasIPv6"`
	serviceId         string `gorm:"column:service_id" json:"serviceID"`
	policyId          string `gorm:"column:policy_id" json:"policyID"`
	accountId         string `gorm:"column:account_id" json:"accountID"`
	wpNasIp           string `gorm:"column:wp_nas_ip" json:"wpNasIP"`
	wpNasPort         string `gorm:"column:wp_nas_port" json:"wpNasPort"`
	proxyName         string `gorm:"column:proxy_name" json:"proxyName"`
	loginTime         string `gorm:"column:login_time" json:"loginTime"`
	logoutTime        string `gorm:"column:logout_time" json:"logoutTime"`
	onlineSec         string `gorm:"column:online_sec" json:"onlineSec"`
	terminateCause    string `gorm:"column:terminate_cause" json:"terminateCause"`
	tunnelClient      string `gorm:"column:tunnel_client" json:"tunnelClient"`
	tunnelServer      string `gorm:"column:tunnel_server" json:"tunnelServer"`
	apMac             string `gorm:"column:ap_mac" json:"apMac"`
	areaName          string `gorm:"column:area_name" json:"areaName"`
	ssid              string `gorm:"column:ssid" json:"ssid"`
	isRoaming         string `gorm:"column:is_roaming" json:"isRoaming"`
	totalTraffic      string `gorm:"column:total_traffic" json:"totalTraffic"`
	serviceSuffix     string `gorm:"column:service_suffix" json:"serviceSuffix"`
	accountingRule    string `gorm:"column:accounting_rule" json:"accountingRule"`
	packageName       string `gorm:"column:package_name" json:"packageName"`
	timesegmentId     string `gorm:"column:timesegment_id" json:"timesegmentID"`
	internetAccessFee string `gorm:"column:internet_access_fee" json:"internetAccessFee"`
	reserved0         string `gorm:"column:reserved_0" json:"reserved0"`
	reserved1         string `gorm:"column:reserved_1" json:"reserved1"`
	reserved2         string `gorm:"column:reserved_2" json:"reserved2"`
	reserved3         string `gorm:"column:reserved_3" json:"reserved3"`
	reserved4         string `gorm:"column:reserved_4" json:"reserved4"`
	reserved5         string `gorm:"column:reserved_5" json:"reserved5"`
	reserved6         string `gorm:"column:reserved_6" json:"reserved6"`
	reserved7         string `gorm:"column:reserved_7" json:"reserved7"`
	reserved8         string `gorm:"column:reserved_8" json:"reserved8"`
	reserved9         string `gorm:"column:reserved_9" json:"reserved9"`
	onlineDetailUuid  string `gorm:"column:online_detail_uuid" json:"onlineDetailUUID"`
}

// ParseFile 解析数据
func ParseFile(filename string) ([]NetLog_OnlineDtail, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data []NetLog_OnlineDtail
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		//line := lines[1]
		if strings.TrimSpace(line) == "" {
			continue
		}
		fmt.Println(line)
		userData := NetLog_OnlineDtail{}
		err = parseLine(line, &userData)
		if err != nil {
			return nil, err
		}

		data = append(data, userData)
	}

	return data, nil
}
func parseLine(line string, userData *NetLog_OnlineDtail) error {
	fields := strings.Split(line, "\t")
	userData.UserID = fields[0]
	userData.UserMac = fields[1]
	userData.UserGroupID = fields[2]
	userData.userTemplateId = fields[3]
	userData.suVersion = fields[4]
	userData.userIpv4 = fields[5]
	userData.userIpv6 = fields[6]
	userData.nasIp = fields[7]
	userData.nasIpv6 = fields[8]
	userData.serviceId = fields[9]
	userData.policyId = fields[10]
	userData.accountId = fields[11]
	userData.wpNasIp = fields[12]
	userData.wpNasPort = fields[13]
	userData.proxyName = fields[14]
	userData.loginTime = fields[15]
	userData.logoutTime = fields[16]
	userData.onlineSec = fields[17]
	userData.terminateCause = fields[18]
	userData.tunnelClient = fields[19]
	userData.tunnelServer = fields[20]
	userData.apMac = fields[21]
	userData.areaName = fields[22]
	userData.ssid = fields[23]
	userData.isRoaming = fields[24]
	userData.totalTraffic = fields[25]
	userData.serviceSuffix = fields[26]
	userData.accountingRule = fields[27]
	userData.packageName = fields[28]
	userData.timesegmentId = fields[29]
	userData.internetAccessFee = fields[30]
	userData.reserved0 = fields[31]
	userData.reserved1 = fields[32]
	userData.reserved2 = fields[33]
	userData.reserved3 = fields[34]
	userData.reserved4 = fields[35]
	userData.reserved5 = fields[36]
	userData.reserved6 = fields[37]
	userData.reserved7 = fields[38]
	userData.reserved8 = fields[39]
	userData.reserved9 = fields[40]
	userData.onlineDetailUuid = fields[41]
	fmt.Println("userData:", userData)
	return nil
}

// InsertData 存入数据库
func InsertData(data []NetLog_OnlineDtail) error {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users (user_id, user_mac, user_group_id) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, userData := range data {
		_, err := stmt.Exec(userData.UserID, userData.UserMac, userData.UserGroupID)
		if err != nil {
			return err
		}
	}

	return nil
}

// BulkInsertData 批量插入方法
func BulkInsertData(data []NetLog_OnlineDtail) error {
	sdb := db.GetDb("mysql")
	var err error
	if err != nil {
		return err
	}

	tx := sdb.Begin()
	for _, userData := range data {
		err := tx.Create(&userData).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
func Excute() {
	data, err := ParseFile("E:\\csl\\20230825_online_detail")
	if err != nil {
		fmt.Errorf("ParseFile Error！%s", err)
	}
	batchSize := 1000 // 每批次插入的数据量
	err = db.GetDb("mysql").AutoMigrate(&NetLog_OnlineDtail{})
	if err != nil {
		fmt.Println("建表错误")
		return
	}
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		batch := data[i:end]
		err := BulkInsertData(batch)
		if err != nil {
			// 处理错误
		}
	}
}
