package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

/*Titan crud*/
type InformacionProveedor struct {
	Id                      int                `orm:"column(id_proveedor);pk"`
	Tipopersona             string             `orm:"column(tipopersona)"`
	NumDocumento            float64            `orm:"column(num_documento)"`
	IdCiudadContacto        *Ciudad            `orm:"column(id_ciudad_contacto);rel(fk)"`
	Direccion               string             `orm:"column(direccion)"`
	Correo                  string             `orm:"column(correo)"`
	Web                     string             `orm:"column(web);null"`
	NomAsesor               string             `orm:"column(nom_asesor);null"`
	TelAsesor               string             `orm:"column(tel_asesor);null"`
	Descripcion             string             `orm:"column(descripcion);null"`
	Anexorut                string             `orm:"column(anexorut)"`
	PuntajeEvaluacion       float64            `orm:"column(puntaje_evaluacion);null"`
	ClasificacionEvaluacion string             `orm:"column(clasificacion_evaluacion);null"`
	Estado                  *ParametroEstandar `orm:"column(estado);rel(fk)"`
	TipoCuentaBancaria      string             `orm:"column(tipo_cuenta_bancaria)"`
	NumCuentaBancaria       string             `orm:"column(num_cuenta_bancaria)"`
	IdEntidadBancaria       float64            `orm:"column(id_entidad_bancaria)"`
	FechaRegistro           string             `orm:"column(fecha_registro)"`
	FechaUltimaModificacion string             `orm:"column(fecha_ultima_modificacion)"`
	NomProveedor            string             `orm:"column(nom_proveedor);null"`
}

/* Agora
type InformacionProveedor struct {
	Id                      int                `orm:"column(id_proveedor);pk"`
	Tipopersona             string             `orm:"column(tipopersona)"`
	NumDocumento            float64            `orm:"column(num_documento)"`
	IdCiudadContacto        int                `orm:"column(id_ciudad_contacto)"`
	Direccion               string             `orm:"column(direccion)"`
	Correo                  string             `orm:"column(correo)"`
	Web                     string             `orm:"column(web);null"`
	NomAsesor               string             `orm:"column(nom_asesor);null"`
	TelAsesor               string             `orm:"column(tel_asesor);null"`
	Descripcion             string             `orm:"column(descripcion);null"`
	Anexorut                string             `orm:"column(anexorut)"`
	PuntajeEvaluacion       float64            `orm:"column(puntaje_evaluacion);null"`
	ClasificacionEvaluacion string             `orm:"column(clasificacion_evaluacion);null"`
	Estado                  *ParametroEstandar `orm:"column(estado);rel(fk)"`
	TipoCuentaBancaria      string             `orm:"column(tipo_cuenta_bancaria)"`
	NumCuentaBancaria       string             `orm:"column(num_cuenta_bancaria)"`
	IdEntidadBancaria       float64            `orm:"column(id_entidad_bancaria)"`
	FechaRegistro           string             `orm:"column(fecha_registro)"`
	FechaUltimaModificacion string             `orm:"column(fecha_ultima_modificacion)"`
	NomProveedor            string             `orm:"column(nom_proveedor);null"`
} */

func (t *InformacionProveedor) TableName() string {
	return "informacion_proveedor"
}

func init() {
	orm.RegisterModel(new(InformacionProveedor))
}

// AddInformacionProveedor insert a new InformacionProveedor into database and returns
// last inserted Id on success.
func AddInformacionProveedor(m *InformacionProveedor) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetInformacionProveedorById retrieves InformacionProveedor by Id. Returns error if
// Id doesn't exist
func GetInformacionProveedorById(id int) (v *InformacionProveedor, err error) {
	o := orm.NewOrm()
	v = &InformacionProveedor{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllInformacionProveedor retrieves all InformacionProveedor matches certain condition. Returns empty list if
// no records exist
func GetAllInformacionProveedor(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(InformacionProveedor))
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

	var l []InformacionProveedor
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

// UpdateInformacionProveedor updates InformacionProveedor by Id and returns error if
// the record to be updated doesn't exist
func UpdateInformacionProveedorById(m *InformacionProveedor) (err error) {
	o := orm.NewOrm()
	v := InformacionProveedor{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteInformacionProveedor deletes InformacionProveedor by Id and returns error if
// the record to be deleted doesn't exist
func DeleteInformacionProveedor(id int) (err error) {
	o := orm.NewOrm()
	v := InformacionProveedor{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&InformacionProveedor{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
