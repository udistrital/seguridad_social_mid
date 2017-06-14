package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ConceptoPorPersona struct {
	ValorNovedad  float64               `orm:"column(valor_novedad)"`
	EstadoNovedad string                `orm:"column(estado_novedad)"`
	FechaDesde    time.Time             `orm:"column(fecha_desde);type(date)"`
	FechaHasta    time.Time             `orm:"column(fecha_hasta);type(date)"`
	NumCuotas     int64                 `orm:"column(num_cuotas)"`
	Persona       *InformacionProveedor `orm:"column(persona);rel(fk)"`
	Concepto      *Concepto             `orm:"column(concepto);rel(fk)"`
	Nomina        *Nomina               `orm:"column(nomina);rel(fk)"`
	Id            int                   `orm:"auto;column(id);pk"`
	Tipo          string                `orm:"column(tipo);null"`
}

func (t *ConceptoPorPersona) TableName() string {
	return "concepto_por_persona"
}

func init() {
	orm.RegisterModel(new(ConceptoPorPersona))
}

// AddConceptoPorPersona insert a new ConceptoPorPersona into database and returns
// last inserted Id on success.
func AddConceptoPorPersona(m *ConceptoPorPersona) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetConceptoPorPersonaById retrieves ConceptoPorPersona by Id. Returns error if
// Id doesn't exist
func GetConceptoPorPersonaById(id int) (v *ConceptoPorPersona, err error) {
	o := orm.NewOrm()
	v = &ConceptoPorPersona{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllConceptoPorPersona retrieves all ConceptoPorPersona matches certain condition. Returns empty list if
// no records exist
func GetAllConceptoPorPersona(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ConceptoPorPersona))
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

	var l []ConceptoPorPersona
	qs = qs.OrderBy(sortFields...).RelatedSel(5)
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

// UpdateConceptoPorPersona updates ConceptoPorPersona by Id and returns error if
// the record to be updated doesn't exist
func UpdateConceptoPorPersonaById(m *ConceptoPorPersona) (err error) {
	o := orm.NewOrm()
	v := ConceptoPorPersona{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteConceptoPorPersona deletes ConceptoPorPersona by Id and returns error if
// the record to be deleted doesn't exist
func DeleteConceptoPorPersona(id int) (err error) {
	o := orm.NewOrm()
	v := ConceptoPorPersona{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ConceptoPorPersona{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func ConceptoPorPersonaActivo(id int, v *Preliquidacion) (datos []ConceptoPorPersona, err error) {
	o := orm.NewOrm()
	consulta := `select a.*
								from titan.concepto_por_persona as a
											where a.estado_novedad = 'Activo'
											and a.persona = ?
											and a.nomina = ?
											`
	_, err = o.Raw(consulta, id, v.Nomina.Id).QueryRows(&datos)
	var concepto []Concepto
	for i := 0; i < len(datos); i++ {
		consulta := `select a.* from titan.concepto  as a
									inner join titan.concepto_por_persona as b on a.id = b.concepto
									where b.id = ?`
		_, err = o.Raw(consulta, datos[i].Id).QueryRows(&concepto)
		datos[i].Concepto = &concepto[0]
	}

	return
}
