package models

type Course struct {
	CourseCode  string  `json:"course_code"`
	Title       string  `json:"title"`
	CreditUnits int     `json:"credit_units"`
	Score       float64 `json:"score"`
	Grade       string  `json:"grade"`
}

type Semester struct {
	Name    string   `json:"name"`
	Courses []Course `json:"courses"`
}

type Student struct {
	ID        string     `json:"id"`
	FullName  string     `json:"full_name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Semesters []Semester `json:"semesters"`
}

type Database struct {
	Students []Student `json:"students"`
}
