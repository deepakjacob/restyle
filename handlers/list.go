package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/service"
	"go.uber.org/zap"
)

//List for uploading images and data
type List struct {
	ListService service.ListService
	// TODO for personalization
	// CustomizationService CustomizationService
}

// Handle handles the upload of the image file along with the mandatory params
func (l *List) Handle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	attrs, err := parseSearchAttrs(r.Form)
	if err != nil {
		logger.Log.Error("list:handle:: error missing form fields", zap.Error(err))
		http.Error(w,
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	result, err := l.ListService.List(r.Context(), attrs)
	if err != nil {
		logger.Log.Error("list:handle:: error listing image and/or attributes", zap.Error(err))
		http.Error(w,
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(result)
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func parseSearchAttrs(vals url.Values) (*domain.ImgSearch, error) {
	objType := vals.Get("obj_type")
	material := vals.Get("material")
	speciality := vals.Get("speciality")
	dressCategory := vals.Get("dress_category")
	ageMin := vals.Get("age_min")
	ageMax := vals.Get("age_max")
	tagStr := vals.Get("tags")
	tags := strings.Split(tagStr, ",")
	name := vals.Get("name")

	valsMissing := isEmpty(objType)

	minAge, err := strconv.ParseInt(ageMin, 10, 8)
	if err != nil {
		logger.Log.Error("error converting to number",
			zap.String("AgeMin", ageMin))
		return nil, errors.New("error converting min age to number")
	}
	maxAge, err := strconv.ParseInt(ageMax, 10, 8)
	if err != nil {
		logger.Log.Error("error converting to number",
			zap.String("AgeMax", ageMax))
		return nil, errors.New("error converting max age to number")
	}

	if valsMissing {
		logger.Log.Error("upload missing mandatory params", zap.String("ObjType", objType))
		return nil, errors.New("missing mandatory values")
	}
	attrs := &domain.ImgSearch{
		ObjType:       objType,
		Material:      material,
		Speciality:    speciality,
		DressCategory: dressCategory,
		AgeMin:        int8(minAge),
		AgeMax:        int8(maxAge),
		Tags:          tags,
		Name:          name,
	}

	return attrs, nil
}
