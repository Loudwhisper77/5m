package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Student struct {
	Name string
}

type Event struct {
	Title string
	Date  time.Time
}

type NewsItem struct {
	Title   string
	Content string
	Date    time.Time
}

type ClassData struct {
	Photos      []string
	Events      []Event
	Grades      map[string]int
	Schedule    []string
	TopStudents []Student
	News        []NewsItem
	Students    []Student
}

var classData = ClassData{
	Photos: []string{
		"/static/photo1.jpg",
		"/static/photo2.jpg",
		"/static/photo3.jpg",
		"/static/my_photo.jpg",
	},
	Events: []Event{
		{Title: "Пока что ничего", Date: time.Date(2025, time.April, 20, 10, 0, 0, 0, time.UTC)},
		{Title: "Мероприятие не указано", Date: time.Date(2025, time.May, 5, 14, 0, 0, 0, time.UTC)},
	},
	Grades: map[string]int{
		"Оценок пока нет...": 0,
	},
	Schedule: []string{
		"Понедельник: РОВ, Математика, Русский Язык, Литература, Биология, ОДНКНР",
		"Вторник: Труд, Труд, Математика, Английский Язык, Русский язык, География",
		"Среда: Русский язык, Литература, Математика, История, Английский язык, Логика",
		"Четверг: Математика, Математика, ИЗО, Физ-ра, Русский язык, Литература, Спортивные игры",
		"Пятница: Информатика/Английский язык, История, Английский язык/Информатика, Физ-ра, Математика, Русский язык",
	},
	TopStudents: []Student{
		{Name: "Бекасова Елена Владимировна: Классный руководитель"},
		{Name: "Куштина Виктория: Староста класса"},
	},
	News: []NewsItem{
		{Title: "Новостей пока нет", Content: "В ближайшее время они могут появиться", Date: time.Date(2025, time.April, 10, 12, 0, 0, 0, time.UTC)},
		{Title: "Новостей пока нет", Content: "В ближайшее время они могут появиться", Date: time.Date(2025, time.April, 12, 15, 0, 0, 0, time.UTC)},
		{Title: "Новостей пока нет", Content: "В ближайшее время они могут появиться", Date: time.Date(2025, time.April, 12, 15, 0, 0, 0, time.UTC)},
	},
	Students: []Student{
		{Name: "Абрамова Виктория"},
		{Name: "Абдулаев Никита"},
		{Name: "Британ Александр"},
		{Name: "Горбатов Михаил"},
		{Name: "Грачёва Виктория"},
		{Name: "Додонов Михаил"},
		{Name: "Дюдяков Фёдор"},
		{Name: "Ермачкова Мария"},
		{Name: "Ермолаев Фёдор"},
		{Name: "Жданкин Юрий"},
		{Name: "Животов Егор"},
		{Name: "Иванова Карина"},
		{Name: "Кирюшкин Николай"},
		{Name: "Кострица Маргарита"},
		{Name: "Кузьмин Георгий"},
		{Name: "Купцевич Дмитрий"},
		{Name: "Куштина Виктория"},
		{Name: "Мокеева Анна"},
		{Name: "Михеев Мирон"},
		{Name: "Никульшин Максим"},
		{Name: "Ордиховская Полина"},
		{Name: "Парамонова Виктория"},
		{Name: "Потапкина Виктория"},
		{Name: "Ромашова Виктория"},
		{Name: "Сафронова Любовь"},
		{Name: "Сердечкина Вера"},
		{Name: "Старча Андрей"},
		{Name: "Степанова Екатерина"},
		{Name: "Солонарь Михаил"},
		{Name: "Улитина Дарья"},
	},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetRequest(w, r)
		case http.MethodPost:
			handlePostRequest(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Сервер запущен на http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{

		"ClassData": classData,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	timestamp := time.Now().Format("20060102150405")
	ext := getFileExtension(header.Filename)
	newFileName := fmt.Sprintf("/static/%s_%s%s", "uploaded", timestamp, ext)

	f, err := createUploadedFile(newFileName)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	classData.Photos = append(classData.Photos, newFileName)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createUploadedFile(name string) (io.WriteCloser, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func getFileExtension(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[i:]
		}
	}
	return ""
}

const (
	username = "5MclassTOP"
	password = "5Mschool13Site"
)

// basicAuthMiddleware - middleware для базовой авторизации
func basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		auth := strings.SplitN(authHeader, " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || pair[0] != username || pair[1] != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Если авторизация прошла успешно, передаем запрос дальше
		next.ServeHTTP(w, r)
	})
}
