package service

import (
	"divya07514-hftest-golang/db"
	"divya07514-hftest-golang/logger"
	"divya07514-hftest-golang/model"
)

type GlobalStats struct {
	Db db.DataSource
}

func (s *GlobalStats) Report(postCode string, fromTime string, toTime string, recipeNames []string) (*model.RecipeStats, error) {
	uniqueRecipeCount, err := s.Db.UniqueRecipeCount()
	if err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return nil, err
	}
	uniqueRecipeListAndCount, err := s.Db.UniqueRecipeAndCount()
	if err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return nil, err
	}
	busiestPostCode, err := s.Db.BusiestPostCode()
	if err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return nil, err
	}
	recipeList, err := s.Db.ListRecipeNames(recipeNames...)
	if err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return nil, err
	}
	deliveriesForPostCode, err := s.Db.DeliveriesPostCode(postCode, fromTime, toTime)
	if err != nil {
		logger.ErrorLogger.Error().Msg(err.Error())
		return nil, err
	}
	return &model.RecipeStats{
		UniqueRecipes:           int(uniqueRecipeCount),
		PerRecipe:               uniqueRecipeListAndCount,
		BusiestPostCode:         busiestPostCode,
		CountPerPostCodeAndTime: deliveriesForPostCode,
		MatchByName:             recipeList,
	}, nil
}
