package postgres

import (
	"context"
	"errors"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/profile"
	"github.com/evanebb/gobble/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	UUID             uuid.UUID
	Name             string
	Description      string
	Kernel           string
	Initrd           string
	KernelParameters []string
}

func (r ProfileRepository) GetProfiles() ([]profile.Profile, error) {
	var profiles []profile.Profile

	stmt := "SELECT id, uuid, name, description, kernel, initrd, kernelParameters FROM profile"
	rows, err := r.db.Query(context.Background(), stmt)
	if err != nil {
		return profiles, err
	}

	for rows.Next() {
		var pr profile.Profile
		var pp postgresProfile

		err = rows.Scan(&pp.Id, &pp.UUID, &pp.Name, &pp.Description, &pp.Kernel, &pp.Initrd, &pp.KernelParameters)
		if err != nil {
			return profiles, err
		}

		kp, err := kernelparameters.ParseStringSlice(pp.KernelParameters)
		if err != nil {
			return profiles, err
		}

		pr, err = profile.New(pp.UUID, pp.Name, pp.Description, pp.Kernel, pp.Initrd, kp)
		if err != nil {
			return profiles, err
		}

		profiles = append(profiles, pr)
	}

	return profiles, nil
}

func (r ProfileRepository) GetProfileById(id uuid.UUID) (profile.Profile, error) {
	var pr profile.Profile
	var pp postgresProfile

	stmt := "SELECT id, uuid, name, description, kernel, initrd, kernelParameters FROM profile WHERE uuid = $1"
	err := r.db.QueryRow(context.Background(), stmt, id).Scan(&pp.Id, &pp.UUID, &pp.Name, &pp.Description, &pp.Kernel, &pp.Initrd, &pp.KernelParameters)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr, repository.ErrNotFound
		}
		return pr, err
	}

	// If this errors someone directly inserted garbage into the database :(
	kp, err := kernelparameters.ParseStringSlice(pp.KernelParameters)
	if err != nil {
		return pr, err
	}

	return profile.New(pp.UUID, pp.Name, pp.Description, pp.Kernel, pp.Initrd, kp)
}

func (r ProfileRepository) SetProfile(p profile.Profile) error {
	stmt := "INSERT INTO profile (uuid, name, description, kernel, initrd, kernelParameters) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (uuid) DO UPDATE set name = $2, description = $3, kernel = $4, initrd = $5, kernelParameters = $6"
	_, err := r.db.Exec(context.Background(), stmt, p.Id, p.Name, p.Description, p.Kernel, p.Initrd, p.KernelParameters.StringSlice())
	if err != nil {
		return err
	}

	return nil
}

func (r ProfileRepository) DeleteProfileById(id uuid.UUID) error {
	stmt := "DELETE FROM profile WHERE uuid = $1"
	_, err := r.db.Exec(context.Background(), stmt, id)
	if err != nil {
		return err
	}

	return nil
}
