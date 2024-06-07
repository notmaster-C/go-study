package db

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	gorm.Model
}

// 宿舍晚未归配置
type SSWWG struct {
	Id            int64  `gorm:"column:Id;primarykey;autoIncrement"  column:"-"`
	LatebeginDate string `gorm:"column:LatebeginDate;type:varchar(64);default:'230000';" column:"LatebeginDate"` //晚归判断机制开始时间
	LateendDate   string `gorm:"column:LateendDate;type:varchar(64);default:'060000'" column:"LateendDate"`      //晚归判断机制结束时间
	LateRate      string `gorm:"column:LateRate;type:int(16);default:1;" column:"LateRate"`                      //晚归通知频率 (1 天,7 周,30 月)
	OutbeginDate  string `gorm:"column:OutbeginDate;type:varchar(64);default:'230000';" column:"OutbeginDate"`   //夜不归宿判断机制开始时间
	OutendDate    string `gorm:"column:OutendDate;type:varchar(64);default:'060000';" column:"OutendDate"`       //夜不归宿判断机制结束时间
	OutRate       string `gorm:"column:OutRate;type:int(16);default:1;" column:"OutRate"`                        //夜不归宿通知频率
	StayHour      string `gorm:"column:StayHour;type:int(16);default:24;" column:"StayHour"`                     //不出不归判断机制(单位：小时)
	StayRate      string `gorm:"column:StayRate;type:int(16);default:1;" column:"StayRate"`                      //不出不归通知频率
}
type SSWWG11 struct {
	Xh         string `form:"xh" json:"xh" xml:"xh"`                         // 学院代码 就是学院的code 数组
	XueyuanId  string `form:"xueyuanIds" json:"xueyuanIds" xml:"xueyuanIds"` // 学院代码 就是学院的code 数组
	NianjiName string `form:"nianjiName" json:"nianjiName" xml:"nianjiName"` // 年级名称 比如2022
	StuType    int64  `form:"stuType" json:"stuType" xml:"stuType"`          // 学生类型 0本科生 1研究生
	Key        string `form:"key" json:"key" xml:"key"`                      // 学生类型 college本科生或者graduate研究生
	StartTime  string `form:"startTime" json:"startTime" xml:"startTime"`
	EndTime    string `form:"endTime" json:"endTime" xml:"endTime"`
}
type SSWWG22 struct {
	XueyuanIds []string `form:"xueyuanIds" json:"xueyuanIds" xml:"xueyuanIds"` // 学院代码 就是学院的code 数组
	NianjiName []string `form:"nianjiName" json:"nianjiName" xml:"nianjiName"` // 年级名称 比如2022
	StartTime  string   `form:"startTime" json:"startTime" xml:"startTime"`
	EndTime    string   `form:"endTime" json:"endTime" xml:"endTime"`
	Key        string   `form:"key" json:"key" xml:"key"` // 学生类型 college本科生或者graduate研究生
}
type SSWWG33 struct {
	XueyuanId  string `form:"xueyuanIds" json:"xueyuanIds" xml:"xueyuanIds"` // 学院代码 就是学院的code 数组
	NianjiName string `form:"nianjiName" json:"nianjiName" xml:"nianjiName"` // 年级名称 比如2022
	StartTime  string `form:"startTime" json:"startTime" xml:"startTime"`
	EndTime    string `form:"endTime" json:"endTime" xml:"endTime"`
	Key        string `form:"key" json:"key" xml:"key"` // 学生类型 college本科生或者graduate研究生
}

