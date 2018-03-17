package repository

import (
	"database/sql"
	"fmt"

	"github.com/bxcodec/go-clean-arch/author"
	models "github.com/ivzb/semaphore_server/user"
)

type postgreUserRepository struct {
	Conn *sql.DB
}

func NewPostgreUserRepository(Conn *sql.DB) UserRepository {

	return &postgreUserRepository{Conn}
}

func (m *postgreUserRepository) fetch(query string, args ...interface{}) ([]*models.User, error) {

	rows, err := m.Conn.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*models.User, 0)
	for rows.Next() {
		t := new(models.User)
		authorID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&authorID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		t.Author = author.Author{
			ID: authorID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *postgreUserRepository) Fetch(cursor string, num int64) ([]*models.User, error) {

	query := `SELECT id,title,content, author_id, updated_at, created_at
																																																																													FROM article WHERE ID > ? LIMIT ?`

	return m.fetch(query, cursor, num)

}
func (m *postgreUserRepository) GetByID(id int64) (*models.User, error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
																																																																																					FROM article WHERE ID = ?`

	list, err := m.fetch(query, id)
	if err != nil {
		return nil, err
	}

	a := &models.User{}
	if len(list) > 0 {
		a = list[0]
	} else {
		return nil, models.NOT_FOUND_ERROR
	}

	return a, nil
}

func (m *postgreUserRepository) GetByTitle(title string) (*models.User, error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
																																																																																																													FROM article WHERE title = ?`

	list, err := m.fetch(query, title)
	if err != nil {
		return nil, err
	}

	a := &models.User{}
	if len(list) > 0 {
		a = list[0]
	} else {
		return nil, models.NOT_FOUND_ERROR
	}
	return a, nil
}

func (m *postgreUserRepository) Store(a *models.User) (int64, error) {

	query := `INSERT  article SET title=? , content=? , author_id=?, updated_at=? , created_at=?`
	stmt, err := m.Conn.Prepare(query)
	if err != nil {

		return 0, err
	}

	res, err := stmt.Exec(a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.CreatedAt)
	if err != nil {

		return 0, err
	}
	return res.LastInsertId()
}

func (m *postgreUserRepository) Delete(id int64) (bool, error) {
	query := "DELETE FROM article WHERE id = ?"

	stmt, err := m.Conn.Prepare(query)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(id)
	if err != nil {

		return false, err
	}
	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAfected)
		return false, err
	}

	return true, nil
}
func (m *postgreUserRepository) Update(ar *models.User) (*models.User, error) {
	query := `UPDATE article set title=?, content=?, author_id=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.Prepare(query)
	if err != nil {
		return nil, nil
	}

	res, err := stmt.Exec(ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.ID)
	if err != nil {
		return nil, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)
		return nil, err
	}

	return ar, nil
}
