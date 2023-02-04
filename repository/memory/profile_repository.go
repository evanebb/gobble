package memory

import (
	"github.com/evanebb/gobble/profile"
	"github.com/google/uuid"
)

type ProfileRepository struct {
}

func NewProfileRepository() (ProfileRepository, error) {
	return ProfileRepository{}, nil
}

func (p ProfileRepository) GetProfiles() ([]profile.Profile, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProfileRepository) GetProfileById(id uuid.UUID) (profile.Profile, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProfileRepository) SetProfile(pr profile.Profile) error {
	//TODO implement me
	panic("implement me")
}

func (p ProfileRepository) DeleteProfileById(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