// 统计数据来源表
type StatDataSource struct {
	Id         int64     `gorm:"primarykey"  json:"-"`
	StatId     int64     `gorm:"type:int(16);not null;index;" json:"stat_id"`  // 统计表自增ID
	Type       string    `gorm:"type:char(6);not null" json:"-"`               // 统计类型 db,api
	Key        string    `gorm:"type:varchar(64);not null;index" json:"key"`   // 统计key
	Param      string    `gorm:"type:text;not null;" json:"param"`             // 统计具体内容 「sql语句/api接口」
	ExtendKey  string    `gorm:"type:varchar(24);not null;" json:"extend_key"` // 数据来源扩展信息
	CreatedAt  time.Time `gorm:"type:datetime" json:"-"`                       // 创建时间
	UpdatedAt  time.Time `gorm:"type:datetime" json:"-"`                       // 更新时间
	ResultType int       `gorm:"type:int(8);" json:"-"`                        //返回结果 0就是之前的单个 1就是json字符串
}
type ApiDbMiddleSource struct {
	Id        int64     `gorm:"primarykey"  json:"-"`
	Key       string    `gorm:"type:varchar(64);not null;index" json:"key"` // 统计key
	Field1    string    `gorm:"type:text;not null;" json:"field1"`          // 内容
	Field2    string    `gorm:"type:text;not null;" json:"field2"`          // 内容
	Field3    string    `gorm:"type:text;not null;" json:"field3"`          // 内容
	Field4    string    `gorm:"type:text;not null;" json:"field4"`          // 内容
	Field5    string    `gorm:"type:text;not null;" json:"field5"`          // 内容
	Field6    string    `gorm:"type:text;not null;" json:"field6"`          // 内容
	Field7    string    `gorm:"type:text;not null;" json:"field7"`          // 内容
	Field8    string    `gorm:"type:text;not null;" json:"field8"`          // 内容
	Field9    string    `gorm:"type:text;not null;" json:"field9"`          // 内容
	Field10   string    `gorm:"type:text;not null;" json:"field10"`         // 内容
	Field11   string    `gorm:"type:text;not null;" json:"field11"`         // 内容
	Field12   string    `gorm:"type:text;not null;" json:"field12"`         // 内容
	Field13   string    `gorm:"type:text;not null;" json:"field13"`         // 内容
	Field14   string    `gorm:"type:text;not null;" json:"field14"`         // 内容
	Field15   string    `gorm:"type:text;not null;" json:"field15"`         // 内容
	Field16   string    `gorm:"type:text;not null;" json:"field16"`         // 内容
	Field17   string    `gorm:"type:text;not null;" json:"field17"`         // 内容
	Field18   string    `gorm:"type:text;not null;" json:"field18"`         // 内容
	Field19   string    `gorm:"type:text;not null;" json:"field19"`         // 内容
	Field20   string    `gorm:"type:text;not null;" json:"field20"`         // 内容
	Field21   string    `gorm:"type:text;not null;" json:"field21"`         // 内容
	Field22   string    `gorm:"type:text;not null;" json:"field22"`         // 内容
	Field23   string    `gorm:"type:text;not null;" json:"field23"`         // 内容
	Field24   string    `gorm:"type:text;not null;" json:"field24"`         // 内容
	CreatedAt time.Time `gorm:"type:datetime" json:"-"`                     // 创建时间

}

var (
	dbClient = make(map[string]*gorm.DB)
	//pwd = "Baiduyun@123"
	//dsn = fmt.Sprintf("baiduyun:%s@tcp(192.168.10.47:3306)/bigdata_core?charset=utf8&parseTime=true", pwd)
)

// 应理菜单
type AppMenu struct {
	Id          uint64     `gorm:"primarykey"  json:"-"`
	Key         string     `gorm:"type:varchar(64);not null;index" json:"key"`   // 菜单key
	Icon        string     `gorm:"type:varchar(64);not null;" json:"icon"`       // icon图标
	Name        string     `gorm:"type:varchar(2048);not null;" json:"name"`     // 菜单名字，支持多语言格式 {"en_US":"name","type":"i18n","zh_CN":"名字"}
	Path        string     `gorm:"type:varchar(512);not null;" json:"path"`      // 路径
	IsShow      bool       `gorm:"type:tinyint(1);not null;default 0;" json:"-"` // 是否显示 0-不显示 1-显示
	ParentId    uint64     `gorm:"parent_id"  json:"-"`                          // 父亲节点Id
	OrderId     int16      `gorm:"type:int(16);index;not null" json:"-"`         // 排序id
	DisplayArea string     `gorm:"type:varchar(32);not null" json:"displayArea"` // 显示区域 web h5
	Privilege   string     `gorm:"type:text;not null;" json:"privilege"`         // 权限信息
	Children    []*AppMenu `gorm:"-" json:"children,omitempty"`                  // 子菜单信息
	CreatedAt   time.Time  `gorm:"type:datetime" json:"-"`                       // 创建时间
	UpdatedAt   time.Time  `gorm:"type:datetime" json:"-"`                       // 更新时间
}

