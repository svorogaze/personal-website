package api

type Blog struct {
	Id           int64  `json:"id"`
	Title        string `json:"title"`
	Contents     string `json:"contents"`
	Description  string `json:"description"`
	CreationDate string `json:"creationDate" db:"creation_date"`
	AuthorId     int64  `json:"authorId"`
	PictureLink  string `json:"pictureLink"`
}

type BlogCard struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PictureLink string `json:"pictureLink"`
}

func (api *API) getBlogById(id int64) (Blog, error) {
	row := api.db.QueryRowx("SELECT * FROM blog WHERE id=$1", id)
	var b Blog
	if err := row.StructScan(&b); err != nil {
		return b, err
	}
	return b, nil
}
