package postgres

import (
	"context"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/system"
	"github.com/jackc/pgx/v5/pgxpool"
	"net"
)

type SystemRepository struct {
	db *pgxpool.Pool
}

func NewSystemRepository(db *pgxpool.Pool) (SystemRepository, error) {
	return SystemRepository{db: db}, nil
}

type postgresSystem struct {
	Id               uint
	Name             string
	Description      string
	Profile          uint
	Mac              net.HardwareAddr
	KernelParameters []string
}

func (s SystemRepository) GetSystems() ([]system.System, error) {
	var systems []system.System

	stmt := "SELECT id, name, description, profile, mac, kernelParameters FROM system"
	rows, err := s.db.Query(context.Background(), stmt)
	if err != nil {
		return systems, err
	}

	for rows.Next() {
		var sys system.System
		var ps postgresSystem

		err = rows.Scan(&ps.Id, &ps.Name, &ps.Description, &ps.Profile, &ps.Mac, &ps.KernelParameters)
		if err != nil {
			return systems, err
		}

		kp, err := kernelparameters.New(ps.KernelParameters)
		if err != nil {
			return systems, err
		}

		sys = system.New(ps.Id, ps.Name, ps.Description, ps.Profile, ps.Mac, kp)

		systems = append(systems, sys)
	}

	return systems, nil
}

func (s SystemRepository) GetSystemByMacAddress(mac net.HardwareAddr) (system.System, error) {
	var sys system.System
	var ps postgresSystem

	// TODO: type conversion for MAC, 99% sure that it will throw an error
	stmt := "SELECT id, name, description, profile, mac, kernelParameters FROM system WHERE mac = $1"
	err := s.db.QueryRow(context.Background(), stmt, mac).Scan(&ps.Id, &ps.Name, &ps.Description, &ps.Profile, &ps.Mac, &ps.KernelParameters)
	if err != nil {
		return sys, err
	}

	// If this errors someone directly inserted garbage into the database :(
	kp, err := kernelparameters.New(ps.KernelParameters)
	if err != nil {
		return sys, err
	}

	return system.New(ps.Id, ps.Name, ps.Description, ps.Profile, ps.Mac, kp), nil
}

func (s SystemRepository) GetSystemById(id uint) (system.System, error) {
	var sys system.System
	var ps postgresSystem

	stmt := "SELECT id, name, description, profile, mac, kernelParameters FROM system WHERE id = $1"
	err := s.db.QueryRow(context.Background(), stmt, id).Scan(&ps.Id, &ps.Name, &ps.Description, &ps.Profile, &ps.Mac, &ps.KernelParameters)
	if err != nil {
		return sys, err
	}

	// If this errors someone directly inserted garbage into the database :(
	kp, err := kernelparameters.New(ps.KernelParameters)
	if err != nil {
		return sys, err
	}

	return system.New(ps.Id, ps.Name, ps.Description, ps.Profile, ps.Mac, kp), nil
}

func (s SystemRepository) SetSystem(sys system.System) error {
	// TODO: check if reusing arguments is possible
	stmt := "INSERT INTO system (id, name, description, profile, mac, kernelParameters) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO UPDATE set name = $2, description = $3, profile = $4, mac = $5, kernelParameters = $6"
	_, err := s.db.Exec(context.Background(), stmt, sys.Id(), sys.Name(), sys.Description(), sys.Profile(), sys.Mac(), sys.KernelParameters())
	if err != nil {
		return err
	}

	return nil
}

func (s SystemRepository) DeleteSystemById(id uint) error {
	stmt := "DELETE FROM system WHERE id = $1"
	_, err := s.db.Exec(context.Background(), stmt, id)
	if err != nil {
		return err
	}

	return nil
}
