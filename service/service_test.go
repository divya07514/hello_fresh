package service

import (
	"divya07514-hftest-golang/db"
	"divya07514-hftest-golang/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_report(t *testing.T) {
	dbMock := &db.MockDataSource{}
	dbMock.On("UniqueRecipeCount").Return(int64(1), nil)
	recipes := []*model.PerRecipeStats{
		{Recipe: "Cajun Chicken", Count: 150},
		{Recipe: "Beef Stroganoff", Count: 85},
	}
	dbMock.On("UniqueRecipeAndCount").Return(recipes, nil)
	postCode := &model.PostCodeStats{
		Postcode:      "10120",
		DeliveryCount: int64(123),
	}
	dbMock.On("BusiestPostCode").Return(postCode, nil)
	postCodeAndTime := &model.PostCodeAndTimeStats{
		PostCode:      "10120",
		FromTime:      "10AM",
		ToTime:        "3PM",
		DeliveryCount: int64(557),
	}
	dbMock.On("DeliveriesPostCode", mock.Anything, mock.Anything, mock.Anything).Return(postCodeAndTime, nil)
	dbMock.On("ListRecipeNames", mock.Anything, mock.Anything, mock.Anything).Return([]string{"Potato", "Tomato"}, nil)
	stats := GlobalStats{Db: dbMock}
	report, err := stats.Report("10120", "10AM", "3PM", []string{"Potato", "Steak"})
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, report.BusiestPostCode.Postcode, "10120")
	assert.Equal(t, report.BusiestPostCode.DeliveryCount, int64(123))
	assert.Equal(t, report.CountPerPostCodeAndTime.PostCode, "10120")
	assert.Equal(t, report.CountPerPostCodeAndTime.FromTime, "10AM")
	assert.Equal(t, report.CountPerPostCodeAndTime.ToTime, "3PM")
	assert.Equal(t, report.CountPerPostCodeAndTime.DeliveryCount, int64(557))
	mock.AssertExpectationsForObjects(t, dbMock)
	dbMock.AssertNumberOfCalls(t, "UniqueRecipeCount", 1)
	dbMock.AssertNumberOfCalls(t, "UniqueRecipeAndCount", 1)
	dbMock.AssertNumberOfCalls(t, "BusiestPostCode", 1)
	dbMock.AssertNumberOfCalls(t, "DeliveriesPostCode", 1)
	dbMock.AssertNumberOfCalls(t, "ListRecipeNames", 1)
}
