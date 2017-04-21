package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"],
		beego.ControllerComments{
			Method:           "CalcularSegSocial",
			Router:           `CalcularSegSocial/:id`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialDetalleController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialDetalleController"],
		beego.ControllerComments{
			Method:           "GetNovedadesPorPersona",
			Router:           `GetNovedadesPorPersona/:persona`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:FuncionarioProveedorController"] = append(beego.GlobalControllerRouter["github.com/miguelramirez93/titan_api_crud/controllers:FuncionarioProveedorController"],
		beego.ControllerComments{
			Method:           "ConsultarIDProveedor",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			Params:           nil})

	beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"] = append(beego.GlobalControllerRouter["ss_api_mid/controllers:SegSocialController"],
		beego.ControllerComments{
			Method:           "GetConceptosIbc",
			Router:           `ConceptosIbc`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})
}
