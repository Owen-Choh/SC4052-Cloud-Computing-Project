package conversation

import (
	"database/sql"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
)

type APIFileStore struct {
	db *sql.DB
}

func NewAPIFileStore(db *sql.DB) *APIFileStore {
	return &APIFileStore{db: db}
}

func (s *APIFileStore) GetAPIFileByFilepath(filepath string) (*types.APIfile, error) {
	rows, err := s.db.Query("SELECT * FROM apifiles WHERE filepath=?", filepath)
	if err != nil {
		return nil, err
	}
	theFile, err := scanRowsIntoApifile(rows)
	if err != nil {
		return nil, err
	}
	return theFile, nil
}

func (s *APIFileStore) StoreAPIFile(apiFilePayload string) (int, error) {
	return 0, nil
}

func scanRowsIntoApifile(rows *sql.Rows) (*types.APIfile, error) {
	apifile := new(types.APIfile)

	err := rows.Scan(
		&apifile.Fileid,
		&apifile.Chatbotid,
		&apifile.Createddate,
		&apifile.Filepath,
		&apifile.Fileuri,
	)

	if err != nil {
		return nil, err
	}
	return apifile, nil
}
