package main

import (
	"Personal-Web/connection"
	"Personal-Web/middleware"
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	connection.DatabaseConnect()

	e := echo.New()

	// static file from directory
	e.Static("/public", "public")
	e.Static("/upload", "upload")

	// initialitation to use session
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	// routing
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/showProject", project)
	e.POST("/addProject", addProject)
	e.GET("/projectDetail/:id", projectDetail)
	e.GET("/deleteProject/:id", deleteProject)
	e.GET("/editProject/:id", editProject)
	e.PUT("/editProject/:id", resultProject)
	e.POST("/addProject", middleware.UploadFile(addProject))

	e.GET("/showregister", showRegister)
	e.GET("/showlogin", showLogin)
	e.GET("/logout", logout)

	e.POST("/register", register)
	e.POST("/login", login)

	fmt.Println("server started on port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}

type Project struct {
	ID          int
	ProjectName string
	Duration    string
	StartDate   time.Time
	EndDate     time.Time
	Description string
	Tech        []string
	Image       string
}

type ProjectR struct {
	ID          int
	ProjectName string
	Duration    string
	StartDate   string
	EndDate     string
	Description string
	Tech        []string
	Image       string
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type SessionData struct {
	IsLogin bool
	Name    string
}

// var userData = SessionData{}

type M map[string]interface{}

func home(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
	}

	project, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, tech, image FROM tb_projects")

	var result []Project

	for project.Next() {
		var each = Project{}
		var arrayTech pgtype.VarcharArray

		err := project.Scan(&each.ID, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Description, &arrayTech, &each.Image)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
		}

		each.Tech = make([]string, len(arrayTech.Elements))

		for i, e := range arrayTech.Elements {
			if e.String == "NodeJs" {
				each.Tech[i] = "fa-brands fa-node"
			}
			if e.String == "ReactJs" {
				each.Tech[i] = "fa-brands fa-react"
			}
			if e.String == "Golang" {
				each.Tech[i] = "fa-brands fa-golang"
			}
			if e.String == "Python" {
				each.Tech[i] = "fa-brands fa-python"
			}
		}

		each.Duration = distanceDate(each.StartDate, each.EndDate)

		result = append(result, each)
	}

	// sess, _ := session.Get("session", c)

	// if sess.Values["isLogin"] != true {
	// 	userData.IsLogin = false
	// } else {
	// 	userData.IsLogin = sess.Values["isLogin"].(bool)
	// 	userData.Name = sess.Values["name"].(string)
	// }

	projects := M{
		// "FlashStatus":  sess.Values["isLogin"],
		// "FlashMessage": sess.Values["message"],
		// "FlashName":    sess.Values["name"],
		"Project": result,
		// "DataSession":  userData,
	}

	// delete(sess.Values, "message")
	// delete(sess.Values, "status")
	// sess.Save(c.Request(), c.Response())

	return tmpl.Execute(c.Response(), projects)
}

func contact(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/contact-me.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func project(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/addProject.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func projectDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	tmpl, err := template.ParseFiles("views/add-project-detail.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
	}

	var projectdetail = Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, tech, image FROM tb_projects WHERE id = $1", id).Scan(&projectdetail.ID, &projectdetail.ProjectName, &projectdetail.StartDate, &projectdetail.EndDate, &projectdetail.Description, &projectdetail.Tech, &projectdetail.Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
	}

	var resultDate = distanceDate(projectdetail.StartDate, projectdetail.EndDate)

	var show = ProjectR{}

	for i, e := range projectdetail.Tech {
		if e == "NodeJs" {
			projectdetail.Tech[i] = "fa-brands fa-node"
		}
		if e == "ReactJs" {
			projectdetail.Tech[i] = "fa-brands fa-react"
		}
		if e == "Golang" {
			projectdetail.Tech[i] = "fa-brands fa-golang"
		}
		if e == "Python" {
			projectdetail.Tech[i] = "fa-brands fa-python"
		}
	}

	show = ProjectR{
		ProjectName: projectdetail.ProjectName,
		Duration:    resultDate,
		StartDate:   dateFormat(projectdetail.StartDate, ""),
		EndDate:     dateFormat(projectdetail.EndDate, ""),
		Description: projectdetail.Description,
		Image:       projectdetail.Image,
		Tech:        projectdetail.Tech,
	}

	projectData := M{
		"Project": show,
	}
	return tmpl.Execute(c.Response(), projectData)
}

func dateFormat(date time.Time, dateType string) string {
	if dateType == "RFC822" {
		return date.Format("02 Jan 2006")
	} else {
		return fmt.Sprintf("%d-%02d-%02d", date.Year(), int(date.Month()), date.Day())
	}
}

func addProject(c echo.Context) error {
	err := c.Request().ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	projectName := c.FormValue("projectName")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	desc := c.FormValue("desc")
	image := c.FormValue("uploadImage")
	nodeJs := c.FormValue("nodeJs")
	reactJs := c.FormValue("reactJs")
	golang := c.FormValue("golang")
	python := c.FormValue("python")

	var tech []string

	//sesuai value
	if nodeJs == "NodeJs" {
		tech = append(tech, "NodeJs")
	}
	if reactJs == "ReactJs" {
		tech = append(tech, "ReactJs")
	}
	if golang == "Golang" {
		tech = append(tech, "Golang")
	}
	if python == "Python" {
		tech = append(tech, "Python")
	}

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, start_date, end_date, description, tech, image) VALUES($1,$2,$3,$4,$5,$6)", projectName, startDate, endDate, desc, tech, image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func editProject(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/edit-Project.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
	}

	var updateProject = Project{}

	id, _ := strconv.Atoi(c.Param("id"))

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, tech, image FROM tb_projects WHERE id=$1", id).Scan(&updateProject.ID, &updateProject.ProjectName, &updateProject.StartDate, &updateProject.EndDate, &updateProject.Description, &updateProject.Tech, &updateProject.Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
	}

	for i, e := range updateProject.Tech {
		if e == "NodeJs" {
			updateProject.Tech[i] = "fa-brands fa-node"
		}
		if e == "ReactJs" {
			updateProject.Tech[i] = "fa-brands fa-react"
		}
		if e == "Golang" {
			updateProject.Tech[i] = "fa-brands fa-golang"
		}
		if e == "Python" {
			updateProject.Tech[i] = "fa-brands fa-python"
		}
	}

	show := ProjectR{
		ProjectName: updateProject.ProjectName,
		StartDate:   dateFormat(updateProject.StartDate, ""),
		EndDate:     dateFormat(updateProject.EndDate, ""),
		Description: updateProject.Description,
		Tech:        updateProject.Tech,
		Image:       updateProject.Image,
	}

	projectData := M{
		"Project": show,
	}
	return tmpl.Execute(c.Response(), projectData)
}

