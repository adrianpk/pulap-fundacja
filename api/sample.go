package api

import (
	"encoding/json"

	"github.com/gorilla/mux"

	"github.com/adrianpk/fundacja/app"
	"github.com/adrianpk/fundacja/logger"
	"github.com/adrianpk/fundacja/models"
  "github.com/adrianpk/fundacja/repo"
	"net/http"
	"net/url"
	"path"

	_ "github.com/lib/pq" // Import pq without side effects
)

type (
	// SampleResource - Resource
	SampleResource struct {
		Data models.Sample `json:"data"`
	}

	// SamplesResource - Resource
	SamplesResource struct {
		Data []models.Sample `json:"data"`
	}
)

// GetSamples - Returns a collection containing all samples.
// Handler for HTTP Get - "/samples"
func GetSamples(w http.ResponseWriter, r *http.Request) {
	// Get repo
	sampleRepo, err := repo.MakeSampleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityNotFound, err, http.StatusInternalServerError)
		return
	}
	// Select
	samples, err := sampleRepo.GetAll()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(SamplesResource{Data: samples})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// CreateSample - Creates a new Sample.
// Handler for HTTP Post - "/samples/create"
func CreateSample(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res SampleResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	sample := &res.Data
	// Get repo
	sampleRepo, err := repo.MakeSampleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Persist
	sampleRepo.Create(sample)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(SampleResource{Data: *sample})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GetSample - Returns a single Sample by its id or name.
// Handler for HTTP Get - "/samples/{ sample }"
func GetSample(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	key := vars["sample"]
	if len(key) == 36 {
		GetSampleByID(w, r)
	} else {
		GetSampleByName(w, r)
	}
}

// GetSampleByID - Returns a single Sample by its id or name.
// Handler for HTTP Get - "/samples/{ sample }"
func GetSampleByID(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["sample"]
	// Get repo
	sampleRepo, err := repo.MakeSampleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	sample, err := sampleRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(sample)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repsond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetSampleByName - Returns a single Sample by its name.
// Handler for HTTP Get - "/samples/{ sample }"
func GetSampleByName(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	sampleName := vars["sample"]
	// Get repo
	sampleRepo, err := repo.MakeSampleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	sample, err := sampleRepo.GetByName(sampleName)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(sample)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdateSample - Update an existing Sample.
// Handler for HTTP Put - "/samples/:id"
func UpdateSample(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["sample"]
	// Decode
	var res SampleResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	sample := &res.Data
	logger.Debugf("%s", sample)
	sample.ID = models.ToNullsString(id)
	logger.Debugf("%s", sample)
	// Get repo
	sampleRepo, err := repo.MakeSampleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Check against current sample
	currentSample, err := sampleRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Avoid ID spoofing
	err = verifyID(sample.IdentifiableModel, currentSample.IdentifiableModel)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Update
	err = sampleRepo.Update(sample)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(SampleResource{Data: *sample})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(j)
}

// DeleteSample - Delete an existing Sample
// Handler for HTTP Delete - "/samples/{ sample }"
func DeleteSample(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["sample"]
	// Get repo
	sampleRepo, err := repo.MakeSampleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = sampleRepo.Delete(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityDelete, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}

func sampleIDfromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	id := path.Base(dir)
	logger.Debugf("Sample id in url is %s", id)
	return id
}

func sampleNameFromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	sampleName := path.Base(dir)
	logger.Debugf("SampleName in url is %s", sampleName)
	return sampleName
}
