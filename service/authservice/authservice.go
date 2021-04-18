package authservice


import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/prometheus/common/log"
	"github.com/xiet16/authcenter/common/lib"
	"github.com/xiet16/authcenter/dao"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"net/http"
)

func GetOAuthServerAndManager() (*server.Server ,*manage.Manager,error){
	// config oauth manager
	mgr := manage.NewDefaultManager()
	mgr.SetAuthorizeCodeTokenCfg(manage.DefaultPasswordTokenCfg)

	//token store in memory
	//mgr.MustTokenStorage(store.NewMemoryTokenStore())
	// or use redis token store
	//mgr.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
	//	Addr:     addr,
	//	Password: pwd,
	//}))

	//get redis addr by rand
	addr, pwd, err := lib.RedisConfFactory("default")
	if err != nil {
		log.Error("read redis config error:", err)
		return nil,nil, err
	}
	mgr.MapTokenStorage(NewRedisStore(&redis.Options{
		Addr:     addr,
		Password: pwd,
	}))

	//access token generate method:jwt
	mgr.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))

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
	srv := server.NewServer(server.NewConfig(), mgr)
	srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	srv.SetAuthorizeScopeHandler(authorizeScopeHandler)
	srv.SetInternalErrorHandler(internalErrorHandler)
	srv.SetResponseErrorHandler(responseErrorHandler)

	return srv,mgr,nil
}

func passwordAuthorizationHandler(username, password string) (userID string, err error) {
	user := &dao.User{
		Name:     username,
		Password: password,
	}
	out, err := user.GetUserIDByPwd(user)
	if err != nil {
		return "", err
	}
	return out.Name, nil
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

func responseErrorHandler(re *errors.Response) {
	fmt.Println("Response Error:", re.Error.Error())
}
