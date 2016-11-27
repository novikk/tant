package routers

import (
	"github.com/astaxie/beego"
	"github.com/novikk/tant/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/form", &controllers.MainController{}, "get:ShowForm")
	beego.Router("/set_info", &controllers.MainController{}, "post:SetInfo")
	beego.Router("/login", &controllers.MainController{}, "get:Logins")
	beego.Router("/recommendations", &controllers.MainController{}, "get:Recommendations")
	beego.Router("/login/instagram", &controllers.MainController{}, "get:InstagramLogin")
	beego.Router("/login/instagram/callback", &controllers.MainController{}, "get:InstagramCallback")
	beego.Router("/getLikes", &controllers.MainController{}, "get:GetLikes")
	//beego.Router("/map", &controllers.MainController{}, "get:MapBrands")
	//beego.Router("/suggestBrands", &controllers.MainController{}, "get:GetSuggestedBrands")
	beego.Router("/callZalando", &controllers.MainController{}, "get:CallZalando")
	beego.Router("/getBranches", &controllers.ZalandoController{})
	beego.Router("/login/deezer", &controllers.MainController{}, "get:DeezerLogin")
	beego.Router("/login/deezer/callback", &controllers.MainController{}, "get:DeezerCallback")

}
