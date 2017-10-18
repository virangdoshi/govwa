package main

import (
	"fmt"
	"html/template"
	"net/http"
)
import "github.com/gorilla/mux"
import "secureCodingLab/util"
//import validation "secureCodingLab/vulnerability/inputvalidation"
import sqli "secureCodingLab/vulnerability/sqli"
import xss "secureCodingLab/vulnerability/xss"

//input validation
/* func validateHandler(w http.ResponseWriter, r *http.Request) {

	var data = validation.WithNoValidation(r) //default
	if util.CheckLevel(r) {                   //if level == high
		data = validation.WithValidation(r)
	}
	datares := struct {
		Res string
	}{
		Res: data,
	}
	//fmt.Println(checkLevel(r))
	util.SafeRender(w, "validation", datares)
} */

//sql injection and escaping
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	data, err := sqli.UnsafeGetData(r) //default function set to unsafe
	if util.CheckLevel(r) {            //if level == high
		data, err = sqli.SafeGetData(r)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	util.RenderAsJson(w, data)
}

//cros site scripting
func getName(w http.ResponseWriter, r *http.Request) {
	//name := r.URL.Query().Get("name")
	name := r.FormValue("name")
	data := xss.New()
	if !util.CheckLevel(r) { //if level == low
		data.UnEscapeString(name)
	} else {
		data.EscapeString(name)
	}
	fmt.Println(data.Name)

	param := make(map[string]interface{})
	param["name"] = template.HTML(name)

	util.UnSafeRender(w, "template.xss", param)
}

//index and set cookie
func indexHandler(w http.ResponseWriter, r *http.Request) {
	cookie := util.SetCookieLevel(w, r)
	data := make(map[string]interface{})
	data["cookie"] = cookie
	util.SafeRender(w,"template.index", data)
}

func main() {
	s := http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))) //public directory
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	//r.HandleFunc("/validate", validateHandler)
	r.HandleFunc("/getuser", getUserHandler)
	r.HandleFunc("/getinfo", getName)
	r.PathPrefix("/public/").Handler(s)

	fmt.Println("staring server at port 8082")

	http.ListenAndServe(":8082", r)
}
