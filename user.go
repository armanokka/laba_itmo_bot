package main

type User struct {
	Users
	error func(error)
}

// NewUser return User with such id
func NewUser(id int64, errfunc func(error)) User {
	return User{
		Users: Users{ID: id},
		error: errfunc,
	}
}

func (u User) Exists() bool {
	var exists bool
	if err := db.Raw("SELECT EXISTS(SELECT id FROM users WHERE id=?)", u.ID).Find(&exists).Error; err != nil {
		u.error(err)
	}
	return exists
}
// Create creates user in db. Also fills a user, e.g. Fill()
func (u User) Create(user Users) {
	if err := db.Create(&user).Error; err != nil {
		u.error(err)
	} else {
		u.Users = user
	}
}

func (u *User) Fill() {
	if err := db.Model(&Users{}).Where("id = ?", u.ID).Find(&u.Users).Error; err != nil {
		u.error(err)
	}
}

func (u *User) Update(user Users) {
	if err := db.Model(&u.Users).Updates(user).Error; err != nil {
		u.error(err)
	}
	u.Fill()
}

func (u User) Localize(text string, placeholders ...interface{}) string {
	return localize(text, u.Lang, placeholders...)
}