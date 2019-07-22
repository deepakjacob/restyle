package handlers

import (
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

//Upload for uploading images and data
type Upload struct {
	UploadService service.UploadService
	// TODO for personalization
	// CustomizationService CustomizationService
}

// Handle handles the upload of the image file along with the mandatory params
func (u *Upload) Handle(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("img_file")
	if err != nil {
		logger.Log.Error("upload:handle::unable to read form file", zap.Error(err))
		http.Error(w,
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	attrs, err := parseImgAttrs(r.Form)
	if err != nil {
		logger.Log.Error("upload:handle:: error missing form fields", zap.Error(err))
		http.Error(w,
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = u.UploadService.Upload(r.Context(), attrs, file)
	if err != nil {
		logger.Log.Error("upload:handle:: error saving image and/or attributes", zap.Error(err))
		http.Error(w,
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func parseImgAttrs(vals url.Values) (*domain.ImgAttrs, error) {
	objType := vals.Get("obj_type")
	material := vals.Get("material")
	speciality := vals.Get("speciality")
	dressCategory := vals.Get("dress_category")
	ageMin := vals.Get("age_min")
	ageMax := vals.Get("age_max")
	tagStr := vals.Get("tags")
	tags := strings.Split(tagStr, ",")
	name := vals.Get("name")

	valsMissing :=
		isEmpty(objType) ||
			isEmpty(material) ||
			isEmpty(speciality) ||
			isEmpty(dressCategory) ||
			isEmpty(ageMin) ||
			isEmpty(ageMax) ||
			isEmpty(name) ||
			isEmpty(tagStr) ||
			tags == nil ||
			len(tags) == 0

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
		logger.Log.Error("upload missing mandatory params",
			zap.String("ObjType", objType),
			zap.String("Material", material),
			zap.String("Speciality", speciality),
			zap.String("Dress Category", dressCategory),
			zap.String("AgeMin", ageMin),
			zap.String("AgeMax", ageMax),
			zap.String("Name", name),
		)
		return nil, errors.New("missing mandatory values")
	}
	attrs := &domain.ImgAttrs{
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

func isEmpty(s string) bool {
	if len(s) == 0 {
		return true
	}
	return false
}
