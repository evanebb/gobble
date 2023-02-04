package profile

type Repository interface {
	GetProfiles() ([]Profile, error)
	GetProfileById(id uint) (Profile, error)
	SetProfile(p Profile) error
	DeleteProfileById(id uint) error
}
