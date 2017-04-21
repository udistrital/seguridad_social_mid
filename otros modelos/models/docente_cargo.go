package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type Docente_x_Cargo struct {
	Id                int       `orm:"column(id)"`
	Asignacion_basica int       `orm:"column(asignacion_basica)"`
	FechaInicio       time.Time `orm:"column(emp_desde);type(date);null"`
	FechaFin          time.Time `orm:"column(emp_hasta);type(date);null"`
	Puntos            float64   `orm:"column(puntos)"`
	Regimen           string    `orm:"column(regimen)"`
	Cargo             string    `orm:"column(cargo)"`
}

type ContratoDocente struct {
	Contrato string
}

type CedulaDocente struct {
	Cedula int `orm:"column(contratista)"`
}

func init() {
	orm.RegisterModel(new(Docente_x_Cargo))
}

func GetCedulaDocente(num_contrato string) (res int) {
	o := orm.NewOrm()
	var temp []CedulaDocente
	_, err := o.Raw("SELECT contrato_general.contratista FROM argo.contrato_general WHERE numero_contrato = '" + num_contrato + "' and contrato_general.objeto_contrato ='Docente de planta'").QueryRows(&temp)
	if err == nil {
		fmt.Println("Consulta exitosa")
	}
	return temp[0].Cedula
}

// last inserted Id on success.
func GetAsignacionBasicaDocente(id_proveedor int) (asignacion_basica []Docente_x_Cargo) {
	o := orm.NewOrm()
	var temp []Docente_x_Cargo
	idProveedorString := strconv.Itoa(id_proveedor)
	_, err := o.Raw("SELECT cargo.id,cargo.asignacion_basica, per.emp_desde, per.emp_hasta,informacion_proveedor.puntos,per.regimen,persona_natural.cargo FROM personal.cargo AS cargo,agora.informacion_proveedor, personal.persona as per,agora.informacion_persona_natural as persona_natural WHERE persona_natural.num_documento_persona= informacion_proveedor.num_documento AND per.id =" + idProveedorString + " AND per.estado = 'A' AND per.id_cargo = cargo.id AND informacion_proveedor.id_proveedor=per.id").QueryRows(&temp)
	if err == nil {
		fmt.Println(temp[0].Regimen)
		fmt.Println("Consulta exitosa")
	}

	return temp
}
