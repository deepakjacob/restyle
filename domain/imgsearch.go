package domain

import "time"

// ImgSearch img search params
type ImgSearch struct {
	ObjType         string    `json:"obj_type"          firestore:"obj_type"`
	Material        string    `json:"material"          firestore:"material"`
	Speciality      string    `json:"speciality"        firestore:"speciality"`
	DressCategory   string    `json:"dress_category"    firestore:"dress_category"` // Women / Girls / Men/ Boys
	AgeMin          int8      `json:"age_min"           firestore:"age_min"`
	AgeMax          int8      `json:"age_max"           firestore:"age_max"`
	Tags            []string  `json:"tags"              firestore:"tags"`
	Name            string    `json:"batch_name"        firestore:"batch_name"`
	Location        string    `json:"location"          firestore:"location"`
	Branches        []string  `json:"branches"          firestore:"branches"`
	DateUpload      time.Time `json:"date_upload"       firestore:"date_upload"`
	DateFromVisible time.Time `json:"visible_from_date" firestore:"visible_from_date"`
	DateToVisible   time.Time `json:"visible_to_date"   firestore:"visible_to_date"`
	UploadCount     int       `json:"upload_count"      firestore:"load_count"`
}