// 统计数据来源扩展表
type StatDataSourceExtend struct {
	Id        int64     `gorm:"primarykey"  json:"-"`
	Type      string    `gorm:"type:char(6);not null;index" json:"-"`        // 扩展类型 db,api
	Key       string    `gorm:"type:varchar(64);not null;unique" json:"key"` // 扩展key
	Info      string    `gorm:"type:varchar(1024);not null;" json:"info"`    // 来源信息 数据库连接json串 API认证json串
	CreatedAt time.Time `gorm:"type:datetime" json:"-"`                      // 创建时间
	UpdatedAt time.Time `gorm:"type:datetime" json:"-"`                      // 更新时间
}
type MiddleData struct {
	Id       uint64 `gorm:"primarykey" json:"id"`                                                   // 学号
	Result   string `gorm:"column:result;type:varchar(100);comment:数据字典名称比如:学院;" json:"result"`     // 姓名
	KeyInfo  string `gorm:"column:key_info;type:varchar(100);comment:数据字典名称比如:学院;" json:"key_info"` // 姓名
	Ex1      string `gorm:"column:ex1;type:varchar(100);comment:特殊用处1;" json:"ex1"`                 // 性别码 // 出生日期
	Ex1Name  string `gorm:"column:ex1_name;type:varchar(100);comment:ex1描述;" json:"ex1_name"`       // 性别码 // 出生日期
	Ex2      string `gorm:"column:ex2;type:varchar(100);comment:特殊用处1;" json:"ex12"`                // 性别码 // 出生日期
	Ex2Name  string `gorm:"column:ex2_name;type:varchar(100);comment:ex1描述;" json:"ex2_name"`       // 性别码 // 出生日期
	Ex3      string `gorm:"column:ex3;type:varchar(100);comment:特殊用处1;" json:"ex3"`                 // 性别码 // 出生日期
	Ex3Name  string `gorm:"column:ex3_name;type:varchar(100);comment:ex1描述;" json:"ex3_name"`       // 性别码 // 出生日期
	Ex4      string `gorm:"column:ex4;type:varchar(100);comment:特殊用处1;" json:"ex4"`                 // 性别码 // 出生日期
	Ex4Name  string `gorm:"column:ex4_name;type:varchar(100);comment:ex1描述;" json:"ex4_name"`       // 性别码 // 出生日期
	Ex5      string `gorm:"column:ex5;type:varchar(100);comment:特殊用处1;" json:"ex5"`                 // 性别码 // 出生日期
	Ex5Name  string `gorm:"column:ex5_name;type:varchar(100);comment:ex1描述;" json:"ex5_name"`       // 性别码 // 出生日期
	Ex6      string `gorm:"column:ex6;type:varchar(100);comment:特殊用处1;" json:"ex6"`                 // 性别码 // 出生日期
	Ex6Name  string `gorm:"column:ex6_name;type:varchar(100);comment:ex1描述;" json:"ex6_name"`       // 性别码 // 出生日期
	Ex7      string `gorm:"column:ex7;type:varchar(100);comment:特殊用处1;" json:"ex7"`                 // 性别码 // 出生日期
	Ex7Name  string `gorm:"column:ex7_name;type:varchar(100);comment:ex1描述;" json:"ex7_name"`       // 性别码 // 出生日期
	Ex8      string `gorm:"column:ex8;type:varchar(100);comment:特殊用处1;" json:"ex8"`                 // 性别码 // 出生日期
	Ex8Name  string `gorm:"column:ex8_name;type:varchar(100);comment:ex1描述;" json:"ex8_name"`       // 性别码 // 出生日期
	Ex9      string `gorm:"column:ex9;type:varchar(100);comment:特殊用处1;" json:"ex9"`                 // 性别码 // 出生日期
	Ex9Name  string `gorm:"column:ex9_name;type:varchar(100);comment:ex1描述;" json:"ex9_name"`       // 性别码 // 出生日期
	Ex10     string `gorm:"column:ex10;type:varchar(100);comment:特殊用处1;" json:"ex10"`               // 性别码 // 出生日期
	Ex10Name string `gorm:"column:ex10_name;type:varchar(100);comment:ex1描述;" json:"ex10_name"`     // 性别码 // 出生日期

	// 身份证号
}

