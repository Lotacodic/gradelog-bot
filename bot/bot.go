package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gradelog-sys/models"
	"gradelog-sys/services"

	tele "gopkg.in/telebot.v3"
)

// Session tracks where a user is in a multi-step conversation
type Session struct {
	Step         string
	FullName     string
	Email        string
	Password     string
	StudentID    string
	StudentName  string
	SemesterName string
	CourseCode   string
	NewCode      string
	Title        string
	CreditUnit   int
}

// sessions stores one session per Telegram user ID
var sessions = map[int64]*Session{}

var (
	mainMenu  = &tele.ReplyMarkup{}
	dashboard = &tele.ReplyMarkup{}

	// Main menu buttons
	btnRegister = mainMenu.Data("📝 Register", "register")
	btnLogin    = mainMenu.Data("🔐 Login", "login")

	// Dashboard buttons
	btnAddSem       = dashboard.Data("➕ Add Semester", "add_semester")
	btnViewSem      = dashboard.Data("📋 View Semesters", "view_semesters")
	btnAddCourse    = dashboard.Data("➕ Add Course", "add_course")
	btnViewCourses  = dashboard.Data("📚 View Courses", "view_courses")
	btnUpdateCourse = dashboard.Data("✏️ Update Course", "update_course")
	btnDeleteCourse = dashboard.Data("🗑️ Delete Course", "delete_course")
	btnGPA          = dashboard.Data("📊 Calculate GPA", "calc_gpa")
	btnCGPA         = dashboard.Data("🏆 Calculate CGPA", "calc_cgpa")
	btnLogout       = dashboard.Data("🚪 Logout", "logout")
)

// getSession returns the session for a user, creating one if it doesn't exist
func getSession(userID int64) *Session {
	if sessions[userID] == nil {
		sessions[userID] = &Session{}
	}
	return sessions[userID]
}

// Start initializes and runs the Telegram bot
func Start(token string) {
	settings := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(settings)
	if err != nil {
		log.Fatal("Failed to start bot:", err)
		return
	}
	fmt.Println("Bot is running...")

	b.Handle("/start", func(c tele.Context) error {
		return c.Send(
			"🎓 *Welcome to GradeLog!*\n"+
				"Your personal academic record manager.\n\n"+
				"Please choose an option:",
			mainMenuKeyboard(),
			tele.ModeMarkdown,
		)
	})

	// Register each button individually
	b.Handle(&btnRegister, func(c tele.Context) error {
		session := getSession(c.Sender().ID)
		session.Step = "register_name"
		return c.Send("📝 *Register*\n\nPlease enter your full name:", tele.ModeMarkdown)
	})

	b.Handle(&btnLogin, func(c tele.Context) error {
		session := getSession(c.Sender().ID)
		session.Step = "login_email"
		return c.Send("🔐 *Login*\n\nPlease enter your email:", tele.ModeMarkdown)
	})

	b.Handle(&btnAddSem, func(c tele.Context) error {
		session := getSession(c.Sender().ID)
		session.Step = "add_semester_name"
		return c.Send("➕ *Add Semester*\n\nEnter semester name (e.g. Year 1 Semester 1):", tele.ModeMarkdown)
	})

	b.Handle(&btnViewSem, func(c tele.Context) error {
		session := getSession(c.Sender().ID)
		return handleViewSemesters(c, session)
	})

	b.Handle(&btnAddCourse, func(c tele.Context) error {
		session := getSession(c.Sender().ID)
		session.Step = "add_course_semester"
		return handleAskSemesterForCourse(c, session)
	})

	b.Handle(&btnDeleteCourse, func(c tele.Context) error {
		session := getSession(c.Sender().ID)
		session.Step = "delete_course_semester"
		return handleAskSemesterForDelete(c, session)
	})

	b.Handle(&btnUpdateCourse, func(c tele.Context) error {
		session := getSession(c.Sender().ID)
		session.Step = "update_course_semester"
		return handleAskSemesterForUpdate(c, session)
	})

	b.Handle(&btnViewCourses, func(c tele.Context) error {
		session := getSession(c.Sender().ID)
		session.Step = "view_courses_semester"
		return handleAskSemesterForView(c, session)
	})

	b.Handle(&btnGPA, func(c tele.Context) error {
		session := getSession(c.Sender().ID)
		session.Step = "calc_gpa_semester"
		return handleAskSemesterForGPA(c, session)
	})

	b.Handle(&btnCGPA, func(c tele.Context) error {
		session := getSession(c.Sender().ID)
		return handleCalcCGPA(c, session)
	})

	b.Handle(&btnLogout, func(c tele.Context) error {
		sessions[c.Sender().ID] = &Session{}
		return c.Send("👋 Logged out successfully!", mainMenuKeyboard())
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		return handleMessage(c, b)
	})

	b.Start()
}

