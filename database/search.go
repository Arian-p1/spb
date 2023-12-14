package database

// import (
// 	"github.com/meilisearch/meilisearch-go"
// )

// func search() {
// 	client := meilisearch.NewClient(meilisearch.ClientConfig{
// 		Host:   "http://localhost:7700",
// 		APIKey: "dsjfadfh182er129p0",
// 	})
// }

// TODO: for search it should work with meilisearch
func Search(i interface{}, q string) error {
	return DB.Where("name = ?", "%"+q+"%").Find(&i).Error
}
