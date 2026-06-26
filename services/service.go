package services

import (
	"errors"
	"strings"

	"gradelog-sys/models"
	"gradelog-sys/storage"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterStudent(fullName, email, password string) error {
	// Load the current database
	db, err := storage.LoadDB()
	if err != nil {
		return err
	}

	email = strings.ToLower(email)
	// Check if email already exists
	for _, student := range db.Students {
		if strings.EqualFold(student.Email, email) {
			return errors.New("a student with this email already exists")
		}
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create the new student
	newStudent := models.Student{
		ID:        generateID(),
		FullName:  fullName,
		Email:     email,
		Password:  string(hashedPassword),
		Semesters: []models.Semester{},
	}

	// Append and save
	db.Students = append(db.Students, newStudent)
	return storage.SaveDB(db)
}

func LoginStudent(email, password string) (models.Student, error) {
	// Load the database
	db, err := storage.LoadDB()
	if err != nil {
		return models.Student{}, err
	}

	email = strings.ToLower(email)
	// Find the student by email
	for _, student := range db.Students {
		if strings.EqualFold(student.Email, email) {
			// Compare the typed password against the stored hash
			err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(password))
			if err != nil {
				// Hash didn't match — wrong password
				return models.Student{}, errors.New("incorrect password")
			}
			// Password matched — return the student
			return student, nil
		}
	}

	// No student found with that email
	return models.Student{}, errors.New("no account found with this email")
}

func generateID() string {
	return uuid.New().String()
}

// AddSemester adds a new semester to a student's record
func AddSemester(studentID, semesterName string) error {
	// Load the database
	db, err := storage.LoadDB()
	if err != nil {
		return err
	}

	// Find the student by ID using the index
	for i, student := range db.Students {
		if student.ID == studentID {
			// Check if semester already exist
			for _, sem := range student.Semesters {
				if strings.EqualFold(sem.Name, semesterName) {
					return errors.New("semester already exists")
				}
			}

			// Create and add the new semester
			newSemester := models.Semester{
				Name:    semesterName,
				Courses: []models.Course{},
			}

			db.Students[i].Semesters = append(db.Students[i].Semesters, newSemester)

			// Save the database
			return storage.SaveDB(db)
		}
	}
	return errors.New("student not found")
}

// GetStudent retrieves a student by their ID
func GetStudent(studentID string) (models.Student, error) {
	db, err := storage.LoadDB()
	if err != nil {
		return models.Student{}, err
	}

	for _, student := range db.Students {
		if student.ID == studentID {
			return student, nil
		}
	}
	return models.Student{}, errors.New("student not found")
}

// calculateGrade returns the letter grade for a given score
func calculateGrade(score float64) string {
	switch {
	case score >= 70:
		return "A"
	case score >= 60:
		return "B"
	case score >= 50:
		return "C"
	case score >= 45:
		return "D"
	case score >= 40:
		return "E"
	default:
		return "F"

	}
}

// AddCourse adds a new course to a specific semester for a student
func AddCourse(studentID, semesterName string, course models.Course) error {
	// Load the Database
	db, err := storage.LoadDB()
	if err != nil {
		return err
	}

	// Find the student by ID
	for i, student := range db.Students {
		if student.ID == studentID {
			// Find the semester by name
			for j, sem := range student.Semesters {
				if strings.EqualFold(sem.Name, semesterName) {
					// Check for duplicate course code
					for _, c := range sem.Courses {
						if c.CourseCode == course.CourseCode {
							return errors.New("course already exists in this semester")
						}
					}

					// Calculate grade automatically from score
					course.Grade = calculateGrade(course.Score)

					// Add course and save
					db.Students[i].Semesters[j].Courses = append(db.Students[i].Semesters[j].Courses, course)
					return storage.SaveDB(db)
				}
			}
			return errors.New("semester not found")
		}
	}
	return errors.New("student not found")
}

func gradeToPoint(grade string) float64 {
	switch grade {
	case "A":
		return 5.0
	case "B":
		return 4.0
	case "C":
		return 3.0
	case "D":
		return 2.0
	case "E":
		return 1.0
	default:
		return 0.0
	}
}

// CalculateGPA calculates the GPA for a specific semester
func CalculateGPA(studentID, semesterName string) (float64, error) {
	// Get the student
	student, err := GetStudent(studentID)
	if err != nil {
		return 0, err
	}

	// Find the semester
	for _, sem := range student.Semesters {
		if strings.EqualFold(sem.Name, semesterName) {
			// Check if there are  courses to calculate
			if len(sem.Courses) == 0 {
				return 0, errors.New("no courses found in this semester")
			}

			// Calculate total quality points and total credits
			totalQualityPoints := 0.0
			totalCreditUnits := 0

			for _, course := range sem.Courses {
				point := gradeToPoint(course.Grade)
				totalQualityPoints += float64(course.CreditUnits) * point
				totalCreditUnits += course.CreditUnits
			}

			// Compute and return GPA
			gpa := totalQualityPoints / float64(totalCreditUnits)
			return gpa, nil
		}
	}
	return 0, errors.New("semester not found")
}

// CalculateCGPA calculates the CGPA across all semester
func CalculateCGPA(studentID string) (float64, error) {
	//  Get the student
	student, err := GetStudent(studentID)
	if err != nil {
		return 0, err
	}

	// Check if there are semesters
	if len(student.Semesters) == 0 {
		return 0, errors.New("no semesters found")
	}

	// Accumulate across ALL semesters
	totalQualityPoints := 0.0
	totalCreditUnits := 0

	for _, sem := range student.Semesters {
		for _, course := range sem.Courses {
			point := gradeToPoint(course.Grade)
			totalQualityPoints += float64(course.CreditUnits) * point
			totalCreditUnits += course.CreditUnits
		}
	}

	if totalCreditUnits == 0 {
		return 0, errors.New("no courses found across any semester")
	}

	cgpa := totalQualityPoints / float64(totalCreditUnits)
	return cgpa, nil
}

// CRUDE OPERATION

// DeleteSemester removes a semester and all its courses from a student's record
func DeleteSemester(studentID, semesterName string) error {
	db, err := storage.LoadDB()
	if err != nil {
		return err
	}

	for i, student := range db.Students {
		if student.ID == studentID {
			// Find the semster index
			for j, sem := range student.Semesters {
				if strings.EqualFold(sem.Name, semesterName) {
					// Remove semester at index j using append trick
					db.Students[i].Semesters = append(db.Students[i].Semesters[:j], db.Students[i].Semesters[j+1:]...)
					return storage.SaveDB(db)
				}
			}
			return errors.New("semester not found")
		}
	}
	return errors.New("student not found")
}

// DeleteCourse removes a specific course from a semester
func DeleteCourse(studentID, semesterName, courseCode string) error {
	db, err := storage.LoadDB()
	if err != nil {
		return err
	}

	for i, student := range db.Students {
		if student.ID == studentID {
			for j, sem := range student.Semesters {
				if strings.EqualFold(sem.Name, semesterName) {
					for k, course := range sem.Courses {
						if course.CourseCode == courseCode {
							// Remove course at index k
							db.Students[i].Semesters[j].Courses = append(db.Students[i].Semesters[j].Courses[:k], db.Students[i].Semesters[j].Courses[k+1:]...)
							return storage.SaveDB(db)
						}
					}
					return errors.New("course not found")
				}
			}
			return errors.New("semester not found")
		}
	}
	return errors.New("student not found")
}

// UpdateCourse updates the title, credit units and score of an existing course
func UpdateCourse(studentID, semesterName, courseCode string, updatedCourse models.Course) error {
	db, err := storage.LoadDB()
	if err != nil {
		return err
	}

	for i, student := range db.Students {
		if student.ID == studentID {
			for j, sem := range student.Semesters {
				if strings.EqualFold(sem.Name, semesterName) {
					for k, course := range sem.Courses {
						if course.CourseCode == courseCode {
							// Update fields
							db.Students[i].Semesters[j].Courses[k].CourseCode = updatedCourse.CourseCode
							db.Students[i].Semesters[j].Courses[k].Title = updatedCourse.Title
							db.Students[i].Semesters[j].Courses[k].CreditUnits = updatedCourse.CreditUnits
							db.Students[i].Semesters[j].Courses[k].Score = updatedCourse.Score
							// Recalculate grade automatically
							db.Students[i].Semesters[j].Courses[k].Grade = calculateGrade(updatedCourse.Score)
							return storage.SaveDB(db)
						}
					}
					return errors.New("course not found")
				}
			}
			return errors.New("semester not found")
		}
	}
	return errors.New("student not found")
}