// mainMenuKeyboard returns the main menu buttons
func mainMenuKeyboard() *tele.ReplyMarkup {
	mainMenu.Inline(
		mainMenu.Row(btnRegister),
		mainMenu.Row(btnLogin),
	)
	return mainMenu
}

// dashboardKeyboard returns the logged-in menu buttons
func dashboardKeyboard() *tele.ReplyMarkup {
	dashboard.Inline(
		dashboard.Row(btnAddSem, btnViewSem),
		dashboard.Row(btnAddCourse, btnViewCourses),
		dashboard.Row(btnUpdateCourse, btnDeleteCourse),
		dashboard.Row(btnGPA, btnCGPA),
		dashboard.Row(btnLogout),
	)
	return dashboard
}

// handleMessage processes text input based on session step
func handleMessage(c tele.Context, b *tele.Bot) error {
	userID := c.Sender().ID
	session := getSession(userID)
	text := strings.TrimSpace(c.Text())

	switch session.Step {
	// Registration flow
	case "register_name":
		session.FullName = text
		session.Step = "register_email"
		return c.Send("Enter your email:")

	case "register_email":
		session.Email = text
		session.Step = "register_password"
		return c.Send("Enter your password:")

	case "register_password":
		err := services.RegisterStudent(session.FullName, session.Email, text)
		session.Step = ""
		if err != nil {
			return c.Send("❌ Registration failed: "+err.Error(), mainMenuKeyboard())
		}
		return c.Send("✅ Registration successful! Please login.", mainMenuKeyboard())

		// Login flow
	case "login_email":
		session.Email = text
		session.Step = "login_password"
		return c.Send("Enter your password:")

	case "login_password":
		student, err := services.LoginStudent(session.Email, text)
		session.Step = ""
		if err != nil {
			return c.Send("❌ Login failed: "+err.Error(), mainMenuKeyboard())
		}
		session.StudentID = student.ID
		session.StudentName = student.FullName
		return c.Send(
			fmt.Sprintf("✅ Welcome back, *%s*! What would you like to do?", student.FullName),
			dashboardKeyboard(),
			tele.ModeMarkdown,
		)

		// Add semester flow
	case "add_semester_name":
		err := services.AddSemester(session.StudentID, text)
		session.Step = ""
		if err != nil {
			return c.Send("❌ Failed:"+err.Error(), dashboardKeyboard())
		}
		return c.Send("✅ Semester added successfully!", dashboardKeyboard())

	// Add course flow
	case "add_course_semester":
		session.SemesterName = text
		session.Step = "add_course_code"
		return c.Send("Enter course code (e.g MTH101):")

	case "add_course_code":
		session.CourseCode = text
		session.Step = "add_course_title"
		return c.Send("Enter course title:")

	case "add_course_title":
		session.Title = text
		session.Step = "add_course_credits"
		return c.Send("Enter credit units (1-6):")

	case "add_course_credits":
		val, err := strconv.Atoi(text)
		if err != nil || val < 1 || val > 6 {
			return c.Send("❌ Invalid credit units. Please enter a number between 1 and 6:")
		}
		session.CreditUnit = val
		session.Step = "add_course_score"
		return c.Send("Enter score (0-100):")

	case "add_course_score":
		val, err := strconv.ParseFloat(text, 64)
		if err != nil || val < 0 || val > 100 {
			return c.Send("❌ Invalid score. Please enter a number between 0 and 100:")
		}
		course := models.Course{
			CourseCode:  session.CourseCode,
			Title:       session.Title,
			CreditUnits: session.CreditUnit,
			Score:       val,
		}
		err = services.AddCourse(session.StudentID, session.SemesterName, course)
		session.Step = ""
		if err != nil {
			return c.Send("❌ Failed: "+err.Error(), dashboardKeyboard())
		}
		return c.Send("✅ Course added successfully!", dashboardKeyboard())

	// View courses flow
	case "view_courses_semester":
		return handleViewCourses(c, session, text)

	// GPA flow
	case "calc_gpa_semester":
		return handleCalcGPA(c, session, text)

	// Delete Course Flow
	case "delete_course_semester":
		session.SemesterName = text
		session.Step = "delete_course_code"
		return c.Send("Enter course code(e.g MTH101):")

	case "delete_course_code":
		err := services.DeleteCourse(session.StudentID, session.SemesterName, text)
		if err != nil {
			return c.Send("❌ Failed:"+err.Error(), dashboardKeyboard())
		}
		return c.Send("✅ Course deleted successfully!", dashboardKeyboard())

	// Update course flow
	case "update_course_semester":
		session.SemesterName = text
		session.Step = "update_course_code"
		return c.Send("Enter course code to update(e.g MTH101):")

	case "update_course_code":
		session.CourseCode = text
		session.Step = "update_course_newcode"
		return c.Send("Enter new course code(e.g MTH101):")

	case "update_course_newcode":
		session.NewCode = text
		session.Step = "update_course_title"
		return c.Send("Enter course title:")

	case "update_course_title":
		session.Title = text
		session.Step = "update_course_credits"
		return c.Send("Enter credit units (1-6):")

	case "update_course_credits":
		val, err := strconv.Atoi(text)
		if err != nil || val < 1 || val > 6 {
			return c.Send("❌ Invalid credit units. Please enter a number between 1 and 6:")
		}
		session.CreditUnit = val
		session.Step = "update_course_score"
		return c.Send("Enter score (0-100):")

	case "update_course_score":
		val, err := strconv.ParseFloat(text, 64)
		if err != nil || val < 0 || val > 100 {
			return c.Send("❌ Invalid score. Please enter a number between 0 and 100:")
		}
		updatedCourse := models.Course{
			CourseCode:  session.NewCode,
			Title:       session.Title,
			CreditUnits: session.CreditUnit,
			Score:       val,
		}

		err = services.UpdateCourse(session.StudentID, session.SemesterName, session.CourseCode, updatedCourse)
		session.Step = ""
		if err != nil {
			return c.Send("❌ Failed:"+err.Error(), dashboardKeyboard())
		}
		return c.Send("✅ Course updated successfully!", dashboardKeyboard())

	default:
		return c.Send(
			"Please use the menu buttons or send /start to begin.",
			mainMenuKeyboard(),
		)
	}
}

