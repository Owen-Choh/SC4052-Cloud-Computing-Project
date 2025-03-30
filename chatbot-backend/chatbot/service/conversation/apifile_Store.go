package conversation

import (
	"database/sql"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
)

type APIFileStore struct {
	db *sql.DB
}

func NewAPIFileStore(db *sql.DB) types.APIFileStoreInterface {
	return &APIFileStore{db: db}
}

func (s *APIFileStore) GetAPIFileByFilepath(filepath string) (*types.APIFile, error) {
	rows, err := s.db.Query("SELECT * FROM apifiles WHERE filepath=?", filepath)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	theFile, err := scanRowsIntoAPIFile(rows)
	if err != nil {
		return nil, err
	}
	return theFile, nil
}

func (s *APIFileStore) GetAPIFileByID(apiFileID int) (*types.APIFile, error) {
	row := s.db.QueryRow("SELECT * FROM apifiles WHERE fileid=?", apiFileID)
	apiFile, err := scanRowIntoAPIFile(row)
	if err != nil {
		return nil, err
	}
	return apiFile, nil
}

func (s *APIFileStore) GetAPIFilesByUserID(userID int) ([]types.APIFile, error) {
	rows, err := s.db.Query("SELECT * FROM apifiles WHERE userid=? ORDER BY fileid", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apiFiles := []types.APIFile{}
	for rows.Next() {
		apiFile, err := scanRowsIntoAPIFile(rows)
		if err != nil {
			return nil, err
		}
		apiFiles = append(apiFiles, *apiFile)
	}

	return apiFiles, nil
}

func (s *APIFileStore) CreateAPIFile(apiFilePayload types.NewAPIFile) (int, error) {
	res, dberr := s.db.Exec(
		"INSERT INTO apifiles (chatbotid, createddate, filepath, fileuri) VALUES (?, ?, ?, ?)",
		apiFilePayload.Chatbotid,
		apiFilePayload.Createddate,
		apiFilePayload.Filepath,
		apiFilePayload.Fileuri,
	)
	if dberr != nil {
		return 0, dberr
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *APIFileStore) UpdateAPIFile(apiFilePayload types.UpdateAPIFile) error {
	_, dberr := s.db.Exec(
		"UPDATE apifiles SET chatbotid=?, createddate=?, filepath=?, fileuri=? WHERE fileid=?",
		apiFilePayload.Chatbotid,
		apiFilePayload.Createddate,
		apiFilePayload.Filepath,
		apiFilePayload.Fileuri,
		apiFilePayload.Fileid,
	)
	return dberr
}

func (s *APIFileStore) DeleteAPIFile(apiFileID int) error {
	_, dberr := s.db.Exec("DELETE FROM apifiles WHERE fileid=?", apiFileID)
	return dberr
}

func scanRowsIntoAPIFile(rows *sql.Rows) (*types.APIFile, error) {
	apifile := new(types.APIFile)

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

func scanRowIntoAPIFile(row *sql.Row) (*types.APIFile, error) {
	apiFile := new(types.APIFile)
	err := row.Scan(
		&apiFile.Fileid,
		&apiFile.Chatbotid,
		&apiFile.Createddate,
		&apiFile.Filepath,
		&apiFile.Fileuri,
	)
	if err != nil {
		return nil, err
	}
	return apiFile, nil
}
