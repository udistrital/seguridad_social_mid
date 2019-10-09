package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:ConceptosIbcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:ConceptosIbcController"],
        beego.ControllerComments{
            Method: "ActualizarConceptos",
            Router: `/ActualizarConceptos/:tipo_ibc`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:GeneradorRelgasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:GeneradorRelgasController"],
        beego.ControllerComments{
            Method: "ObtenerHechosCalculo",
            Router: `/ObtenerHechosCalculo`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:GeneradorRelgasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:GeneradorRelgasController"],
        beego.ControllerComments{
            Method: "RegistrarNuevosHechos",
            Router: `/RegistrarNuevosHechos`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:IncapacidadesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:IncapacidadesController"],
        beego.ControllerComments{
            Method: "BuscarPersonas",
            Router: `/BuscarPersonas/:documento`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:IncapacidadesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:IncapacidadesController"],
        beego.ControllerComments{
            Method: "IncapacidadesPorPersona",
            Router: `/incapacidadesPersona/:contrato/:vigencia`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"],
        beego.ControllerComments{
            Method: "CalcularSegSocial",
            Router: `/CalcularSegSocial/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"],
        beego.ControllerComments{
            Method: "CalcularSegSocialHonorarios",
            Router: `/CalcularSegSocialHonorarios/:idPreliquidacion`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"],
        beego.ControllerComments{
            Method: "ConceptosIbc",
            Router: `/ConceptosIbc`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"],
        beego.ControllerComments{
            Method: "GetInfoCabecera",
            Router: `/GetInfoCabecera/:idPreliquidacion/:tipoPlanilla`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"],
        beego.ControllerComments{
            Method: "NovedadesPorPersona",
            Router: `/NovedadesPorPersona/:persona`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"],
        beego.ControllerComments{
            Method: "RegistrarPagos",
            Router: `/RegistrarPagos`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"],
        beego.ControllerComments{
            Method: "SumarPagosSalud",
            Router: `/SumarPagosSalud/:idPeriodoPago`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"],
        beego.ControllerComments{
            Method: "ConceptosIbcParafiscales",
            Router: `/conceptos_ibc_parafiscales`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PlanillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PlanillasController"],
        beego.ControllerComments{
            Method: "GenerarPlanillaActivos",
            Router: `/GenerarPlanillaActivos/:limit/:offset`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:UtilsController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:UtilsController"],
        beego.ControllerComments{
            Method: "GetActualDate",
            Router: `/GetActualDate`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