// 表单填报数据
type FormFillingData struct {
	Id              uint64             `gorm:"primarykey"  json:"id"`
	AppKey          string             `gorm:"type:varchar(64);not null;index" json:"appKey"`        // 应用唯一key
	FormId          string             `gorm:"type:varchar(64);not null;index" json:"formId"`        // 表单唯一id
	Content         string             `gorm:"type:longtext;" json:"content"`                        // 表单填报内容
	Creator         string             `gorm:"type:varchar(32);index" json:"creator"`                // 创建者id
	CreatorName     string             `gorm:"type:varchar(256);not null;index" json:"creatorName"`  // 创建者名称
	Editor          string             `gorm:"type:varchar(32);index" json:"editor"`                 // 编辑者id
	EditorName      string             `gorm:"-" json:"editorName,omitempty"`                        // 编辑者名称
	Approval        string             `gorm:"type:varchar(32);index" json:"approval,omitempty"`     // 审批者id
	ApprovalName    string             `gorm:"-" json:"approvalName,omitempty"`                      // 审批者名称
	Status          string             `gorm:"type:varchar(32);index" json:"status,omitempty"`       // 状态 【提交:submit 暂存:draft】
	ProcessFillStat int                `gorm:"type:int(16);default 0;index" json:"-"`                // 流程填报是否统计 【0:否 1:是】
	IsDelete        bool               `gorm:"type:tinyint(3);default 0;index" json:"isDelete"`      // 是否删除【0:否 1:是】
	Semester        string             `gorm:"type:varchar(255);index:idx_semester" json:"semester"` // 开课学期
	ProcessNodeName string             `gorm:"-" json:"processNodeName,omitempty"`                   // 流程节点名称
	Attachment      []string           `gorm:"-" json:"attachment,omitempty"`                        // 流程附件
	CreatorMobile   string             `gorm:"-" json:"creatorMobile,omitempty"`                     // 创建者手机号
	FillData        []*FormFillingData `gorm:"-" json:"fillData,omitempty"`                          // 子表单数据
	CreatedAt       time.Time          `gorm:"type:datetime;index" json:"createdAt,omitempty"`       // 创建时间
	UpdatedAt       time.Time          `gorm:"type:datetime" json:"updatedAt,omitempty"`             // 更新时间
}

// 统计报表
type StatReport struct {
	Id        uint64      `gorm:"primarykey"  json:"-"`
	Key       string      `gorm:"type:varchar(64);not null;index;" json:"key"`   // 统计key
	Name      string      `gorm:"type:varchar(128);not null;index;" json:"name"` // 统计字段名称
	Memo      string      `gorm:"type:varchar(256);not null;" json:"memo"`       // 统计字段中文描述
	Type      int8        `gorm:"type:int(16);not null" json:"-"`                // 统计字段类型，0字符串，1整型，2double类型
	Value     string      `gorm:"type:longtext;" json:"-"`                       // 统计字段内容
	Result    interface{} `gorm:"-" json:"value,omitempty"`                      // 统计字段内容
	CreatedAt time.Time   `gorm:"type:datetime" json:"-"`                        // 创建时间
	UpdatedAt time.Time   `gorm:"type:datetime" json:"-"`                        // 更新时间
}

type XJT1 struct {
	Code   int    `json:"code"`
	Uuid   string `json:"uuid"`
	Result struct {
		Total   int     `json:"total"`
		Size    int     `json:"size"`
		Current int     `json:"current"`
		Pages   int     `json:"pages"`
		Records Records `json:"records"`
	} `json:"result"`
}
type Records []struct {
	Id                  string `gorm:"column:Id;type:varchar(100);primarykey" json:"id"`
	Name                string `gorm:"column:ex1;type:varchar(64);" json:"name"`
	Description         string `gorm:"column:description;type:varchar(64);" json:"description"`
	DbType              string `gorm:"column:dbType;type:varchar(64);" json:"dbType"`
	DbCategoryAlias     string `gorm:"column:dbCategoryAlias;type:varchar(64);" json:"dbCategoryAlias"`
	IsEnable            bool   `gorm:"column:isEnable;type:varchar(64);" json:"isEnable"`
	IsDeleted           bool   `gorm:"column:isDeleted;type:varchar(64);" json:"isDeleted"`
	State               string `gorm:"column:state;type:varchar(64);" json:"state"`
	CreateTime          string `gorm:"column:createTime;type:varchar(64);" json:"createTime"`
	CreateBy            string `gorm:"column:createBy;type:varchar(64);" json:"createBy"`
	UpdateTime          string `gorm:"column:updateTime;type:varchar(64);" json:"updateTime"`
	DeptId              string `gorm:"column:deptId;type:varchar(64);" json:"deptId"`
	DeptName            string `gorm:"column:deptName;type:varchar(64);" json:"deptName"`
	UserType            string `gorm:"column:userType;type:varchar(64);" json:"userType"`
	BusinessSystem      string `gorm:"column:business_system;type:varchar(64);" json:"business_system"`
	SystemSupplier      string `gorm:"column:system_supplier;type:varchar(64);" json:"system_supplier"`
	BusinessTypeDisplay string `gorm:"column:business_type_display;type:varchar(64);" json:"business_type_display"`
	BusinessType        string `gorm:"column:business_type;type:varchar(64);" json:"business_type"`
	BusinessSystemId    string `gorm:"column:businessSystemId;type:varchar(100);" json:"businessSystemId"`
	BusinessSupplierId  string `gorm:"column:businessSupplierId;type:varchar(100);" json:"businessSupplierId"`
}

