package api

type Author struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	AvatarLink string `json:"avatarLink" db:"avatar_link"`
	Login      string
	Password   string
}

func (api *API) getAuthorById(id int64) (Author, error) {
	row := api.db.QueryRowx("SELECT * FROM author WHERE id=$1", id)
	var a Author
	if err := row.StructScan(&a); err != nil {
		return a, err
	}
	return a, nil
}
