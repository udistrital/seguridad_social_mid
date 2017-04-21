package models

import (
	"time"
	"fmt"
	"errors"
	"strings"
	"strconv"
	"reflect"
	"github.com/astaxie/beego/orm"
)

type InformacionPensionado struct {
	Id                   int                   `orm:"column(id);pk"`
	InformacionProveedor int 										`orm:"column(informacion_proveedor)"`
	Estado               string                `orm:"column(estado)"`
	FechaNacEmpleado     time.Time             `orm:"column(fecha_nac_empleado);type(date);null"`
	EstadoCivil          *EstadoCivil          `orm:"column(estado_civil);rel(fk)"`
	FechaRetiro          time.Time             `orm:"column(fecha_retiro);type(date);null"`
	PersonaFallecido     string                `orm:"column(persona_fallecido)"`
	PensionadoEnExterior string                `orm:"column(pensionado_en_exterior)"`
	TipoPensionado       int      						 `orm:"column(tipo_pensionado)"`
	TipoPension          *TipoPension          `orm:"column(tipo_pension);rel(fk)"`
	ValorPensionAsignada int                   `orm:"column(valor_pension_asignada);null"`
	FechaPension         time.Time             `orm:"column(fecha_pension);type(date);null"`
	Resolucion           *Resolucion           `orm:"column(resolucion);rel(fk)"`
	EstadoPension        string                `orm:"column(estado_pension);null"`
}

func (t *InformacionPensionado) TableName() string {
	return "informacion_pensionado"
}

func init() {
	orm.RegisterModel(new(InformacionPensionado))
}

// AddInformacionPensionado insert a new InformacionPensionado into database and returns
// last inserted Id on success.
func AddInformacionPensionado(m *InformacionPensionado) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetInformacionPensionadoById retrieves InformacionPensionado by Id. Returns error if
// Id doesn't exist
func GetPensionado( idProveedorString int) (v []InformacionPensionado) {
	o := orm.NewOrm()
	var temp []  InformacionPensionado
	id_proveedor := strconv.Itoa(idProveedorString)
	_, err := o.Raw("SELECT informacion_proveedor, valor_pension_asignada, pensionado_en_exterior, persona_fallecido,tipo_pensionado FROM personal.informacion_persona_pensionado WHERE informacion_proveedor ="+id_proveedor).QueryRows(&temp)
	if err == nil {
		fmt.Println("Consulta exitosa")
	}
	return temp
}

// GetAllInformacionPensionado retrieves all InformacionPensionado matches certain condition. Returns empty list if
// no records exist
func GetAllInformacionPensionado(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(InformacionPensionado))
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

	var l []InformacionPensionado
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

// UpdateInformacionPensionado updates InformacionPensionado by Id and returns error if
// the record to be updated doesn't exist
func UpdateInformacionPensionadoById(m *InformacionPensionado) (err error) {
	o := orm.NewOrm()
	v := InformacionPensionado{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteInformacionPensionado deletes InformacionPensionado by Id and returns error if
// the record to be deleted doesn't exist
func DeleteInformacionPensionado(id int) (err error) {
	o := orm.NewOrm()
	v := InformacionPensionado{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&InformacionPensionado{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
