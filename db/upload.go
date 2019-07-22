package db

import (
	"context"

	"github.com/deepakjacob/restyle/domain"
)

// Upload gets the user with the provided email.
func (fs *FireStore) Upload(ctx context.Context, user *domain.User,
	attrs *domain.ImgAttrs) error {
	doc := make(map[string]interface{})
	doc["obj_type"] = attrs.ObjType
	doc["material"] = attrs.Material            // Silk / Cotton
	doc["speciality"] = attrs.Speciality        // Kancheepuram / Banaras
	doc["dress_category"] = attrs.DressCategory // Women/ Girls/ Men/ Boys
	doc["age_min"] = attrs.AgeMin
	doc["age_max"] = attrs.AgeMax
	doc["tags"] = attrs.Tags
	doc["name"] = attrs.Name
	doc["location"] = attrs.Location
	doc["branch_item_available"] = attrs.Branches
	doc["date"] = attrs.DateUpload
	doc["date_from_visible"] = attrs.DateFromVisible
	doc["date_to_visible"] = attrs.DateToVisible
	doc["upload_count"] = attrs.UploadCount
	doc["user_id"] = user.UserID
	ref := fs.Collection("images").Doc(attrs.ObjType).Collection(attrs.Material).NewDoc()
	_, err := ref.Set(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}
