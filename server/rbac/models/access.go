package models

type Access struct {
	Id          int64
	ModuleName  string //模块名称
	ActionName  string //操作名称
	Type        int64  //节点类型 :  1、表示模块    2、表示菜单     3、操作
	Url         string //路由跳转地址
	ModuleId    int64  //此module_id和当前模型的id关联       module_id= 0 表示模块
	Sort        int64
	Description string
	Status      int64
	AddTime     int64
	AccessItem  []Access `gorm:"foreignKey:ModuleId;references:Id"`
	Checked     bool     `gorm:"-"` // 忽略本字段
}

func (Access) TableName() string {
	return "access"
}