func handleViewSemesters(c tele.Context, session *Session) error {
	student, err := services.GetStudent(session.StudentID)
	if err != nil {
		return c.Send("You have no semesters yet!", dashboardKeyboard())
	}
	if len(student.Semesters) == 0 {
		return c.Send("You have no semesters yet!", dashboardKeyboard())
	}
	msg := "📋 *Your Semesters:*\n\n"
	for i, sem := range student.Semesters {
		msg += fmt.Sprintf("%d. %s (%d courses)\n", i+1, sem.Name, len(sem.Courses))
	}
	return c.Send(msg, dashboardKeyboard(), tele.ModeMarkdown)
}

func handleAskSemesterForCourse(c tele.Context, session *Session) error {
	student, err := services.GetStudent(session.StudentID)
	if err != nil {
		return c.Send("❌ Error: "+err.Error(), dashboardKeyboard())
	}
	if len(student.Semesters) == 0 {
		session.Step = ""
		return c.Send("You have no semesters yet! Add one first.", dashboardKeyboard())
	}
	msg := "📋 *Your Semesters:*\n\n"
	for i, sem := range student.Semesters {
		msg += fmt.Sprintf("%d. %s\n", i+1, sem.Name)
	}
	msg += "\nEnter the semester name to add a course to:"
	return c.Send(msg, tele.ModeMarkdown)
}

func handleAskSemesterForView(c tele.Context, session *Session) error {
	student, err := services.GetStudent(session.StudentID)
	if err != nil {
		return c.Send("❌ Error: "+err.Error(), dashboardKeyboard())
	}
	if len(student.Semesters) == 0 {
		session.Step = ""
		return c.Send("You have no semesters yet!", dashboardKeyboard())
	}
	msg := "📋 *Your Semesters:*\n\n"
	for i, sem := range student.Semesters {
		msg += fmt.Sprintf("%d. %s\n", i+1, sem.Name)
	}
	msg += "\nEnter the semester name to view courses:"
	return c.Send(msg, tele.ModeMarkdown)
}

