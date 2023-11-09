package entity

type Subject int8

const (
	IT Subject = iota
	Programming
	OPD // основы профессиональной деятельности
)

func (s Subject) Name() string {
	return map[Subject]string{
		IT:          "Информатика",
		Programming: "Программирование",
		OPD:         "ОПД",
	}[s]
}

func (s Subject) NameGenitiveCase() string {
	return map[Subject]string{
		IT:          "информатике",
		Programming: "программированию",
		OPD:         "ОПД",
	}[s]
}

type User struct {
	ID         int64
	FirstName  *string
	LastName   *string
	Patronymic *string

	OPDThreadID         *int // 3.6 -> 36
	ProgrammingThreadID *int // 2.5 -> 25
	ITThreadID          *int

	Action         *string
	TeacherSubject *Subject
}

// QueueUser is a person in subject queue
type QueueUser struct {
	UserID     int64
	FirstName  string
	LastName   string
	Patronymic *string

	LabName string
	Subject Subject
	Checked bool
	Passed  bool
	Retake  bool
}

type Thread struct {
	ID      int
	Subject Subject `gorm:"not null; uniqueIndex:subject_thread_uindex"`
	Name    string  `gorm:"not null; uniqueIndex:subject_thread_uindex"`
}

type Laboratory struct {
	ID      int     `gorm:"primaryKey; index; unique"`
	Subject Subject `gorm:"not null; uniqueIndex:subject_name_uindex"`
	Name    string  `gorm:"not null; uniqueIndex:subject_name_uindex"`
}
