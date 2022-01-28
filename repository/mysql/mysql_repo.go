package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/cfloress-gb-cl/final-project-bootcamp/repository/user"
	_ "github.com/go-sql-driver/mysql"
)

const (
	INSERTUSER         = "INSERT INTO user(email, name, lastname,age,status) VALUES (?, ?, ?,?,?)"
	SELECTUSERBYID     = "SELECT id, email, name, lastname,age,status FROM user WHERE id = ?"
	SELECTUSEERBYEMAIL = "SELECT id, email, name, lastname,age,status FROM user WHERE email = ?"
	SELECTALLUSERS     = "SELECT id, email, name, lastname,age,status FROM user"
	UPDATEUSER         = "UPDATE user SET name=?, lastname=?,age=?, status=? WHERE id = ?"
	DELETEUSER         = "UPDATE user SET status=0 WHERE id= ?"
)

type config struct {
	User      string `env:"MYSQL_USER" envDefault:"cfloress"`
	Password  string `env:"MYSQL_PASSWORD" envDefault:"cfloress.,2021"`
	Port      string `env:"MYSQL_PORT" envDefault:":3306"`
	Host      string `env:"MYSQL_HOST" envDefault:"localhost"`
	DefaultDB string `env:"MYSQL_DEFAULTDB" envDefault:"globant"`
}

func initMySQLRepository() (*sql.DB, error) {

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s%s)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DefaultDB)
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("initMysqlRepository...")

	return db, nil
}

//MySQLRepository - is a mysql implementation of users repository
type MySQLRepository struct {
	db *sql.DB
}

//NewMySQLUserRepository - returns a MySQLRepository type pointer
func NewMySQLUserRepository() (*MySQLRepository, error) {

	db, err := initMySQLRepository()

	if err != nil {
		return nil, err
	}
	fmt.Println("NewMySQLUserRepository...")

	return &MySQLRepository{
		db: db,
	}, nil
}

//Add - adds a user to the repository
func (r *MySQLRepository) Add(ctx context.Context, usr user.User) (int, error) {

	stmt, err := r.db.Prepare(INSERTUSER)

	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(usr.Email, usr.Name, usr.Lastname, usr.Age, usr.Status)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		//TODO log error
		return 0, err
	}

	return int(id), nil
}

//GetByID - retrieves a user from the repository based on the integer id
func (r *MySQLRepository) GetByID(ctx context.Context, userID int) (user.User, error) {

	usr := user.User{}

	err := r.db.QueryRow(SELECTUSERBYID, userID).
		Scan(&usr.ID, &usr.Email, &usr.Name, &usr.Lastname, &usr.Age, &usr.Status)

	if err == sql.ErrNoRows {
		return usr, nil
	}

	return usr, err

}

//GetByEmail - retrieves a user from the repository based on the email address
func (r *MySQLRepository) GetByEmail(ctx context.Context, email string) (user.User, error) {

	usr := user.User{}
	row := r.db.QueryRow(SELECTUSEERBYEMAIL, email)
	err := row.Scan(&usr.ID, &usr.Email, &usr.Name, &usr.Lastname, &usr.Age, &usr.Status)

	if err == sql.ErrNoRows {
		return usr, nil
	}

	return usr, err
}

//GetAll - retrieves all the users from the repository
func (r *MySQLRepository) GetAll(ctx context.Context) ([]user.User, error) {

	usrs := []user.User{}
	records, err := r.db.Query(SELECTALLUSERS)

	if err != nil {
		fmt.Println(err)
	}

	defer records.Close()

	for records.Next() {
		var userss user.User

		if err := records.Scan(&userss.ID, &userss.Email, &userss.Name, &userss.Lastname, &userss.Age, &userss.Status); err != nil {
			return nil, err
		}

		usrs = append(usrs, userss)
	}
	return usrs, nil
}

//Update -  updates the information of a user
func (r *MySQLRepository) Update(ctx context.Context, usr user.User) error {

	stmt, err := r.db.Prepare(UPDATEUSER)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(usr.Name, usr.Lastname, usr.Age, usr.Status, usr.ID)

	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); rows == 0 || err != nil {
		return errors.New("no records were updated")
	}

	return nil
}

//Delete - deletes a user from the repository
func (r *MySQLRepository) Delete(ctx context.Context, userID int) error {

	stmt, err := r.db.Prepare(DELETEUSER)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(userID)

	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("no records were affected")
	}
	return nil
}