func handleAskSemesterForGPA(c tele.Context, session *Session) error {
	student, err := services.GetStudent(session.StudentID)
	if err != nil {
		return c.Send("❌ Error: "+err.Error(), dashboardKeyboard())
	}
	if len(student.Semesters) == 0 {
		session.Step = ""
		return c.Send("You have no semesters yet!", dashboardKeyboard())
	}
	msg := "📋 *Your Semesters:*\n\n"
	for i, sem := range student.Semesters {
		msg += fmt.Sprintf("%d. %s\n", i+1, sem.Name)
	}
	msg += "\nEnter the semester name to calculate GPA:"
	return c.Send(msg, tele.ModeMarkdown)
}

func handleViewCourses(c tele.Context, session *Session, semesterName string) error {
	session.Step = ""
	student, err := services.GetStudent(session.StudentID)
	if err != nil {
		return c.Send("❌ Error: "+err.Error(), dashboardKeyboard())
	}
	for _, sem := range student.Semesters {
		if strings.EqualFold(sem.Name, semesterName) {
			if len(sem.Courses) == 0 {
				return c.Send("No courses in this semester yet!", dashboardKeyboard())
			}
			msg := fmt.Sprintf("📚 *Courses in %s:*\n\n", sem.Name)
			for _, course := range sem.Courses {
				msg += fmt.Sprintf(
					"*%s* - %s\nCredits: %d | Score: %.1f | Grade: %s\n\n",
					course.CourseCode, course.Title,
					course.CreditUnits, course.Score, course.Grade,
				)
			}
			return c.Send(msg, dashboardKeyboard(), tele.ModeMarkdown)
		}
	}
	return c.Send("❌ Semester not found!", dashboardKeyboard())
}

func handleCalcGPA(c tele.Context, session *Session, semesterName string) error {
	session.Step = ""
	gpa, err := services.CalculateGPA(session.StudentID, semesterName)
	if err != nil {
		return c.Send("❌ "+err.Error(), dashboardKeyboard())
	}
	return c.Send(
		fmt.Sprintf("📊 GPA for *%s*: *%.2f / 5.00*", semesterName, gpa),
		dashboardKeyboard(),
		tele.ModeMarkdown,
	)
}

func handleCalcCGPA(c tele.Context, session *Session) error {
	cgpa, err := services.CalculateCGPA(session.StudentID)
	if err != nil {
		return c.Send("❌ "+err.Error(), dashboardKeyboard())
	}
	return c.Send(
		fmt.Sprintf("🏆 Your CGPA: *%.2f / 5.00*", cgpa),
		dashboardKeyboard(),
		tele.ModeMarkdown,
	)
}

func handleAskSemesterForDelete(c tele.Context, session *Session) error {
	student, err := services.GetStudent(session.StudentID)
	if err != nil {
		return c.Send("❌ Error: "+err.Error(), dashboardKeyboard())
	}

	if len(student.Semesters) == 0 {

		session.Step = ""

		return c.Send("You have no semesters yet! Add one first.", dashboardKeyboard())

	}

	msg := "📋 *Your Semesters:*\n\n"

	for i, sem := range student.Semesters {
		msg += fmt.Sprintf("%d. %s\n", i+1, sem.Name)
	}

	msg += "\nEnter the semester name to delete a course  from:"

	return c.Send(msg, tele.ModeMarkdown)
}

func handleAskSemesterForUpdate(c tele.Context, session *Session) error {
	student, err := services.GetStudent(session.StudentID)
	if err != nil {
		return c.Send("❌ Error: "+err.Error(), dashboardKeyboard())
	}

	if len(student.Semesters) == 0 {

		session.Step = ""

		return c.Send("You have no semesters yet! Add one first.", dashboardKeyboard())

	}

	msg := "📋 *Your Semesters:*\n\n"

	for i, sem := range student.Semesters {
		msg += fmt.Sprintf("%d. %s\n", i+1, sem.Name)
	}

	msg += "\nEnter the semester name to update a course in:"

	return c.Send(msg, tele.ModeMarkdown)
}
