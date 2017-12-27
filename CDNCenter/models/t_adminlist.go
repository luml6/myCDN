package models

import (
	"github.com/astaxie/beego/orm"
)

type TAdminlist struct {
	Username string `orm:"column(userName);pk"`
	Pwd      string `orm:"column(pwd);size(50);null"`
}

func (t *TAdminlist) TableName() string {
	return "t_adminlist"
}

func init() {
	orm.RegisterModel(new(TAdminlist))
}
func ValidateUser(userName string) (v *TAdminlist, err error) {
	o := orm.NewOrm()
	v = &TAdminlist{Username: userName}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}
