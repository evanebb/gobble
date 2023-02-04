package memory

import "github.com/evanebb/gobble/profile"

type ProfileRepository struct {
}

func NewProfileRepository() (ProfileRepository, error) {
	return ProfileRepository{}, nil
}

func (p ProfileRepository) GetProfiles() ([]profile.Profile, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProfileRepository) GetProfileById(id uint) (profile.Profile, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProfileRepository) SetProfile(pr profile.Profile) error {
	//TODO implement me
	panic("implement me")
}

func (p ProfileRepository) DeleteProfileById(id uint) error {
	//TODO implement me
	panic("implement me")
}
