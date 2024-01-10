package ui_handlers

import (
	"errors"
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/evanebb/gobble/profile"
	"github.com/evanebb/gobble/server/handlers"
	"github.com/google/uuid"
	"net/http"
)

func parseProfileFromPostForm(r *http.Request) (profile.Profile, error) {
	var p profile.Profile

	err := r.ParseForm()
	if err != nil {
		return p, err
	}

	requiredKeys := []string{"name", "description", "kernel", "initrd", "kernelParameters"}
	for _, v := range requiredKeys {
		if !r.PostForm.Has(v) {
			return p, errors.New("missing value " + v + " in POST form")
		}
	}

	kp, err := kernelparameters.ParseString(r.PostFormValue("kernelParameters"))
	if err != nil {
		return p, err
	}

	return profile.New(
		// the UUID needs to be set properly afterward by the caller, depending on whether we are creating a new one or updating an existing one
		uuid.Nil,
		r.PostFormValue("name"),
		r.PostFormValue("description"),
		r.PostFormValue("kernel"),
		r.PostFormValue("initrd"),
		kp,
	)
}

type UiProfileHandlerGroup struct {
	profileRepo profile.Repository
}

func NewUiProfileHandlerGroup(pr profile.Repository) UiProfileHandlerGroup {
	return UiProfileHandlerGroup{pr}
}

// Overview will list all profiles.
func (h UiProfileHandlerGroup) Overview(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.profileRepo.GetProfiles()
	if err != nil {
		renderError(w)
		return
	}

	d := templateData{Title: "Profiles", Data: profiles}
	renderTemplate(w, "profiles/overview", d)
}

// Show will show information about a single profile.
func (h UiProfileHandlerGroup) Show(w http.ResponseWriter, r *http.Request) {
	profileId, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		renderError(w)
		return
	}

	p, err := h.profileRepo.GetProfileById(profileId)
	if err != nil {
		renderError(w)
		return
	}

	d := templateData{Title: "Profile Information", Data: p}
	renderTemplate(w, "profiles/show", d)
}

// Create shows the page for creating a new profile.
func (h UiProfileHandlerGroup) Create(w http.ResponseWriter, r *http.Request) {
	d := templateData{Title: "Create profile"}
	renderTemplate(w, "profiles/create", d)
}

// Store will store a newly created profile.
func (h UiProfileHandlerGroup) Store(w http.ResponseWriter, r *http.Request) {
	p, err := parseProfileFromPostForm(r)
	if err != nil {
		renderError(w)
		return
	}

	p.Id = uuid.New()

	err = h.profileRepo.SetProfile(p)
	if err != nil {
		renderError(w)
		return
	}

	http.Redirect(w, r, "/ui/profiles/"+p.Id.String(), http.StatusSeeOther)
}

// Edit shows the page for editing an existing profile.
func (h UiProfileHandlerGroup) Edit(w http.ResponseWriter, r *http.Request) {
	profileId, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		renderError(w)
		return
	}

	p, err := h.profileRepo.GetProfileById(profileId)
	if err != nil {
		renderError(w)
		return
	}

	d := templateData{Title: "Edit Profile", Data: p}
	renderTemplate(w, "profiles/edit", d)
}

// Update will update the specified profile.
func (h UiProfileHandlerGroup) Update(w http.ResponseWriter, r *http.Request) {
	profileId, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		renderError(w)
		return
	}

	p, err := parseProfileFromPostForm(r)
	if err != nil {
		renderError(w)
		return
	}

	p.Id = profileId

	err = h.profileRepo.SetProfile(p)
	if err != nil {
		renderError(w)
		return
	}

	http.Redirect(w, r, "/ui/profiles/"+p.Id.String(), http.StatusSeeOther)
}

// Delete will delete the specified profile.
func (h UiProfileHandlerGroup) Delete(w http.ResponseWriter, r *http.Request) {
	profileId, err := handlers.GetUUIDFromRequest(r)
	if err != nil {
		renderError(w)
		return
	}

	err = h.profileRepo.DeleteProfileById(profileId)
	if err != nil {
		renderError(w)
		return
	}

	http.Redirect(w, r, "/ui/profiles", http.StatusSeeOther)
}
