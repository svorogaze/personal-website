package api

type Author struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	AvatarLink string `json:"avatarLink" db:"avatar_link"`
	Login      string
	Password   string
}

type AuthorCard struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	AvatarLink string `json:"avatarLink" db:"avatar_link"`
}

func (api *API) getAuthorCardById(id int64) (AuthorCard, error) {
	row := api.db.QueryRowx("SELECT id, name, avatar_link FROM author WHERE id=$1", id)
	var a AuthorCard
	if err := row.StructScan(&a); err != nil {
		return a, err
	}
	return a, nil
}
