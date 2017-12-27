package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type DetallePreliquidacion struct {
	Id 								int                    `orm:"auto;column(id);pk"`
	ValorCalculado     float64               `orm:"column(valor_calculado)"`
	NumeroContrato     string                `orm:"column(numero_contrato);null"`
	VigenciaContrato   int                   `orm:"column(vigencia_contrato);null"`
	DiasLiquidados     float64               `orm:"column(dias_liquidados);null"`
	TipoPreliquidacion *TipoPreliquidacion   `orm:"column(tipo_preliquidacion);rel(fk)"`
	Preliquidacion     *Preliquidacion       `orm:"column(preliquidacion);rel(fk)"`
	Concepto           *ConceptoNomina       `orm:"column(concepto);rel(fk)"`
}

func (t *DetallePreliquidacion) TableName() string {
	return "detalle_preliquidacion"
}

func init() {
	orm.RegisterModel(new(DetallePreliquidacion))
}

// AddDetallePreliquidacion insert a new DetallePreliquidacion into database and returns
// last inserted Id on success.
func AddDetallePreliquidacion(m *DetallePreliquidacion) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	fmt.Println(err)
	return
}

// GetDetallePreliquidacionById retrieves DetallePreliquidacion by Id. Returns error if
// Id doesn't exist
func GetDetallePreliquidacionById(id int) (v *DetallePreliquidacion, err error) {
	o := orm.NewOrm()
	v = &DetallePreliquidacion{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllDetallePreliquidacion retrieves all DetallePreliquidacion matches certain condition. Returns empty list if
// no records exist
func GetAllDetallePreliquidacion(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(DetallePreliquidacion)).RelatedSel(5)
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

	var l []DetallePreliquidacion
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

// UpdateDetallePreliquidacion updates DetallePreliquidacion by Id and returns error if
// the record to be updated doesn't exist
func UpdateDetallePreliquidacionById(m *DetallePreliquidacion) (err error) {
	o := orm.NewOrm()
	v := DetallePreliquidacion{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDetallePreliquidacion deletes DetallePreliquidacion by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDetallePreliquidacion(id int) (err error) {
	o := orm.NewOrm()
	res, err := o.Raw("DELETE FROM detalle_preliquidacion WHERE preliquidacion = ?", id).Exec()
	if err == nil {
    num, _ := res.RowsAffected()
    fmt.Println("row affected nums: ", num)
	}
	return
}