func resultProject(c echo.Context) error {
	err := c.Request().ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	projectName := c.FormValue("projectName")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	desc := c.FormValue("desc")
	nodeJs := c.FormValue("nodeJs")
	reactJs := c.FormValue("reactJs")
	golang := c.FormValue("golang")
	python := c.FormValue("python")
	image := c.FormValue("uploadImage")

	var tech []string

	//sesuai value
	if nodeJs == "NodeJs" {
		tech = append(tech, "NodeJs")
	}
	if reactJs == "ReactJs" {
		tech = append(tech, "ReactJs")
	}
	if golang == "Golang" {
		tech = append(tech, "Golang")
	}
	if python == "Python" {
		tech = append(tech, "Python")
	}

	id, _ := strconv.Atoi(c.Param("id"))
	_, err = connection.Conn.Exec(context.Background(), "UPDATE public.tb_projects SET id=$1 name=$2, start_date=$3, end_date=$4, description=$5, tech=$6 WHERE id=$1", id, projectName, startDate, endDate, desc, tech, image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func distanceDate(startDate time.Time, endDate time.Time) string {
	// variable for count duration
	diff := endDate.Sub(startDate)

	// declare variable with package math method floor
	var yearDistance float64 = math.Floor(float64(diff.Milliseconds()) / (12 * 30 * 24 * 60 * 60 * 1000))
	var monthDistance float64 = math.Floor(float64(diff.Milliseconds()) / (30 * 24 * 60 * 60 * 1000))
	var weekDistance float64 = math.Floor(float64(diff.Milliseconds()) / (7 * 24 * 60 * 60 * 1000))
	var dayDistance float64 = math.Floor(float64(diff.Milliseconds()) / (24 * 60 * 60 * 1000))

	// validation duration
	if yearDistance > 0 {
		year := fmt.Sprintf("%d year", int(yearDistance))
		return year
	} else {
		if monthDistance > 0 {
			month := fmt.Sprintf("%d month", int(monthDistance))
			return month
		} else {
			if weekDistance > 0 {
				week := fmt.Sprintf("%d week", int(weekDistance))
				return week
			} else {
				if dayDistance > 0 {
					day := fmt.Sprintf("%d day", int(dayDistance))
					return day
				} else {
					return ""
				}
			}
		}
	}
}

func showRegister(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func showLogin(c echo.Context) error {
	sess, _ := session.Get("session", c)
	flash := M{
		"FlashStatus":  sess.Values["alertStatus"],
		"FlashMessage": sess.Values["message"],
	}

	delete(sess.Values, "message")
	delete(sess.Values, "alertStatus")

	tmpl, err := template.ParseFiles("views/login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), flash)
}

func login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return redirect(c, "Email Salah !", false, "/showlogin")
	}

	fmt.Println(user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return redirect(c, "Password Salah !", false, "/showlogin")
	}

	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 // 3 jam
	sess.Values["message"] = "Login Success"
	sess.Values["status"] = true // show alert
	sess.Values["name"] = user.Name
	sess.Values["id"] = user.ID
	sess.Values["isLogin"] = true // access login
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func register(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := c.FormValue("name")
	email := c.FormValue("email")
	pass := c.FormValue("pass")

	//generate password
	passHash, _ := bcrypt.GenerateFromPassword([]byte(pass), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user (name, email, password) VALUES($1,$2,$3)", name, email, passHash)

	if err != nil {
		redirect(c, "Register failed, please try again", false, "/showregister")
	}

	return redirect(c, "Register Success", true, "/showlogin")
}

func redirect(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, path)
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Values["isLogin"] = false
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
