package ui_handlers

import (
	"errors"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/profile"
	"github.com/evanebb/gobble/server/handlers"
	"github.com/evanebb/gobble/system"
	"github.com/google/uuid"
	"net"
	"net/http"
	"strings"
)

func parseSystemFromPostForm(r *http.Request) (system.System, error) {
	var s system.System

	err := r.ParseForm()
	if err != nil {
		return s, err
	}

	requiredKeys := []string{"name", "description", "profile", "mac", "kernelParameters"}
	for _, v := range requiredKeys {
		if !r.PostForm.Has(v) {
			return s, errors.New("missing value " + v + " in POST form")
		}
	}

	kp, err := kernelparameters.New(strings.Split(r.PostFormValue("kernelParameters"), " "))
	if err != nil {
		return s, err
	}

	profileId, err := uuid.Parse(r.PostFormValue("profile"))
	if err != nil {
		return s, err
	}

	macAddress, err := net.ParseMAC(r.PostFormValue("mac"))
	if err != nil {
		return s, err
	}

	return system.New(
		// the UUID needs to be set properly afterward by the caller, depending on whether we are creating a new one or updating an existing one
		uuid.Nil,
		r.PostFormValue("name"),
		r.PostFormValue("description"),
		profileId,
		macAddress,
		kp,
	)
}

type UiSystemHandlerGroup struct {
	systemRepo  system.Repository
	profileRepo profile.Repository
}

func NewUiSystemHandlerGroup(sr system.Repository, pr profile.Repository) UiSystemHandlerGroup {
	return UiSystemHandlerGroup{sr, pr}
}

// Overview will list all systems.
func (h UiSystemHandlerGroup) Overview(w http.ResponseWriter, r *http.Request) {
	systems, err := h.systemRepo.GetSystems()
	if err != nil {
		renderError(w, err)
		return
	}

	d := templateData{Title: "Systems", Data: systems}
	renderTemplate(w, "systems/overview", d)
}

// Show will show information about a single system.
func (h UiSystemHandlerGroup) Show(w http.ResponseWriter, r *http.Request) {
	systemId, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		renderError(w, err)
		return
	}

	s, err := h.systemRepo.GetSystemById(systemId)
	if err != nil {
		renderError(w, err)
		return
	}

	p, err := h.profileRepo.GetProfileById(s.Profile)
	if err != nil {
		renderError(w, err)
		return
	}

	d := templateData{Title: "System information", Data: struct {
		System  system.System
		Profile profile.Profile
	}{
		System:  s,
		Profile: p,
	}}
	renderTemplate(w, "systems/show", d)
}

// Create shows the page for creating a new system.
func (h UiSystemHandlerGroup) Create(w http.ResponseWriter, r *http.Request) {
	d := templateData{Title: "Create system"}
	renderTemplate(w, "systems/create", d)
}

// Store will store a newly created system.
func (h UiSystemHandlerGroup) Store(w http.ResponseWriter, r *http.Request) {
	s, err := parseSystemFromPostForm(r)
	if err != nil {
		renderError(w, err)
		return
	}

	s.Id = uuid.New()

	err = h.systemRepo.SetSystem(s)
	if err != nil {
		renderError(w, err)
		return
	}

	http.Redirect(w, r, "/ui/systems/"+s.Id.String(), http.StatusSeeOther)
}

// Edit shows the page for editing an existing system.
func (h UiSystemHandlerGroup) Edit(w http.ResponseWriter, r *http.Request) {
	systemId, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		renderError(w, err)
		return
	}

	s, err := h.systemRepo.GetSystemById(systemId)
	if err != nil {
		renderError(w, err)
		return
	}

	profiles, err := h.profileRepo.GetProfiles()
	if err != nil {
		renderError(w, err)
		return
	}

	d := templateData{Title: "Edit system", Data: struct {
		System   system.System
		Profiles []profile.Profile
	}{
		System:   s,
		Profiles: profiles,
	}}
	renderTemplate(w, "systems/edit", d)
}

// Update will update the specified system.
func (h UiSystemHandlerGroup) Update(w http.ResponseWriter, r *http.Request) {
	systemId, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		renderError(w, err)
		return
	}

	s, err := parseSystemFromPostForm(r)
	if err != nil {
		renderError(w, err)
		return
	}

	s.Id = systemId

	err = h.systemRepo.SetSystem(s)
	if err != nil {
		renderError(w, err)
		return
	}

	http.Redirect(w, r, "/ui/systems/"+s.Id.String(), http.StatusSeeOther)
}

// Delete will delete the specified system.
func (h UiSystemHandlerGroup) Delete(w http.ResponseWriter, r *http.Request) {
	systemId, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		renderError(w, err)
		return
	}

	err = h.systemRepo.DeleteSystemById(systemId)
	if err != nil {
		renderError(w, err)
		return
	}

	http.Redirect(w, r, "/ui/systems", http.StatusSeeOther)
}
