package services

import "testing"
import "alidada/models"

type sortTest struct {
	cancellationConditions       *[]models.CancellationCondition
	start                        int
	end                          int
	sortedCancellationConditions []models.CancellationCondition
}

var sortTests = []sortTest{
	sortTest{&[]models.CancellationCondition{
		{1, 10, "10 min", "", 90},
		{2, 20, "20 min", "", 80},
		{3, 20, "20 min", "", 70},
	}, 0, 2,
		[]models.CancellationCondition{
			{3, 20, "20 min", "", 70},
			{2, 20, "20 min", "", 80},
			{1, 10, "10 min", "", 90},
		},
	},
	sortTest{&[]models.CancellationCondition{
		{1, 10, "10 min", "", 10},
		{2, 20, "20 min", "", 20},
		{3, 20, "20 min", "", 30},
	}, 0, 2,
		[]models.CancellationCondition{
			{1, 10, "10 min", "", 10},
			{2, 20, "20 min", "", 20},
			{3, 20, "20 min", "", 30},
		},
	},
	sortTest{&[]models.CancellationCondition{
		{1, 10, "10 min", "", 80},
		{2, 20, "20 min", "", 10},
		{3, 20, "20 min", "", 12},
		{4, 20, "20 min", "", 70},
	}, 0, 3,
		[]models.CancellationCondition{
			{2, 20, "20 min", "", 10},
			{3, 20, "20 min", "", 12},
			{4, 20, "20 min", "", 70},
			{1, 10, "10 min", "", 80},
		},
	},
}

func TestSort(t *testing.T) {

	for j, test := range sortTests {
		output := Sort(test.cancellationConditions, test.start, test.end)
		for i, out := range output {
			if out.ID != test.sortedCancellationConditions[i].ID {
				t.Errorf("Output %d not equal to expected %d test number %d", out.ID, test.sortedCancellationConditions[i].ID, j+1)
			}
		}
	}
}