// SuperviseCourseListenTemp 督导课程听课信息临时处理表
type SuperviseCourseListenTemp struct {
	Id             uint64 `gorm:"type:bigint(32);autoIncrement;primarykey" json:"id"`
	ClassNumber    string `gorm:"type:varchar(255);index" json:"classNumber"`                // 开课号
	TeacherNumber  string `gorm:"type:varchar(64);index" json:"teacherNumber"`               // 任课老师工号
	ListenTime     string `gorm:"type:varchar(255);not null" json:"listenTime"`              // 听课时间
	AddressCode    int64  `gorm:"type:varchar(255);not null" json:"AddressCode"`             // 听课地点
	Weeks          string `gorm:"type:varchar(255);not null" json:"weeks"`                   // 周次
	Week           int    `gorm:"type:varchar(255);not null" json:"week"`                    // 星期几
	Sessions       string `gorm:"type:varchar(100);not null;" json:"sessions"`               // 节次
	WeeksBinary    int64  `gorm:"type:int(11);not null;index" json:"weeksBinary"`            // 周次段
	SessionsBinary int64  `gorm:"type:int(11);not null;index" json:"sessionsBinary"`         // 节次段
	Status         bool   `gorm:"type:tinyint(1);not null;default:0;" json:"status"`         // 更新状态
	StartTime      int64  `gorm:"type:int(11);index:idx_address_start_end" json:"startTime"` // 听课开始时间
	EndTime        int64  `gorm:"type:int(11);index:idx_address_start_end" json:"endTime"`   // 听课结束时间
	Campus         string `gorm:"type:varchar(100);not null;" json:"campus"`                 // 校区
	TotalStudent   int64  `gorm:"type:int(64);" json:"totalStudent"`                         // 总学生
}

type AssitantBooks struct {
	Id   uint64 `gorm:"type:bigint(32);autoIncrement;primarykey" json:"-"`
	Jsgh string `gorm:"type:varchar(255);index" json:"Jsgh"`
	Jcrq string `gorm:"type:varchar(64);index:idx_jsgh_jcrq" json:"Jcrq"`
	Yhrq string `gorm:"type:varchar(64);index:idx_jsgh_yhrq" json:"Yhrq"`
	Tssm string `gorm:"type:varchar(64);" json:"-"`
}

// 督导课程听课信息
type SuperviseCourseListen struct {
	ClassNumber   string `gorm:"type:varchar(255);index" json:"classNumber"`                  // 开课号
	TeacherNumber string `gorm:"type:varchar(64);index" json:"teacherNumber"`                 // 任课老师工号
	AddressCode   int64  `gorm:"type:int(11);index:idx_address_start_end" json:"addressCode"` // 听课地点编码
	Campus        string `gorm:"type:varchar(255)" json:"campus"`                             // 所属校区
	StartTime     int64  `gorm:"type:int(11);index:idx_address_start_end" json:"startTime"`   // 听课开始时间
	EndTime       int64  `gorm:"type:int(11);index:idx_address_start_end" json:"endTime"`     // 听课结束时间
	TotalStudent  int64  `gorm:"type:int(64);" json:"totalStudent"`                           // 总学生
}
