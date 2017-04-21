package models

import (
  "fmt"
  "strconv"
	"github.com/astaxie/beego/orm"
  	"time"
)

type Funcionario_x_Cargo struct {
  Id                      int                `orm:"column(id)"`
  Asignacion_basica         int             `orm:"column(asignacion_basica)"`
  FechaInicio    time.Time `orm:"column(emp_desde);type(date);null"`
	FechaFin       time.Time `orm:"column(emp_hasta);type(date);null"`
}


func init() {
	orm.RegisterModel(new(Funcionario_x_Cargo))
}

// last inserted Id on success.
func GetAsignacionBasica (id_proveedor int)  (asignacion_basica []Funcionario_x_Cargo){
	o := orm.NewOrm()
  var temp []Funcionario_x_Cargo;
  idProveedorString := strconv.Itoa(id_proveedor)
	_, err := o.Raw("SELECT cargo.id,cargo.asignacion_basica, per.emp_desde, per.emp_hasta FROM personal.cargo AS cargo, personal.persona as per WHERE per.id = "+idProveedorString+" AND per.estado = 'A' AND per.id_cargo = cargo.id;").QueryRows(&temp)
  if err == nil{
      fmt.Println(temp[0].FechaInicio)
      fmt.Println("Consulta exitosa")
  }

  return temp
}
