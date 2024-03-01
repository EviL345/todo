package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	"todo/internal/lib/response"
	"todo/internal/session"
	"todo/internal/storage"
)

func Index(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				log.Println("Ошибка парсинга формы", err)
				render.JSON(w, r, response.Error("Не удалось распарсить форму"))

				return
			}

			title := r.FormValue("title")
			description := r.FormValue("description")

			if title == "" {
				log.Println("Не указано название задачи")
				render.JSON(w, r, response.Error("Не указано название задачи"))

				return
			}

			ctx := r.Context()
			userLogin := ctx.Value("userLogin").(string)
			err := db.AddTask(userLogin, title, description)
			if err != nil {
				log.Println("Ошибка добавления задачи", err)
				render.JSON(w, r, response.Error("Не удалось добавить задачу"))

				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			log.Println("Задача добавлена")
		}

		ctx := r.Context()
		userLogin := ctx.Value("userLogin").(string)
		tasks, err := db.GetAllTasks(userLogin)
		if err != nil {
			log.Println("Ошибка получения задач", err)
			render.JSON(w, r, response.Error("Не удалось получить задачи"))

			return
		}
		templ, err := template.ParseFiles("./frontend/index.html")
		if err != nil {
			log.Println("Ошибка загрузки шаблона", err)
			render.JSON(w, r, response.Error("Не удалось загрузить шаблон"))

			return
		}

		err = templ.Execute(w, tasks)
		if err != nil {
			log.Println("Ошибка выполнения шаблона", err)
			render.JSON(w, r, response.Error("Не удалось выполнить шаблон"))

			return
		}
	}
}

func DeleteTask(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId := chi.URLParam(r, "task_id")
		id, err := strconv.Atoi(taskId)
		if err != nil {
			log.Println("Не удалось получить id задачи", err)
			render.JSON(w, r, response.Error("Не удалось получить id задачи"))

			return
		}

		if err := db.DeleteTask(id); err != nil {
			log.Println("Ошибка удаления задачи", err)
			render.JSON(w, r, response.Error("Не удалось удалить задачу"))

			return
		}

		log.Printf("Задача %d удалена\n", id)
		render.JSON(w, r, response.OK())

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func Login(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Println("Ошибка парсинга формы", err)
			render.JSON(w, r, response.Error("Не удалось распарсить форму"))

			return
		}

		inputLogin := r.FormValue("login")
		inputPassword := r.FormValue("password")

		isExists, err := db.IsUserInDB(inputLogin)
		if err != nil {
			log.Println("Ошибка проверки пользователя", err)
			render.JSON(w, r, response.Error("Не удалось проверить пользователя"))

			return
		}

		if !isExists {
			log.Println("Пользователь не найден")
			render.JSON(w, r, response.Error("Пользователь не найден"))

			return
		}

		password, err := db.GetPassword(inputLogin)
		if err != nil {
			log.Println("Ошибка получения пароля", err)
			render.JSON(w, r, response.Error("Не удалось получить пароль"))

			return
		}

		if password != inputPassword {
			log.Println("Пользователь ввёл неверный пароль")
			render.JSON(w, r, response.Error("Неверный пароль"))

			return
		}

		sessionID := session.NewSession(inputLogin)
		var cookie = http.Cookie{
			Name:    "session_id",
			Value:   sessionID,
			Expires: time.Now().Add(24 * time.Hour),
		}

		log.Println("Пользователь авторизован")
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Register(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Println("Ошибка парсинга формы", err)
			render.JSON(w, r, response.Error("Не удалось распарсить форму"))

			return
		}

		inputLogin := r.FormValue("login")
		inputPassword := r.FormValue("password")

		if err := db.AddUser(inputLogin, inputPassword); err != nil {
			log.Println("Ошибка добавления пользователя", err)
			render.JSON(w, r, response.Error("Не удалось добавить пользователя"))

			return
		}

		sessionID := session.NewSession(inputLogin)
		var cookie = http.Cookie{
			Name:    "session_id",
			Value:   sessionID,
			Expires: time.Now().Add(24 * time.Hour),
		}

		http.SetCookie(w, &cookie)

		log.Println("Пользователь добавлен")

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
