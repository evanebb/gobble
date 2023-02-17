package postgres

import (
	"context"
	"errors"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/repository"
	"github.com/evanebb/gobble/system"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	UUID             uuid.UUID
	Name             string
	Description      string
	Profile          uuid.UUID
	Mac              net.HardwareAddr
	KernelParameters []string
}

func (r SystemRepository) GetSystems() ([]system.System, error) {
	var systems []system.System

	stmt := "SELECT id, uuid, name, description, profile, mac, kernelParameters FROM system"
	rows, err := r.db.Query(context.Background(), stmt)
	if err != nil {
		return systems, err
	}

	for rows.Next() {
		var sys system.System
		var ps postgresSystem

		err = rows.Scan(&ps.Id, &ps.UUID, &ps.Name, &ps.Description, &ps.Profile, &ps.Mac, &ps.KernelParameters)
		if err != nil {
			return systems, err
		}

		kp, err := kernelparameters.New(ps.KernelParameters)
		if err != nil {
			return systems, err
		}

		sys, err = system.New(ps.UUID, ps.Name, ps.Description, ps.Profile, ps.Mac, kp)
		if err != nil {
			return systems, err
		}

		systems = append(systems, sys)
	}

	return systems, nil
}

func (r SystemRepository) GetSystemByMacAddress(mac net.HardwareAddr) (system.System, error) {
	var sys system.System
	var ps postgresSystem

	stmt := "SELECT id, uuid, name, description, profile, mac, kernelParameters FROM system WHERE mac = $1"
	err := r.db.QueryRow(context.Background(), stmt, mac).Scan(&ps.Id, &ps.UUID, &ps.Name, &ps.Description, &ps.Profile, &ps.Mac, &ps.KernelParameters)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sys, repository.ErrNotFound
		}
		return sys, err
	}

	// If this errors someone directly inserted garbage into the database :(
	kp, err := kernelparameters.New(ps.KernelParameters)
	if err != nil {
		return sys, err
	}

	return system.New(ps.UUID, ps.Name, ps.Description, ps.Profile, ps.Mac, kp)
}

func (r SystemRepository) GetSystemById(id uuid.UUID) (system.System, error) {
	var sys system.System
	var ps postgresSystem

	stmt := "SELECT id, uuid, name, description, profile, mac, kernelParameters FROM system WHERE uuid = $1"
	err := r.db.QueryRow(context.Background(), stmt, id).Scan(&ps.Id, &ps.UUID, &ps.Name, &ps.Description, &ps.Profile, &ps.Mac, &ps.KernelParameters)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sys, repository.ErrNotFound
		}
		return sys, err
	}

	// If this errors someone directly inserted garbage into the database :(
	kp, err := kernelparameters.New(ps.KernelParameters)
	if err != nil {
		return sys, err
	}

	return system.New(ps.UUID, ps.Name, ps.Description, ps.Profile, ps.Mac, kp)
}

func (r SystemRepository) SetSystem(s system.System) error {
	stmt := "INSERT INTO system (uuid, name, description, profile, mac, kernelParameters) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (uuid) DO UPDATE set name = $2, description = $3, profile = $4, mac = $5, kernelParameters = $6"
	_, err := r.db.Exec(context.Background(), stmt, s.Id, s.Name, s.Description, s.Profile, s.Mac, kernelparameters.FormatKernelParameters(s.KernelParameters))
	if err != nil {
		return err
	}

	return nil
}

func (r SystemRepository) DeleteSystemById(id uuid.UUID) error {
	stmt := "DELETE FROM system WHERE uuid = $1"
	_, err := r.db.Exec(context.Background(), stmt, id)
	if err != nil {
		return err
	}

	return nil
}
