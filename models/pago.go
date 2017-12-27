package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Pago struct {
	Id                 int          `orm:"column(id);pk"`
	DetalleLiquidacion int          `orm:"column(detalle_liquidacion)"`
	Valor              float64      `orm:"column(valor)"`
	TipoPago           int          `orm:"column(tipo_pago)"`
	EntidadPago        int          `orm:"column(entidad_pago)"`
	PeriodoPago        *PeriodoPago `orm:"column(periodo_pago);rel(fk)"`
}

type TotalPagosUnidos struct {
	TipoPago       int
	Total          float64
	NombreConcepto string
	AliasConcepto  string
}

type TotalPagosUd struct {
	TipoPago int
	Total    float64
}

type TotalPagos struct {
	TipoPago       int
	Total          float64
	NombreConcepto string
	AliasConcepto  string
}

type TrPeriodoPago struct {
	Contratos   []string
	PeriodoPago *PeriodoPago
	Pagos       []*Pago
}

func (t *Pago) TableName() string {
	return "pago"
}

func init() {
	orm.RegisterModel(new(Pago))
}

// AddPago insert a new Pago into database and returns
// last inserted Id on success.
func AddPago(m *Pago) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPagoById retrieves Pago by Id. Returns error if
// Id doesn't exist
func GetPagoById(id int) (v *Pago, err error) {
	o := orm.NewOrm()
	v = &Pago{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPago retrieves all Pago matches certain condition. Returns empty list if
// no records exist
func GetAllPago(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Pago))
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

	var l []Pago
	qs = qs.OrderBy(sortFields...).RelatedSel("Periodo_Pago")
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

// UpdatePago updates Pago by Id and returns error if
// the record to be updated doesn't exist
func UpdatePagoById(m *Pago) (err error) {
	o := orm.NewOrm()
	v := Pago{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePago deletes Pago by Id and returns error if
// the record to be deleted doesn't exist
func DeletePago(id int) (err error) {
	o := orm.NewOrm()
	v := Pago{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Pago{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
