package postgres

import (
	"context"
	"github.com/evanebb/gobble/distro"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DistroRepository struct {
	db *pgxpool.Pool
}

func NewDistroRepository(db *pgxpool.Pool) (DistroRepository, error) {
	return DistroRepository{db: db}, nil
}

type postgresDistro struct {
	Id               uint
	Name             string
	Description      string
	Kernel           string
	Initrd           string
	KernelParameters []string
}

func (r DistroRepository) GetDistros() ([]distro.Distro, error) {
	//TODO implement me
	panic("implement me")
}

func (r DistroRepository) GetDistroById(id uint) (distro.Distro, error) {
	var d distro.Distro
	var pd postgresDistro

	stmt := "SELECT id, name, description, kernel, initrd, kernelParameters FROM distro WHERE id = $1"
	err := r.db.QueryRow(context.Background(), stmt, id).Scan(&pd.Id, &pd.Name, &pd.Description, &pd.Kernel, &pd.Initrd, &pd.KernelParameters)
	if err != nil {
		return d, err
	}

	kp, err := kernelparameters.New(pd.KernelParameters)
	if err != nil {
		return d, err
	}

	return distro.New(pd.Id, pd.Name, pd.Description, pd.Kernel, pd.Initrd, kp), nil
}

func (r DistroRepository) SetDistro(d distro.Distro) error {
	//TODO implement me
	panic("implement me")
}

func (r DistroRepository) DeleteDistroById(id uint) error {
	//TODO implement me
	panic("implement me")
}
