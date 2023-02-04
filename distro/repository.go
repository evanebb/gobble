package distro

type Repository interface {
	GetDistros() ([]Distro, error)
	GetDistroById(id uint) (Distro, error)
	SetDistro(d Distro) error
	DeleteDistroById(id uint) error
}
