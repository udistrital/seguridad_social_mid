package models

import (
  "fmt"
  "strconv"
	"github.com/astaxie/beego/orm"

)

type Funcionario_x_Pritec struct {
    Porcentaje_Prima_Tecnica         int             `orm:"column(valor_porcentaje)"`
}


func init() {
//	orm.RegisterModel(new(Funcionario_x_Pritec))
}

// last inserted Id on success.
func GetPorcentajePT (id_proveedor int)  (porcentajePT int){
	o := orm.NewOrm()
  var temp []int;
  idProveedorString := strconv.Itoa(id_proveedor)
	_, err := o.Raw("SELECT valor_porcentaje FROM personal.porcentaje_prima_tecnica WHERE id_proveedor = "+idProveedorString+" AND estado_porcentaje = 'A';").QueryRows(&temp)
  if err == nil{
      fmt.Println("Consulta exitosa")
  }else{
      temp[0] = 0;
  }

  return temp[0]
}
