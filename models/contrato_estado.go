package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ContratoEstado struct {
	NumeroContrato string          `orm:"column(numero_contrato);null"`
	Vigencia       int             `orm:"column(vigencia);null"`
	FechaRegistro  time.Time       `orm:"column(fecha_registro);type(timestamp without time zone);null"`
	Id             int             `orm:"column(id);pk;auto"`
	Estado         *EstadoContrato `orm:"column(estado);rel(fk)"`
	Usuario        string          `orm:"column(usuario);null"`
}

func (t *ContratoEstado) TableName() string {
	return "contrato_estado"
}

func init() {
	orm.RegisterModel(new(ContratoEstado))
}

// AddContratoEstado insert a new ContratoEstado into database and returns
// last inserted Id on success.
func AddContratoEstado(m *ContratoEstado) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetContratoEstadoById retrieves ContratoEstado by Id. Returns error if
// Id doesn't exist
func GetContratoEstadoById(id int) (v *ContratoEstado, err error) {
	o := orm.NewOrm()
	v = &ContratoEstado{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllContratoEstado retrieves all ContratoEstado matches certain condition. Returns empty list if
// no records exist
func GetAllContratoEstado(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ContratoEstado)).RelatedSel(5)
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
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

	var l []ContratoEstado
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
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

// UpdateContratoEstado updates ContratoEstado by Id and returns error if
// the record to be updated doesn't exist
func UpdateContratoEstadoById(m *ContratoEstado) (err error) {
	o := orm.NewOrm()
	v := ContratoEstado{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteContratoEstado deletes ContratoEstado by Id and returns error if
// the record to be deleted doesn't exist
func DeleteContratoEstado(id int) (err error) {
	o := orm.NewOrm()
	v := ContratoEstado{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ContratoEstado{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
