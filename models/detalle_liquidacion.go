package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type DetalleLiquidacion struct {
	Id              int          `orm:"column(id);pk"`
	ValorCalculado  int64        `orm:"column(valor_calculado)"`
	EstadoConcepto  string       `orm:"column(estado_concepto)"`
	Liquidacion     *Liquidacion `orm:"column(liquidacion);rel(fk)"`
	Persona         int          `orm:"column(persona)"`
	Concepto        *Concepto    `orm:"column(concepto);rel(fk)"`
	NumeroContrato  string       `orm:"column(numero_contrato)"`
	DiasLiquidados  string       `orm:"column(dias_liquidados)"`
	TipoLiquidacion string       `orm:"column(tipo_liquidacion)"`
}

func (t *DetalleLiquidacion) TableName() string {
	return "detalle_liquidacion"
}

func init() {
	orm.RegisterModel(new(DetalleLiquidacion))
}

// AddDetalleLiquidacion insert a new DetalleLiquidacion into database and returns
// last inserted Id on success.
func AddDetalleLiquidacion(m *DetalleLiquidacion) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetDetalleLiquidacionById retrieves DetalleLiquidacion by Id. Returns error if
// Id doesn't exist
func GetDetalleLiquidacionById(id int) (v *DetalleLiquidacion, err error) {
	o := orm.NewOrm()
	v = &DetalleLiquidacion{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllDetalleLiquidacion retrieves all DetalleLiquidacion matches certain condition. Returns empty list if
// no records exist
func GetAllDetalleLiquidacion(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(DetalleLiquidacion))
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

	var l []DetalleLiquidacion
	qs = qs.OrderBy(sortFields...).RelatedSel("numero_contrato")
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

// UpdateDetalleLiquidacion updates DetalleLiquidacion by Id and returns error if
// the record to be updated doesn't exist
func UpdateDetalleLiquidacionById(m *DetalleLiquidacion) (err error) {
	o := orm.NewOrm()
	v := DetalleLiquidacion{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDetalleLiquidacion deletes DetalleLiquidacion by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDetalleLiquidacion(id int) (err error) {
	o := orm.NewOrm()
	v := DetalleLiquidacion{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&DetalleLiquidacion{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
