package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/utils"
)

type MySqlClient struct {
	instance     *sql.DB
	lock         *sync.Mutex
	maxRecordAge time.Duration
}

type UrlEntry struct {
	CreatedAt int    `json:"createdAt"`
	Id        string `json:"id"`
	Link      string `json:"url"`
}

func InitMySqlConnection(
	maxRecordAge time.Duration,
) *MySqlClient {
	connStr := ""
	if consts.Local {
		fmt.Println("MySQL config: local")
		connStr = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s",
			os.Getenv("MYSQL_UN_DEV"),
			os.Getenv("MYSQL_PW_DEV"),
			os.Getenv("MYSQL_URL_DEV"),
			os.Getenv("MYSQL_DB_DEV"),
		)
	} else {
		fmt.Println("MySQL config: remote")
		connStr = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s",
			os.Getenv("MYSQL_UN"),
			os.Getenv("MYSQL_PW"),
			os.Getenv("MYSQL_URL"),
			os.Getenv("MYSQL_DB"),
		)
	}

	mySqlInstance, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}

	client := &MySqlClient{
		instance:     mySqlInstance,
		lock:         &sync.Mutex{},
		maxRecordAge: maxRecordAge,
	}

	return client
}

func GetMySqlClient(ctx context.Context) (*MySqlClient, error) {
	mySqlClient, ok := ctx.Value(consts.MySqlClientKey).(*MySqlClient)
	if !ok {
		return nil, fmt.Errorf("couldn't find %s in request context", consts.MySqlClientKey)
	}
	return mySqlClient, nil
}

/* -------------------- Methods --------------------*/

func (mySqlClient *MySqlClient) CheckExists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := mySqlClient.instance.
		QueryRow("SELECT EXISTS(SELECT * FROM Urls WHERE Id = ?);", id).
		Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (mySqlClient *MySqlClient) GetShortestNewKey(
	ctx context.Context,
	key string,
) (string, error) {
	shortenedKey := ""

	for i := 3; i <= len(key)+1; i++ {
		shortenedKey = key[:i]
		exists, err := mySqlClient.CheckExists(ctx, shortenedKey)
		if err != nil {
			return "", err
		}
		if !exists {
			return shortenedKey, nil
		}
	}

	return mySqlClient.GetShortestNewKey(ctx, utils.RandString(8))
}

func (mySqlClient *MySqlClient) GetUrlRecordById(ctx context.Context, id string) (*UrlEntry, error) {
	var entry UrlEntry
	err := mySqlClient.instance.
		QueryRow("SELECT * FROM Urls where Id = ?", id).
		Scan(&entry.CreatedAt, &entry.Id, &entry.Link)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &entry, nil
}

func (mySqlClient *MySqlClient) WriteUrlRecord(ctx context.Context, id string, url string) (*UrlEntry, error) {
	createdAt := time.Now().Unix()
	query := fmt.Sprintf("INSERT INTO Urls(CreatedAt, Id, Link) VALUES(%d, '%s', '%s');", createdAt, id, url)
	insert, err := mySqlClient.instance.Query(query)
	if err != nil {
		return nil, err
	}
	defer insert.Close()

	return &UrlEntry{CreatedAt: int(createdAt), Id: id, Link: url}, nil
}

func (mySqlClient *MySqlClient) SetClearInterval(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		func() {
			mySqlClient.lock.Lock()
			createdAt := time.Now().Unix() - int64(mySqlClient.maxRecordAge.Seconds())
			query := fmt.Sprintf("DELETE FROM Urls WHERE CreatedAt < %d", createdAt)
			queryRun, err := mySqlClient.instance.Query(query)
			if err != nil {
				fmt.Println(err)
			}
			queryRun.Close()
			mySqlClient.lock.Unlock()
		}()
	}
}

func (mySqlClient *MySqlClient) LogAll() {
	results, err := mySqlClient.instance.Query("SELECT * FROM Urls")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var urlEntry UrlEntry
		err = results.Scan(&urlEntry.CreatedAt, &urlEntry.Id, &urlEntry.Link)
		if err != nil {
			panic(err.Error())
		}
		log.Println(urlEntry.Link)
	}
}
