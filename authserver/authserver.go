package authserver

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/prometheus/common/log"
	"github.com/xiet16/authcenter/dao"
	"github.com/xiet16/authcenter/lib"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"net/http"
	"net/url"
	"html/template"
	"time"
)

var srv *server.Server
var mgr *manage.Manager

type TplData struct {
	Client *lib.ConnClientConf
	// 用户申请resource scope
	Scope []lib.Scope
	Error string
}

func Run() {

	// config oauth manager
	mgr =  manage.NewDefaultManager()
	mgr.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	//token store
	mgr.MustTokenStorage(store.NewMemoryTokenStore())
	// or use redis token store
	// mgr.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
	//     Addr: config.Get().Redis.Default.Addr,
	//     DB: config.Get().Redis.Default.Db,
	// }))

	//access token generate method:jwt
	mgr.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"),jwt.SigningMethodHS512))


   clientStore := store.NewClientStore()
	for _, v := range lib.ConfConnCientMap.List {
		clientStore.Set(v.ID, &models.Client{
			ID:     v.ID,
			Secret: v.Secret,
			//Domain: v.Domain,
		})
	}

   mgr.MapClientStorage(clientStore)


   // config oauth2 server
   srv = server.NewServer(server.NewConfig(),mgr)
   srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	srv.SetAuthorizeScopeHandler(authorizeScopeHandler)
	srv.SetInternalErrorHandler(internalErrorHandler)
	srv.SetResponseErrorHandler(responseErrorHandler)

    //授权接口
	http.HandleFunc("/authorize", authorizeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

    //获取access_token 接口
	http.HandleFunc("/token", tokenHandler)

	http.HandleFunc("/authenticate", authenticateHandle)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Fatal(http.ListenAndServe(":9096", nil))
}

func authorizeHandler(w http.ResponseWriter, r *http.Request) {
    var form url.Values
    if v,_ := lib.GetSession(r,"session_name");v!=nil{
    	r.ParseForm()
    	if r.Form.Get("client_id") ==""{
			form = v.(url.Values)
		}
	}
	r.Form = form

	if err:= lib.DeleteSession(w,r,"RequestForm");err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	if err:=srv.HandleAuthorizeRequest(w,r);err!=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
	}
}

func loginHandler1(w http.ResponseWriter, r *http.Request) {
	form,err := lib.GetSession(r,"RequestForm")
	if err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	if form == nil {
		http.Error(w,"Invalid Request",http.StatusBadRequest)
	}

	clientID := form.(url.Values).Get("client_id")
	scope := form.(url.Values).Get("scope")

	//页面数据
	data :=TplData{
		Client:lib.GetClient(clientID),
		Scope: lib.ScopeFilter(clientID,scope),
	}

	if data.Scope == nil {
		http.Error(w,"Invalid Scope", http.StatusBadRequest)
	}

	if r.Method == "POST" {
		if r.Form == nil {
			if err:=r.ParseForm(); err!=nil {
				http.Error(w,err.Error(),http.StatusInternalServerError)
				return
			}
		}
	}

	var userID string

	//账号密码验证
	if r.Form.Get("type") == "password" {
		search := &dao.User{
			Name: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}
		user ,err:= search.GetUserIDByPwd(search)
		if err!=nil || user.Name == ""{
			t, err := template.ParseFiles("tpl/login.html")
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data.Error = "用户名密码错误!"
			t.Execute(w, data)
		}
	}

	// 扫码验证
	// 手机验证码验证

	if err:=lib.SetSession(w,r,"LoggedInUserID",userID);err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.Header().Set("location","/authorize")
	w.WriteHeader(http.StatusFound)

	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    //form,err := lib.GetSession(r,"RequestForm")
	//if err!=nil {
	//	http.Error(w,err.Error(),http.StatusInternalServerError)
	//	return
	//}
    r.ParseForm()
	form := r.PostForm

	if form == nil {
		http.Error(w,"Invalid Request",http.StatusBadRequest)
		return
	}

/*	clientID := form.(url.Values).Get("client_id")
	scope := form.(url.Values).Get("scope")*/

	clientID :=  form.Get("client_id")
	scope := r.Form.Get("scope")

	//页面数据
	data :=TplData{
		Client:lib.GetClient(clientID),
		Scope: lib.ScopeFilter(clientID,scope),
	}

	if data.Scope == nil {
		http.Error(w,"Invalid Scope", http.StatusBadRequest)
	}

	if r.Method == "POST" {
		if r.Form == nil {
			if err:=r.ParseForm(); err!=nil {
				http.Error(w,err.Error(),http.StatusInternalServerError)
				return
			}
		}
	}

	var userID string

	//账号密码验证
	if r.Form.Get("type") == "password" {
		search := &dao.User{
			Name: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}
		user ,err:= search.GetUserIDByPwd(search)
		if err!=nil || user.Name == ""{
			t, err := template.ParseFiles("tpl/login.html")
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data.Error = "用户名密码错误!"
			t.Execute(w, data)
		}
	}

	// 扫码验证
	// 手机验证码验证

	if err:=lib.SetSession(w,r,"LoggedInUserID",userID);err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	req, err := srv.ValidationAuthorizeRequest(r)
	if err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	ti, err := srv.GetAuthorizeToken(req)
	if err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	code := srv.GetAuthorizeData(req.ResponseType, ti)
	fmt.Println(code)
	//w.Header().Set("location","/authorize")
	//w.WriteHeader(http.StatusFound)

	return
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Form ==nil {
		if err:=r.ParseForm();err!=nil {
			http.Error(w,err.Error(),http.StatusBadRequest)
			return
		}
	}

	redirectURI := r.Form.Get("redirect_uri")
	if _,err :=url.Parse(redirectURI);err!=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
	}

	if err:=lib.DeleteSession(w,r,"LoggedInUserID");err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", redirectURI)
	w.WriteHeader(http.StatusFound)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	err := srv.HandleTokenRequest(w,r)
	if err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
}

func authenticateHandle(w http.ResponseWriter, r *http.Request) {
	token ,err := srv.ValidationBearerToken(r)
	if err!=nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cli,err := mgr.GetClient(token.GetClientID())
	if err!=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	data := map[string] interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"user_id": token.GetUserID(),
		"client_id": token.GetClientID(),
		"scope": token.GetScope(),
		"domain": cli.GetDomain(),
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(data)
}

func passwordAuthorizationHandler(username, password string) (userID string, err error)  {
    user:= &dao.User{
    	Name: username,
    	Password: password,
	}
   out,err := user.GetUserIDByPwd(user)
	if err!=nil {
		return "", err
	}
   return out.Name,nil
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	v, _ := lib.GetSession(r, "LoggedInUserID")
	if v == nil {
		if r.Form == nil {
			r.ParseForm()
		}
		lib.SetSession(w, r, "RequestForm", r.Form)

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)

		return
	}
	userID = v.(string)

	// 不记住用户
	// store.Delete("LoggedInUserID")
	// store.Save()

	return
}

func authorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	if r.Form == nil {
		r.ParseForm()
	}
	s := lib.ScopeFilter(r.Form.Get("client_id"), r.Form.Get("scope"))
	if s == nil {
		http.Error(w, "Invalid Scope", http.StatusBadRequest)
		return
	}
	scope = lib.ScopeJoin(s)

	return
}

func internalErrorHandler(err error) (re *errors.Response) {
	fmt.Println("Internal Error:", err.Error())
	return
}

func responseErrorHandler(re *errors.Response)  {
	fmt.Println("Response Error:", re.Error.Error())
}