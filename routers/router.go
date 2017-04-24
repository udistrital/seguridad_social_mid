// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/ss_mid_api/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/upc_adicional",
			beego.NSInclude(
				&controllers.UpcAdicionalController{},
			),
		),

		beego.NSNamespace("/seg_social",
			beego.NSInclude(
				&controllers.SegSocialController{},
			),
		),

		beego.NSNamespace("/seg_social_detalle",
			beego.NSInclude(
				&controllers.SegSocialDetalleController{},
			),
		),

		beego.NSNamespace("/zona",
			beego.NSInclude(
				&controllers.ZonaController{},
			),
		),

		beego.NSNamespace("/tipo_upc",
			beego.NSInclude(
				&controllers.TipoUpcController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
