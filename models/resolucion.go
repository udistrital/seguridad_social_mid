package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Resolucion struct {
	Objeto                  string          `orm:"column(objeto);null"`
	FechaRegistro           time.Time       `orm:"column(fecha_registro);type(date)"`
	Estado                  bool            `orm:"column(estado)"`
	ConsideracionResolucion string          `orm:"column(consideracion_resolucion)"`
	PreambuloResolucion     string          `orm:"column(preambulo_resolucion)"`
	IdTipoResolucion        *TipoResolucion `orm:"column(id_tipo_resolucion);rel(fk)"`
	IdDependencia           int             `orm:"column(id_dependencia)"`
	Vigencia                int             `orm:"column(vigencia)"`
	FechaExpedicion         *time.Time      `orm:"column(fecha_expedicion);type(date);null"`
	NumeroResolucion        string          `orm:"column(numero_resolucion)"`
	Id                      int             `orm:"column(id_resolucion);pk;auto"`
}

func (t *Resolucion) TableName() string {
	return "resolucion"
}

func init() {
	orm.RegisterModel(new(Resolucion))
}

func CancelarResolucion(m *Resolucion) (err error) {
	o := orm.NewOrm()
	o.Begin()
	v := ResolucionVinculacionDocente{Id: m.Id}
	if err = o.Read(&v); err == nil {
		var vinculacion_docente []*VinculacionDocente
		_, err = o.QueryTable("vinculacion_docente").Filter("id_resolucion", m.Id).Filter("estado", true).All(&vinculacion_docente)
		for _, vd := range vinculacion_docente {
			var contratos_generales []*ContratoGeneral
			if vd.NumeroContrato != "" && vd.Vigencia != 0 {
				_, err = o.QueryTable("contrato_general").Filter("numero_contrato", vd.NumeroContrato).Filter("vigencia", vd.Vigencia).All(&contratos_generales)
				if err == nil {
					for _, c := range contratos_generales {
						aux1 := c.Id
						aux2 := c.VigenciaContrato
						e := ContratoEstado{}
						e.NumeroContrato = aux1
						e.Vigencia = aux2
						e.FechaRegistro = time.Now()
						e.Estado = &EstadoContrato{Id: 7}
						if _, err = o.Insert(&e); err != nil {
							o.Rollback()
							return
						}
					}
				} else {
					o.Rollback()
					return
				}
			}
		}
		var num int64
		if num, err = o.Update(m); err == nil {
			var e ResolucionEstado
			e.Resolucion = m
			e.Estado = &EstadoResolucion{Id: 3}
			e.FechaRegistro = time.Now()
			_, err = o.Insert(&e)
			if err == nil {
				fmt.Println("Number of records updated in database:", num)
			} else {
				o.Rollback()
				return
			}
		} else {
			o.Rollback()
			return
		}
	} else {
		o.Rollback()
		return
	}
	o.Commit()
	return
}

func GenerarResolucion(m *Resolucion) (id int64, err error) {
	o := orm.NewOrm()
	o.Begin()
	m.Vigencia, _, _ = time.Now().Date()
	m.FechaRegistro = time.Now()
	m.Estado = true
	m.IdTipoResolucion = &TipoResolucion{Id: 1}
	id, err = o.Insert(m)
	if err == nil {
		var e ResolucionEstado
		e.Resolucion = m
		e.Estado = &EstadoResolucion{Id: 1}
		e.FechaRegistro = time.Now()
		_, err = o.Insert(&e)
		if err != nil {
			o.Rollback()
			return
		}
	} else {
		o.Rollback()
		return
	}
	o.Commit()
	return
}

func RestaurarResolucion(m *Resolucion) (err error) {
	o := orm.NewOrm()
	o.Begin()
	var num int64
	if num, err = o.Update(m); err == nil {
		var e ResolucionEstado
		e.Resolucion = m
		e.Estado = &EstadoResolucion{Id: 1}
		e.FechaRegistro = time.Now()
		_, err = o.Insert(&e)
		if err == nil {
			fmt.Println("Number of records updated in database:", num)
		} else {
			o.Rollback()
			return
		}
	} else {
		o.Rollback()
		return
	}
	o.Commit()
	return
}

// AddResolucion insert a new Resolucion into database and returns
// last inserted Id on success.
func AddResolucion(m *Resolucion) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetResolucionById retrieves Resolucion by Id. Returns error if
// Id doesn't exist
func GetResolucionById(id int) (v *Resolucion, err error) {
	o := orm.NewOrm()
	v = &Resolucion{Id: id}
	if err = o.Read(v); err == nil {
		fmt.Println("si")
		return v, nil
	}
	fmt.Println("si")
	return nil, err
}

// GetAllResolucion retrieves all Resolucion matches certain condition. Returns empty list if
// no records exist
func GetAllResolucion(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Resolucion)).RelatedSel(5)
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

	var l []Resolucion
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

// UpdateResolucion updates Resolucion by Id and returns error if
// the record to be updated doesn't exist
func UpdateResolucionById(m *Resolucion) (err error) {
	o := orm.NewOrm()
	v := Resolucion{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteResolucion deletes Resolucion by Id and returns error if
// the record to be deleted doesn't exist
func DeleteResolucion(id int) (err error) {
	o := orm.NewOrm()
	v := Resolucion{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Resolucion{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
