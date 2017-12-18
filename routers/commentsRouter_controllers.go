package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:AportanteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:AportanteController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:AportanteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:AportanteController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:AportanteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:AportanteController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:AportanteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:AportanteController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:AportanteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:AportanteController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"],
		beego.ControllerComments{
			Method: "RegistrarPagos",
			Router: `/RegistrarPagos`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"],
		beego.ControllerComments{
			Method: "SumarPagosSalud",
			Router: `SumarPagosSalud/:idPeriodoPago`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PeriodoPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PeriodoPagoController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PeriodoPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PeriodoPagoController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PeriodoPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PeriodoPagoController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PeriodoPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PeriodoPagoController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PeriodoPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PeriodoPagoController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PlanillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PlanillasController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PlanillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PlanillasController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PlanillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PlanillasController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PlanillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PlanillasController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PlanillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PlanillasController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:RangoEdadUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:RangoEdadUpcController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:RangoEdadUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:RangoEdadUpcController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:RangoEdadUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:RangoEdadUpcController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:RangoEdadUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:RangoEdadUpcController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:RangoEdadUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:RangoEdadUpcController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoPagoController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:TipoUpcController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:UpcAdicionalController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaUpcController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaUpcController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaUpcController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaUpcController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaUpcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:ZonaUpcController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

}
