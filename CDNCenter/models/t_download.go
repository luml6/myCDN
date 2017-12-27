package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"reflect"
	"strings"
	"time"
)

type TDownload struct {
	Id           int       `orm:"column(id);auto"`
	Ip           string    `orm:"column(ip);size(50);null"`
	Download     string    `orm:"column(download);size(50);null"`
	DownSize     int       `orm:"column(downSize);null"`
	DownloadTime time.Time `orm:"column(downloadTime);null"`
}

func (t *TDownload) TableName() string {
	return "t_downloadlog"
}

func init() {
	orm.RegisterModel(new(TDownload))
}
func AddAllDownload(m []TDownload) (err error) {
	o := orm.NewOrm()
	qs := o.QueryTable("t_downloadlog")
	i, _ := qs.PrepareInsert()
	for _, logger := range m {
		_, err := i.Insert(&logger)
		if err != nil {
			beego.Debug(err)
			break
		}
	}
	i.Close() // 别忘记关闭 statement
	return
}
func FindTDownloadAll() (count int, channel []TDownload) {
	o := orm.NewOrm()
	var ch []TDownload
	sSql := "SELECT ip,download,downSize,downloadTime FROM t_downloadlog ORDER BY ip "
	num, err := o.Raw(sSql).QueryRows(&ch)
	if err == nil && num > 0 {
		return int(num), ch
	}
	beego.Debug(err)
	return 0, nil
}
func FindTDownloadData(dateStart, dateEnd string) (count int, channel []TDownload) {
	o := orm.NewOrm()
	var ch []TDownload
	sSql := "SELECT ip,download,downSize,downloadTime FROM t_downloadlog WHERE downloadTime BETWEEN '" + dateStart + "' AND '" + dateEnd + "' ORDER BY ip "
	num, err := o.Raw(sSql).QueryRows(&ch)
	if err == nil && num > 0 {
		return int(num), ch
	}
	beego.Debug(err)
	return 0, nil
}
func FindTDownloadWithDatePage(pageNum, pageCount, totalCount int, dateStart, dateEnd string) (num int, channel []TDownload) {
	var ch []TDownload
	if pageCount > totalCount {
		o := orm.NewOrm()
		sSql := "SELECT ip,download,downSize,downloadTime FROM t_downloadlog  WHERE downloadTime BETWEEN '" + dateStart + "' AND '" + dateEnd + "'  ORDER BY ip DESC"
		num, err := o.Raw(sSql).QueryRows(&ch)
		if err == nil && num > 0 {
			return pageNum, ch
		}
		beego.Debug(err)
	} else {
		if pageNum*pageCount > totalCount {
			listCount := pageCount - (pageNum*pageCount - totalCount)
			o := orm.NewOrm()
			sSql := "SELECT ip,download,downSize,downloadTime FROM t_downloadlog  WHERE downloadTime BETWEEN '" + dateStart + "' AND '" + dateEnd + "' AND  ORDER BY ip LIMIT ? OFFSET 0"
			num, err := o.Raw(sSql, listCount).QueryRows(&ch)
			if err == nil && num > 0 {
				return pageNum, ch
			}
			beego.Debug(err)
		} else {
			o := orm.NewOrm()
			sSql := "SELECT ip,download,downSize,downloadTime FROM t_downloadlog  WHERE downloadTime BETWEEN '" + dateStart + "' AND '" + dateEnd + "'  ORDER BY ip DESC LIMIT ? OFFSET ?"
			num, err := o.Raw(sSql, pageCount, (pageNum-1)*pageCount).QueryRows(&ch)
			if err == nil && num > 0 {
				return pageNum, ch
			}
			beego.Debug(err)
		}
	}
	return 0, nil
}
func FindTDownloadIP(ip string) (count int, channel []TDownload) {
	o := orm.NewOrm()
	var ch []TDownload
	sSql := "SELECT ip,download,downSize,downloadTime FROM t_downloadlog WHERE ip like '%" + ip + "%' ORDER BY ip "
	num, err := o.Raw(sSql).QueryRows(&ch)
	if err == nil && num > 0 {
		return int(num), ch
	}
	beego.Debug(err)
	return 0, nil
}
func FindTDownloadWithIPPage(pageNum, pageCount, totalCount int, ip string) (num int, channel []TDownload) {
	var ch []TDownload
	if pageCount > totalCount {
		o := orm.NewOrm()
		sSql := "SELECT ip,download,downSize,downloadTime FROM t_downloadlog  WHERE ip like '%" + ip + "%'  ORDER BY ip DESC"
		num, err := o.Raw(sSql).QueryRows(&ch)
		if err == nil && num > 0 {
			return pageNum, ch
		}
		beego.Debug(err)
	} else {
		if pageNum*pageCount > totalCount {
			listCount := pageCount - (pageNum*pageCount - totalCount)
			o := orm.NewOrm()
			sSql := "SELECT ip,download,downSize,downloadTime FROM t_downloadlog  WHERE ip like '%" + ip + "%' AND  ORDER BY ip LIMIT ? OFFSET 0"
			num, err := o.Raw(sSql, listCount).QueryRows(&ch)
			if err == nil && num > 0 {
				return pageNum, ch
			}
			beego.Debug(err)
		} else {
			o := orm.NewOrm()
			sSql := "SELECT ip,download,downSize,downloadTime FROM t_downloadlog  WHERE ip like '%" + ip + "%'  ORDER BY ip DESC LIMIT ? OFFSET ?"
			num, err := o.Raw(sSql, pageCount, (pageNum-1)*pageCount).QueryRows(&ch)
			if err == nil && num > 0 {
				return pageNum, ch
			}
			beego.Debug(err)
		}
	}
	return 0, nil
}
func FindTDownloadWithPage(pageNum, pageCount, totalCount int) (num int, channel []TDownload) {
	var ch []TDownload
	o := orm.NewOrm()
	if pageCount > totalCount {
		sql := "SELECT ip,download,downSize,downloadTime FROM t_downloadlog ORDER BY ip "
		// 执行SQL语句
		num, err := o.Raw(sql).QueryRows(&ch)

		if err == nil && num > 0 {
			return pageNum, ch
		}
		beego.Debug(err)
	} else {
		if pageNum*pageCount > totalCount {
			listCount := pageCount - (pageNum*pageCount - totalCount)
			// qb, _ := orm.NewQueryBuilder("mysql")
			// qb.Select("*").From("t_channel").Limit(listCount).Offset(0)
			sql := "SELECT ip,download,downSize,downloadTime FROM t_downloadlog ORDER BY ip LIMIT ? OFFSET 0"
			// 执行SQL语句
			num, err := o.Raw(sql, listCount).QueryRows(&ch)

			if err == nil && num > 0 {
				return pageNum, ch
			}
			beego.Debug(err)
		} else {
			// qb, _ := orm.NewQueryBuilder("mysql")
			// qb.Select("*").From("t_channel").OrderBy("Id").Desc().Limit(pageCount).Offset((pageNum - 1) * pageCount)
			sql := "SELECT ip,download,downSize,downloadTime FROM t_downloadlog ORDER BY ip LIMIT ? OFFSET ? "
			num, err := o.Raw(sql, pageCount, (pageNum-1)*pageCount).QueryRows(&ch)

			if err == nil && num > 0 {
				return pageNum, ch
			}
			beego.Debug(err)
		}
	}
	return 0, nil
}
func FindDownloadCount() (m []orm.Params) {
	var ch []orm.Params
	o := orm.NewOrm()
	sSql := "SELECT ip,COUNT(ip) FROM t_downloadlog GROUP BY ip"
	num, err := o.Raw(sSql).Values(&ch)
	if err == nil && num > 0 {
		return ch
	}
	beego.Debug(err)
	return nil

}

// AddTDownload insert a new TDownload into database and returns
// last inserted Id on success.
func AddTDownload(m *TDownload) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetTDownloadById retrieves TDownload by Id. Returns error if
// Id doesn't exist
func GetTDownloadById(id int) (v *TDownload, err error) {
	o := orm.NewOrm()
	v = &TDownload{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTDownload retrieves all TDownload matches certain condition. Returns empty list if
// no records exist
func GetAllTDownload(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(TDownload))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []TDownload
	qs = qs.OrderBy(sortFields...)
	if _, err := qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateTDownload updates TDownload by Id and returns error if
// the record to be updated doesn't exist
func UpdateTDownloadById(m *TDownload) (err error) {
	o := orm.NewOrm()
	v := TDownload{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			beego.Debug("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTDownload deletes TDownload by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTDownload(id int) (err error) {
	o := orm.NewOrm()
	v := TDownload{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&TDownload{Id: id}); err == nil {
			beego.Debug("Number of records deleted in database:", num)
		}
	}
	return
}
