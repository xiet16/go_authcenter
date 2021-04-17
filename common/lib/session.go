package lib

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"net/http"
	"net/url"
)

var store *sessions.CookieStore


func InitSession()  {
	gob.Register(url.Values{})
	//这里要从配置文件取
	store =sessions.NewCookieStore([]byte("abc123"))
	store.Options =&sessions.Options{
		Path: "/",
		MaxAge: 60*20,
		HttpOnly: true,
	}
}

func GetSession(r *http.Request,name string)(val interface{},err error) {
	session,err := store.Get(r,"session_name")
	if err!=nil {
		return
	}

	val = session.Values[name]
	return
}

func SetSession(w http.ResponseWriter,r *http.Request,name string,val interface{}) error {
	session,err := store.Get(r,"session_name")
	if err!=nil {
		return err
	}

	session.Values[name] = val
	err =session.Save(r,w)
	return err
}

func DeleteSession(w http.ResponseWriter,r *http.Request,name string) error {
	session ,err := store.Get(r,"session_name")
	if err!=nil {
		return err
	}

	delete(session.Values,name)
	err = session.Save(r,w)
	return err
}