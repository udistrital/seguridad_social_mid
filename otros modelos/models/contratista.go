package models

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type Contratista struct {
	Id                int `orm:"column(id_proveedor)"`
	NumeroContrato    int `orm:"column(numero_contrato)"`
	Asignacion_basica int `orm:"column(asignacion_basica)"`
}

func GetContratista(idProveedorString int) (v []Contratista) {
	o := orm.NewOrm()
	var temp []Contratista
	id_proveedor := strconv.Itoa(idProveedorString)
	_, err := o.Raw("SELECT contrato_general.contratista FROM argo.contrato_general WHERE numero_contrato = '" + id_proveedor + "' and contrato_general.objeto_contrato ='Contratista'").QueryRows(&temp)
	if err == nil {
		fmt.Println("Consulta exitosa")
	}

	return temp
}

//informacionproveedor.id_proveedor, contratos.contratista,contratos.numero_contrato,informacionproveedor.nom_proveedor

// last inserted Id on success.
