package dbops

import (
	"database/sql"
	"log"
	"myproject/video_server/api/defs"
	"time"

	"github.com/pborman/uuid"

	_ "github.com/go-sql-driver/mysql"
)

func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO users (login_name, pwd) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE login_name = ?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	var pwd string
	// row := stmtOut.QueryRow(loginName)
	// err = row.Scan(&pwd)
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	defer stmtOut.Close()

	return pwd, nil
}

func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM users WHERE login_name=? and pwd=?")
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

// Videos
func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	vid := uuid.New()
	t := time.Now()
	ctime := t.Format("jan 02 2006, 15:04:05")
	stmtIns, err := dbConn.Prepare(`INSERT INTO video_info (id, author_id, name, display_ctime) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}
	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}

	res := &defs.VideoInfo{ID: vid, AuthorID: aid, Name: name, DisplayCTime: ctime}
	defer stmtIns.Close()
	return res, nil
}

func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`select author_id, name, display_ctime from video_info WHERE id=? `)

	var (
		aid  int
		dct  string
		name string
	)

	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	defer stmtOut.Close()

	res := &defs.VideoInfo{ID: vid, AuthorID: aid, Name: name, DisplayCTime: dct}
	return res, nil
}

func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id=?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}

	defer stmtDel.Close()
	return nil
}

// comments
func AddNewComments(vid string, aid int, content string) error {
	id := uuid.New()

	stmtIns, err := dbConn.Prepare("insert into comments (id, video_id, author_id, content) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err2 := stmtIns.Exec(id, vid, aid, content)
	if err2 != nil {
		return err2
	}
	defer stmtIns.Close()
	return nil
}

func ListComments(vid string, from, to int) ([]*defs.Comments, error) {
	var res []*defs.Comments

	stmtOut, err := dbConn.Prepare(`select comments.id, users.login_name, comments.content from comments inner join users on comments.author_id = users.id where comments.video_id = ? and comments.time > from_unixtime(?) and comments.time <= from_unixtime(?)`)
	if err != nil {
		return res, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return res, err
		}

		c := &defs.Comments{ID: id, VideoID: vid, Author: name, Content: content}
		res = append(res, c)

	}
	return res, nil
}
