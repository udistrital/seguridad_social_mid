package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type VinculacionDocente struct {
	Id                   int                           `orm:"column(id);pk;auto"`
	NumeroContrato       string                        `orm:"column(numero_contrato);null"`
	Vigencia             int                           `orm:"column(vigencia);null"`
	IdPersona            int                           `orm:"column(id_persona)"`
	NumeroHorasSemanales int                           `orm:"column(numero_horas_semanales)"`
	NumeroSemanas        int                           `orm:"column(numero_semanas)"`
	IdPuntoSalarial      int                           `orm:"column(id_punto_salarial);null"`
	IdSalarioMinimo      int                           `orm:"column(id_salario_minimo);null"`
	IdResolucion         *ResolucionVinculacionDocente `orm:"column(id_resolucion);rel(fk)"`
	IdDedicacion         *Dedicacion                   `orm:"column(id_dedicacion);rel(fk)"`
	IdProyectoCurricular int16                         `orm:"column(id_proyecto_curricular)"`
	Estado               bool                          `orm:"column(estado)"`
	FechaRegistro        time.Time                     `orm:"column(fecha_registro);type(date)"`
}

func (t *VinculacionDocente) TableName() string {
	return "vinculacion_docente"
}

func init() {
	orm.RegisterModel(new(VinculacionDocente))
}

func AddConjuntoVinculaciones(m []VinculacionDocente) (err error) {
	o := orm.NewOrm()
	o.Begin()
	for _, vinculacion := range m {
		vinculacion.Estado = true
		vinculacion.FechaRegistro = time.Now()
		_, err = o.Insert(&vinculacion)
		if err != nil {
			o.Rollback()
			return
		}
	}
	o.Commit()
	return
}

// AddVinculacionDocente insert a new VinculacionDocente into database and returns
// last inserted Id on success.
func AddVinculacionDocente(m *VinculacionDocente) (id int64, err error) {
	o := orm.NewOrm()
	m.Estado = true
	m.FechaRegistro = time.Now()
	id, err = o.Insert(m)
	return
}

// GetVinculacionDocenteById retrieves VinculacionDocente by Id. Returns error if
// Id doesn't exist
func GetVinculacionDocenteById(id int) (v *VinculacionDocente, err error) {
	o := orm.NewOrm()
	v = &VinculacionDocente{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllVinculacionDocente retrieves all VinculacionDocente matches certain condition. Returns empty list if
// no records exist
func GetAllVinculacionDocente(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(VinculacionDocente)).RelatedSel()
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

	var l []VinculacionDocente
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

// UpdateVinculacionDocente updates VinculacionDocente by Id and returns error if
// the record to be updated doesn't exist
func UpdateVinculacionDocenteById(m *VinculacionDocente) (err error) {
	o := orm.NewOrm()
	v := VinculacionDocente{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		m.FechaRegistro = v.FechaRegistro
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteVinculacionDocente deletes VinculacionDocente by Id and returns error if
// the record to be deleted doesn't exist
func DeleteVinculacionDocente(id int) (err error) {
	o := orm.NewOrm()
	v := VinculacionDocente{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&VinculacionDocente{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
