package postgres

import (
	"context"
	"github.com/evanebb/gobble/distro"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
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
	UUID             uuid.UUID
	Name             string
	Description      string
	Kernel           string
	Initrd           string
	KernelParameters []string
}

func (r DistroRepository) GetDistros() ([]distro.Distro, error) {
	var distros []distro.Distro

	stmt := "SELECT id, uuid, name, description, kernel, initrd, kernelParameters FROM distro"
	rows, err := r.db.Query(context.Background(), stmt)
	if err != nil {
		return distros, err
	}

	for rows.Next() {
		var dis distro.Distro
		var pd postgresDistro

		err = rows.Scan(&pd.Id, &pd.UUID, &pd.Name, &pd.Description, &pd.Kernel, &pd.Initrd, &pd.KernelParameters)
		if err != nil {
			return distros, err
		}

		kp, err := kernelparameters.New(pd.KernelParameters)
		if err != nil {
			return distros, err
		}

		dis = distro.New(pd.UUID, pd.Name, pd.Description, pd.Kernel, pd.Initrd, kp)
		distros = append(distros, dis)
	}

	return distros, nil
}

func (r DistroRepository) GetDistroById(id uuid.UUID) (distro.Distro, error) {
	var d distro.Distro
	var pd postgresDistro

	stmt := "SELECT id, uuid, name, description, kernel, initrd, kernelParameters FROM distro WHERE uuid = $1"
	err := r.db.QueryRow(context.Background(), stmt, id).Scan(&pd.Id, &pd.UUID, &pd.Name, &pd.Description, &pd.Kernel, &pd.Initrd, &pd.KernelParameters)
	if err != nil {
		return d, err
	}

	kp, err := kernelparameters.New(pd.KernelParameters)
	if err != nil {
		return d, err
	}

	return distro.New(pd.UUID, pd.Name, pd.Description, pd.Kernel, pd.Initrd, kp), nil
}

func (r DistroRepository) SetDistro(d distro.Distro) error {
	stmt := "INSERT INTO distro (uuid, name, description, kernel, initrd, kernelParameters) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (uuid) DO UPDATE set name = $2, description = $3, kernel = $4, initrd = $5, kernelParameters = $6"
	_, err := r.db.Exec(context.Background(), stmt, d.Id(), d.Name(), d.Description(), d.Kernel(), d.Initrd(), kernelparameters.FormatKernelParameters(d.KernelParameters()))
	if err != nil {
		return err
	}

	return nil
}

func (r DistroRepository) DeleteDistroById(id uuid.UUID) error {
	stmt := "DELETE FROM distro WHERE uuid = $1"
	_, err := r.db.Exec(context.Background(), stmt, id)
	if err != nil {
		return err
	}

	return nil
}
