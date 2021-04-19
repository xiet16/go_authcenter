package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"github.com/xiet16/authcenter/common/lib"
	"github.com/xiet16/authcenter/dao"
	"github.com/xiet16/authcenter/service/authservice"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

type TplData struct {
	Client *lib.ConnClientConf
	// 用户申请resource scope
	Scope []lib.Scope
	Error string
}

var srv *server.Server
var mgr *manage.Manager
type IndentityServerContoller struct {}

//api register
func IndentityServerRegister(group *gin.RouterGroup)  {

	authServer,authManger,err := authservice.GetOAuthServerAndManager()
	if err!=nil {
		log.Info("config auhtservice err:",err)
	}
	srv = authServer
	mgr = authManger

	indentityController := &IndentityServerContoller{}
	 group.POST("/authorize",indentityController.authorizeHandler)
	 group.POST("/token",indentityController.tokenHandler)
	 group.POST("/login",indentityController.loginHandler)
	 group.POST("/logout",indentityController.logoutHandler)
	 group.POST("/authenticate",indentityController.authenticateHandler)
}

func (identity *IndentityServerContoller)loginHandler(c *gin.Context) {
	r:=c.Request
	w:=c.Writer
	user := &dao.User{
		Name: "daibo",
		Password: "daibo123",
	}
	user,_= user.GetUserIDByPwd(user)

	user.Update(user)
	//form,err := lib.GetSession(r,"RequestForm")
	//if err!=nil {
	//	http.Error(w,err.Error(),http.StatusInternalServerError)
	//	return
	//}

	r.ParseForm()
	form := r.PostForm

	if form == nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	/*	clientID := form.(url.Values).Get("client_id")
		scope := form.(url.Values).Get("scope")*/

	clientID := form.Get("client_id")
	scope := r.Form.Get("scope")

	//页面数据
	data := TplData{
		Client: lib.GetClient(clientID),
		Scope:  lib.ScopeFilter(clientID, scope),
	}

	if data.Scope == nil {
		http.Error(w, "Invalid Scope", http.StatusBadRequest)
	}

	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	var userID string

	//账号密码验证
	if r.Form.Get("type") == "password" {
		search := &dao.User{
			Name:     r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}
		user, err := search.GetUserIDByPwd(search)
		if err != nil || user.Name == "" {
			t, err := template.ParseFiles("tpl/login.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data.Error = "用户名密码错误!"
			t.Execute(w, data)
		}
	}

	// 扫码验证
	// 手机验证码验证

	if err := lib.SetSession(w, r, "LoggedInUserID", userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := srv.ValidationAuthorizeRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	ti, err := srv.GetAuthorizeToken(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	code := srv.GetAuthorizeData(req.ResponseType, ti)
	fmt.Println(code)
	//w.Header().Set("location","/authorize")
	//w.WriteHeader(http.StatusFound)

	return
}

func (identity *IndentityServerContoller)logoutHandler(c *gin.Context) {
	r:=c.Request
	w:=c.Writer
	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	redirectURI := r.Form.Get("redirect_uri")
	if _, err := url.Parse(redirectURI); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := lib.DeleteSession(w, r, "LoggedInUserID"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", redirectURI)
	w.WriteHeader(http.StatusFound)
}

func (identity *IndentityServerContoller)tokenHandler(c *gin.Context) {
	r:=c.Request
	w:=c.Writer
	err := srv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (identity *IndentityServerContoller)authenticateHandler(c *gin.Context) {
	r:=c.Request
	w:=c.Writer
	token, err := srv.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cli, err := mgr.GetClient(token.GetClientID())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"user_id":    token.GetUserID(),
		"client_id":  token.GetClientID(),
		"scope":      token.GetScope(),
		"domain":     cli.GetDomain(),
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(data)
}

func (identity *IndentityServerContoller)authorizeHandler(c *gin.Context) {
	r:=c.Request
	w:=c.Writer
	var form url.Values
	if v, _ := lib.GetSession(r, "session_name"); v != nil {
		r.ParseForm()
		if r.Form.Get("client_id") == "" {
			form = v.(url.Values)
		}
	}
	r.Form = form
	if err := lib.DeleteSession(w, r, "RequestForm"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := srv.HandleAuthorizeRequest(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
