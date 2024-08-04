package db

import (
	"database/sql"
	"divya07514-hftest-golang/logger"
	"divya07514-hftest-golang/model"
	"divya07514-hftest-golang/util"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

//go:generate mockery --name=DataSource --inpackage --case=underscore
type DataSource interface {
	CreateTable() error
	Insert(data model.RecipeData)
	InsertBatch(dataChannel <-chan []model.RecipeData)
	UniqueRecipeCount() (int64, error)
	UniqueRecipeAndCount() ([]*model.PerRecipeStats, error)
	BusiestPostCode() (*model.PostCodeStats, error)
	DeliveriesPostCode(postCode string, fromTime string, toTime string) (*model.PostCodeAndTimeStats, error)
	ListRecipeNames(names ...string) ([]string, error)
}

type MySqlDB struct {
	Connection *sql.DB
}

func (db *MySqlDB) CreateTable() error {
	_, err := db.Connection.Exec(CreateTable)
	if err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return err
	}
	return nil
}

func (db *MySqlDB) Insert(data model.RecipeData) {
	from, to, dayOfWeek := util.GetFromAndToTimes(data.Delivery)
	_, err := db.Connection.Exec(InsertRecipe, data.Postcode, data.Recipe, from, to, dayOfWeek)
	if err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return
	}
}

func (db *MySqlDB) InsertBatch(dataChannel <-chan []model.RecipeData) {
	stmt, err := db.Connection.Prepare(InsertRecipe)
	if err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
	}
	defer stmt.Close()
	for batch := range dataChannel {
		if err := insertBatch(db.Connection, stmt, batch); err != nil {
			logger.ErrorLogger.Error().Msg(err.Error())
		}
		log.Info().Msg("Done inserting batch")
	}
}

// UniqueRecipeCount Get count of unique recipe names
func (db *MySqlDB) UniqueRecipeCount() (int64, error) {
	rows := db.Connection.QueryRow(UniqueRecipeCount)
	var uniqueCount int64
	err := rows.Scan(&uniqueCount)
	if err != nil {
		log.Error().Msg(err.Error())
		return 0, err
	}
	return uniqueCount, nil
}

// UniqueRecipeAndCount Get unique recipes with their counts
func (db *MySqlDB) UniqueRecipeAndCount() ([]*model.PerRecipeStats, error) {
	rows, _ := db.Connection.Query(UniqueRecipeAndCount)
	var resultSet []*model.PerRecipeStats
	for rows.Next() {
		var recipe model.PerRecipeStats
		if err := rows.Scan(&recipe.Count, &recipe.Recipe); err != nil {
			logger.ErrorLogger.Error().Msg(err.Error())
			return nil, err
		}
		resultSet = append(resultSet, &recipe)
	}
	return resultSet, nil
}

// BusiestPostCode Get busiest post code and it's count
func (db *MySqlDB) BusiestPostCode() (*model.PostCodeStats, error) {
	rows := db.Connection.QueryRow(BusiestPostCode)
	var recipe model.PostCodeStats
	if err := rows.Scan(&recipe.DeliveryCount, &recipe.Postcode); err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return &model.PostCodeStats{}, err
	}
	return &recipe, nil
}

// DeliveriesPostCode Get stats about deliveries for a given postcode and between a given time range. Time range is inclusive of given values
func (db *MySqlDB) DeliveriesPostCode(postCode string, fromTime string, toTime string) (*model.PostCodeAndTimeStats, error) {
	fromTimeVal := util.FormatTimeToInteger(fromTime)
	toTimeVal := util.FormatTimeToInteger(toTime)
	if toTimeVal < fromTimeVal {

		return nil, errors.New("to_time cannot be less than from_time")
	}
	rows := db.Connection.QueryRow(DeliveriesToPostcode, postCode, fromTimeVal, toTimeVal)
	var recipe model.PostCodeAndTimeStats
	if err := rows.Scan(&recipe.DeliveryCount); err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return &model.PostCodeAndTimeStats{}, err
	}
	recipe.FromTime = fromTime
	recipe.ToTime = toTime
	recipe.PostCode = postCode
	return &recipe, nil
}

// ListRecipeNames Return recipe names where input is like names ["Potato","Veggies"]
// query will return data where names are like above sample input
func (db *MySqlDB) ListRecipeNames(searchTerms ...string) ([]string, error) {
	likeClauses := make([]string, len(searchTerms))
	for i, _ := range searchTerms {
		likeClauses[i] = "recipe LIKE ?"
	}
	query := fmt.Sprintf(
		RecipesLike,
		strings.Join(likeClauses, " OR "),
	)
	args := make([]interface{}, len(searchTerms))
	for i, term := range searchTerms {
		args[i] = "%" + term + "%"
	}
	rows, err := db.Connection.Query(query, args...)
	if err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return nil, err
	}
	var resultSet []string
	for rows.Next() {
		var recipe string
		if err = rows.Scan(&recipe); err != nil {
			logger.ErrorLogger.Error().Msg(err.Error())
			return nil, err
		}
		resultSet = append(resultSet, recipe)
	}
	return resultSet, nil
}

func GetMySqlDbConnection(dataSource string) (*sql.DB, error) {
	var dataBase *sql.DB
	dataBase, err := sql.Open("mysql", dataSource)
	if err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return nil, err
	}
	_, err = dataBase.Exec(CreateDb)
	if err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return nil, err
	}
	return dataBase, err
}

func insertBatch(db *sql.DB, stmt *sql.Stmt, batch []model.RecipeData) error {
	tx, err := db.Begin()
	if err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return err
	}
	defer tx.Rollback()
	for _, data := range batch {
		from, to, dayOfWeek := util.GetFromAndToTimes(data.Delivery)
		_, err := stmt.Exec(data.Postcode, data.Recipe, from, to, dayOfWeek)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
