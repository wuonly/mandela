package maize

//import (
//	"html/template"
//	"path/filepath"
//)

////var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

//type Template struct {
//	template_url []string
//	templates    template.Must
//}

//func (t *Template) Init(dir string) {
//	t.template_url = make([]string, 0)
//	filepath.Walk(template_PATH, t.addTemplate)
//	template.Must(template.ParseFiles(t.template_url))
//}
//func (t *Template) addTemplate(path string, f os.FileInfo, err error) {
//	t.template_url = append(t.template_url, path)
//}
//func (t *Template) RenderTemplate(w http.ResponseWriter, tmpl string, locals interface{}) {
//	err := templates.ExecuteTemplate(w, tmpl+".html", p)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//	}
//}
