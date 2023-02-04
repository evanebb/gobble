package postgres

import (
	"context"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/profile"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) (ProfileRepository, error) {
	return ProfileRepository{db: db}, nil
}

type postgresProfile struct {
	Id               uint
	Name             string
	Description      string
	Distro           uint
	KernelParameters []string
}

func (r ProfileRepository) GetProfiles() ([]profile.Profile, error) {
	//TODO implement me
	panic("implement me")
}

func (r ProfileRepository) GetProfileById(id uint) (profile.Profile, error) {
	var pr profile.Profile
	var pp postgresProfile

	stmt := "SELECT id, name, description, distro, kernelParameters FROM profile WHERE id = $1"
	err := r.db.QueryRow(context.Background(), stmt, id).Scan(&pp.Id, &pp.Name, &pp.Description, &pp.Distro, &pp.KernelParameters)
	if err != nil {
		return pr, err
	}

	// If this errors someone directly inserted garbage into the database :(
	kp, err := kernelparameters.New(pp.KernelParameters)
	if err != nil {
		return pr, err
	}

	return profile.New(pp.Id, pp.Name, pp.Description, pp.Distro, kp), nil
}

func (r ProfileRepository) SetProfile(p profile.Profile) error {
	//TODO implement me
	panic("implement me")
}

func (r ProfileRepository) DeleteProfileById(id uint) error {
	//TODO implement me
	panic("implement me")
}
