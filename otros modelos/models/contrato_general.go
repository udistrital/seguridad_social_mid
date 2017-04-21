package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ContratoGeneral struct {
	Id                           string                 `orm:"column(numero_contrato);pk"`
	Vigencia                     int                 `orm:"column(vigencia)"`
	ObjetoContrato               string              `orm:"column(objeto_contrato);null"`
	PlazoEjecucion               int                 `orm:"column(plazo_ejecucion)"`
	FormaPago                    *Parametros         `orm:"column(forma_pago);rel(fk)"`
	OrdenadorGasto               *ArgoOrdenadores    `orm:"column(ordenador_gasto);rel(fk)"`
	ClausulaRegistroPresupuestal bool                `orm:"column(clausula_registro_presupuestal);null"`
	SedeSolicitante              string              `orm:"column(sede_solicitante);null"`
	DependenciaSolicitante       string              `orm:"column(dependencia_solicitante);null"`
	NumeroSolicitudNecesidad     int                 `orm:"column(numero_solicitud_necesidad)"`
	NumeroCdp                    int                 `orm:"column(numero_cdp)"`
	Contratista                  *InformacionProveedor   `orm:"column(contratista);rel(fk)"`
	UnidadEjecucion              *Parametros         `orm:"column(unidad_ejecucion);rel(fk)"`
	ValorContrato                float64             `orm:"column(valor_contrato)"`
	Justificacion                string              `orm:"column(justificacion)"`
	DescripcionFormaPago         string              `orm:"column(descripcion_forma_pago)"`
	Condiciones                  string              `orm:"column(condiciones)"`
	UnidadEjecutora              *UnidadEjecutora    `orm:"column(unidad_ejecutora);rel(fk)"`
	FechaRegistro                time.Time           `orm:"column(fecha_registro);type(date)"`
	TipologiaContrato            int                 `orm:"column(tipologia_contrato)"`
	TipoCompromiso               int                 `orm:"column(tipo_compromiso)"`
	ModalidadSeleccion           int                 `orm:"column(modalidad_seleccion)"`
	Procedimiento                int                 `orm:"column(procedimiento)"`
	RegimenContratacion          int                 `orm:"column(regimen_contratacion)"`
	TipoGasto                    int                 `orm:"column(tipo_gasto)"`
	TemaGastoInversion           int                 `orm:"column(tema_gasto_inversion)"`
	OrigenPresupueso             int                 `orm:"column(origen_presupueso)"`
	OrigenRecursos               int                 `orm:"column(origen_recursos)"`
	TipoMoneda                   int                 `orm:"column(tipo_moneda)"`
	ValorContratoMe              float64             `orm:"column(valor_contrato_me);null"`
	ValorTasaCambio              float64             `orm:"column(valor_tasa_cambio);null"`
	TipoControl                  int                 `orm:"column(tipo_control);null"`
	Observaciones                string              `orm:"column(observaciones);null"`
	Supervisor                   *SupervisorContrato `orm:"column(supervisor);rel(fk)"`
	ClaseContratista             int                 `orm:"column(clase_contratista)"`
	Convenio                     string              `orm:"column(convenio);null"`
	NumeroConstancia             int                 `orm:"column(numero_constancia);null"`
	Estado                       bool                `orm:"column(estado);null"`
	ResgistroPresupuestal        int                 `orm:"column(resgistro_presupuestal);null"`
	TipoContrato                 *TipoContrato       `orm:"column(tipo_contrato);rel(fk)"`
	LugarEjecucion               *LugarEjecucion     `orm:"column(lugar_ejecucion);rel(fk)"`
}

func (t *ContratoGeneral) TableName() string {
	return "contrato_general"
}

func init() {
	orm.RegisterModel(new(ContratoGeneral))
}

// AddContratoGeneral insert a new ContratoGeneral into database and returns
// last inserted Id on success.
func AddContratoGeneral(m *ContratoGeneral) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetContratoGeneralById retrieves ContratoGeneral by Id. Returns error if
// Id doesn't exist
func GetContratoGeneralById(id string) (v *ContratoGeneral, err error) {
	o := orm.NewOrm()
	v = &ContratoGeneral{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllContratoGeneral retrieves all ContratoGeneral matches certain condition. Returns empty list if
// no records exist
func GetAllContratoGeneral(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ContratoGeneral))
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

	var l []ContratoGeneral
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

// UpdateContratoGeneral updates ContratoGeneral by Id and returns error if
// the record to be updated doesn't exist
func UpdateContratoGeneralById(m *ContratoGeneral) (err error) {
	o := orm.NewOrm()
	v := ContratoGeneral{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteContratoGeneral deletes ContratoGeneral by Id and returns error if
// the record to be deleted doesn't exist
func DeleteContratoGeneral(id string) (err error) {
	o := orm.NewOrm()
	v := ContratoGeneral{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ContratoGeneral{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
