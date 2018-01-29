package neatly_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/neatly"
	"github.com/viant/toolbox/data"
	"github.com/viant/toolbox/url"
	"testing"
	"time"
)

type Request1 struct {
	URL     string
	Method  string
	Cookies map[string]string
}

type ExpectedResponse1 struct {
	StatusCode int
	Body       string
}

type UseCase1 struct {
	UseCase  string
	Requests []*Request1
	Expect   []*ExpectedResponse1
	Comments string
}

func TestDao_LoadUseCase1(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()

	var useCase1 = &UseCase1{}
	err := dao.Load(context, url.NewResource("test/use_case1.csv"), useCase1)
	assert.Nil(t, err)

	assert.Equal(t, "case 1", useCase1.UseCase)
	assert.Equal(t, "{Root123}", useCase1.Comments)
	assert.Equal(t, 2, len(useCase1.Requests))
	{
		request := useCase1.Requests[0]
		assert.Equal(t, "http://127.0.0.1/test1", request.URL)
		assert.Equal(t, "GET", request.Method)
		assert.Equal(t, 2, len(request.Cookies))
		assert.Equal(t, "value1", request.Cookies["Cookie1"])
		assert.Equal(t, "value2", request.Cookies["Cookie2"])
	}
	{
		request := useCase1.Requests[1]
		assert.Equal(t, "http://127.0.0.1/test2", request.URL)
		assert.Equal(t, "GET", request.Method)
		assert.Equal(t, 2, len(request.Cookies))
		assert.Equal(t, "value1", request.Cookies["Cookie1"])
		assert.Equal(t, "value2", request.Cookies["Cookie2"])
	}
	assert.Equal(t, 2, len(useCase1.Expect))

	{
		expect := useCase1.Expect[0]
		assert.Equal(t, 200, expect.StatusCode)
		assert.Equal(t, "test1Content", expect.Body)
	}

	{
		expect := useCase1.Expect[1]
		assert.Equal(t, 404, expect.StatusCode)
	}

}

type LineItem struct {
	Product  string
	Quantity int
	Price    float64
}

type Order struct {
	Id        int
	Name      string
	SubTotal  float64
	LineItems []*LineItem
}

type UseCase2 struct {
	CreateTime time.Time
	Orders     []Order
}

func TestDao_LoadUseCase2(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "yyyy-MM-dd h:mm:ss", nil)
	var context = data.NewMap()
	var useCase2 = &UseCase2{}
	err := dao.Load(context, url.NewResource("test/use_case2.csv"), useCase2)
	assert.Nil(t, err)

	assert.Equal(t, 2017, useCase2.CreateTime.Year())
	assert.Equal(t, 2, len(useCase2.Orders))
	{
		order := useCase2.Orders[0]
		assert.Equal(t, 1, order.Id)
		assert.Equal(t, "Order 1", order.Name)
		assert.Equal(t, 100.0, order.SubTotal)
		assert.Equal(t, 2, len(order.LineItems))
		{
			lineItem := order.LineItems[0]
			assert.Equal(t, "Magic Mouse", lineItem.Product)
			assert.Equal(t, 5, lineItem.Quantity)
			assert.Equal(t, 10.0, lineItem.Price)

		}
		{
			lineItem := order.LineItems[1]
			assert.Equal(t, "TrackPad", lineItem.Product)
			assert.Equal(t, 5, lineItem.Quantity)
			assert.Equal(t, 10.0, lineItem.Price)

		}

	}
	{
		order := useCase2.Orders[1]
		assert.Equal(t, 2, order.Id)
		assert.Equal(t, "Order 2", order.Name)
		assert.Equal(t, 150.0, order.SubTotal)
		assert.Equal(t, 2, len(order.LineItems))
		{
			lineItem := order.LineItems[0]
			assert.Equal(t, "Keyboard", lineItem.Product)
			assert.Equal(t, 10, lineItem.Quantity)
			assert.Equal(t, 10.0, lineItem.Price)

		}
		{
			lineItem := order.LineItems[1]
			assert.Equal(t, "TrackPad", lineItem.Product)
			assert.Equal(t, 5, lineItem.Quantity)
			assert.Equal(t, 10.0, lineItem.Price)

		}

	}
}

type Merit struct {
	EmpNo       int
	Description string
}

type Bonus struct {
	EmpNo  int
	Name   string
	Amount float64
}

type UseCase3 struct {
	Merits  []*Merit
	Bonus   []*Bonus
	Created time.Time
}

func TestDao_LoadUseCase3(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "yyyy-MM-dd h:mm:ss", nil)
	var context = data.NewMap()
	var useCase3 = &UseCase3{}
	err := dao.Load(context, url.NewResource("test/use_case3.csv"), useCase3)
	assert.Nil(t, err)
	assert.Equal(t, 2017, useCase3.Created.Year())
	assert.Equal(t, 3, len(useCase3.Merits))

	for i := 0; i < 3; i++ {
		var merit = useCase3.Merits[i]
		assert.Equal(t, i+1, merit.EmpNo)
	}

	if assert.Equal(t, 3, len(useCase3.Bonus)) {
		var bonuses = []float64{10000, 8000, 4000}
		for i := 0; i < 3; i++ {
			var bonus = useCase3.Bonus[i]
			assert.Equal(t, i+1, bonus.EmpNo)
			assert.Equal(t, bonuses[i], bonus.Amount)
		}
	}

}

func TestDao_LoadUseCase4(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "yyyy-MM-dd h:mm:ss", nil)
	var context = data.NewMap()
	var useCase3 = &UseCase3{}
	err := dao.Load(context, url.NewResource("test/use_case4.csv"), useCase3)
	assert.Nil(t, err)
	assert.Equal(t, 2017, useCase3.Created.Year())
	assert.Equal(t, 3, len(useCase3.Merits))
	assert.Equal(t, 3, len(useCase3.Bonus))

	for i := 0; i < 3; i++ {
		var merit = useCase3.Merits[i]
		assert.Equal(t, i+1, merit.EmpNo)
	}

	var bonuses = []float64{10000, 8000, 4000}
	for i := 0; i < 3; i++ {
		var bonus = useCase3.Bonus[i]
		assert.Equal(t, i+1, bonus.EmpNo)
		assert.Equal(t, bonuses[i], bonus.Amount)
	}

}

type Repeated struct {
	Id   int
	Name string
}

type UseCase5 struct {
	Repeated []*Repeated
}

func TestDao_LoadUseCase5(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()
	var useCase5 = &UseCase5{}
	err := dao.Load(context, url.NewResource("test/use_case5.csv"), useCase5)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(useCase5.Repeated))

	for i := 1; i <= 5; i++ {
		assert.Equal(t, useCase5.Repeated[i-1].Id, i)
		assert.Equal(t, useCase5.Repeated[i-1].Name, fmt.Sprintf("Name %02d", i))
	}
}

type SubjectScore struct {
	Subject string
	Score   float64
}

type Student struct {
	Id     int
	Name   string
	Scores []*SubjectScore
}

type UseCase6 struct {
	Students []*Student
}

func TestDao_LoadUseCase6(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()
	var useCase6 = &UseCase6{}
	err := dao.Load(context, url.NewResource("test/use_case6.csv"), useCase6)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(useCase6.Students))
	{
		var student = useCase6.Students[0]
		assert.Equal(t, 1, student.Id)
		assert.Equal(t, "Smith", student.Name)
		assert.Equal(t, 2, len(student.Scores))
		assert.Equal(t, "Math", student.Scores[0].Subject)
		assert.Equal(t, 3.2, student.Scores[0].Score)
		assert.Equal(t, "English", student.Scores[1].Subject)
		assert.Equal(t, 3.5, student.Scores[1].Score)

	}
	{
		var student = useCase6.Students[1]
		assert.Equal(t, 2, student.Id)
		assert.Equal(t, "Kowalczyk", student.Name)
		assert.Equal(t, 2, len(student.Scores))
		assert.Equal(t, "Math", student.Scores[0].Subject)
		assert.Equal(t, 3.7, student.Scores[0].Score)
		assert.Equal(t, "English", student.Scores[1].Subject)
		assert.Equal(t, 3.2, student.Scores[1].Score)

	}

}

type UseCase7Item struct {
	Id          int
	Description string
}

type UseCase7 struct {
	Setup    map[string]map[string][]map[string]interface{}
	UseCases []*UseCase7Item
}

func TestDao_LoadUseCase7(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()
	var useCase7 = &UseCase7{}
	err := dao.Load(context, url.NewResource("test/use_case7.csv"), useCase7)
	if ! assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 2, len(useCase7.UseCases))

	for i := 0; i < 2; i++ {
		assert.Equal(t, i+1, useCase7.UseCases[i].Id)
		assert.Equal(t, fmt.Sprintf("use case %d", i+1), useCase7.UseCases[i].Description)
	}

	mydb, ok := useCase7.Setup["MyDb"]
	if assert.True(t, ok) {
		customers, ok := mydb["Customer"]
		if assert.True(t, ok) {
			assert.Equal(t, 2, len(customers))

			{
				var customer = customers[0]
				assert.Equal(t, 1.0, customer["ID"])
				assert.EqualValues(t, "Smith", customer["NAME"])
				assert.EqualValues(t, 100.0, customer["DAILY_CAP"])
				assert.EqualValues(t, 1000.0, customer["OVERALL_CAP"])

			}
			{
				var customer = customers[1]
				assert.Equal(t, 2.0, customer["ID"])
				assert.EqualValues(t, "Kowalczyk", customer["NAME"])
				assert.EqualValues(t, 400.0, customer["DAILY_CAP"])
				assert.EqualValues(t, 8000.0, customer["OVERALL_CAP"])

			}
		}
	}

}

func TestDao_LoadUseCase8(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()
	var useCase7 = &UseCase7{}
	err := dao.Load(context, url.NewResource("test/use_case8.csv"), useCase7)
	if ! assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 2, len(useCase7.UseCases))

	for i := 0; i < 2; i++ {
		assert.Equal(t, i+1, useCase7.UseCases[i].Id)
		assert.Equal(t, fmt.Sprintf("use case %d", i+1), useCase7.UseCases[i].Description)
	}

	mydb, ok := useCase7.Setup["MyDb"]
	if assert.True(t, ok) {
		customers, ok := mydb["Customer"]
		if assert.True(t, ok) {
			assert.Equal(t, 2, len(customers))

			{
				var customer = customers[0]
				assert.Equal(t, 1.0, customer["ID"])
				assert.EqualValues(t, "Smith", customer["NAME"])
				assert.EqualValues(t, 100.0, customer["DAILY_CAP"])
				assert.EqualValues(t, 1000.0, customer["OVERALL_CAP"])

			}
			{
				var customer = customers[1]
				assert.Equal(t, 2.0, customer["ID"])
				assert.EqualValues(t, "Kowalczyk", customer["NAME"])
				assert.EqualValues(t, 100.0, customer["DAILY_CAP"])
				assert.EqualValues(t, 1000.0, customer["OVERALL_CAP"])

			}
		}
	}

}

func TestDao_LoadUseCase9(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()
	var useCase7 = &UseCase7{}
	err := dao.Load(context, url.NewResource("test/use_case9.csv"), useCase7)
	if ! assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 2, len(useCase7.UseCases))

	for i := 0; i < 2; i++ {
		assert.Equal(t, i+1, useCase7.UseCases[i].Id)
		assert.Equal(t, fmt.Sprintf("use case %d", i+1), useCase7.UseCases[i].Description)
	}

	mydb, ok := useCase7.Setup["MyDb"]
	if assert.True(t, ok) {
		customers, ok := mydb["Customer"]
		if assert.True(t, ok) {
			assert.Equal(t, 2, len(customers))

			{
				var customer = customers[0]
				assert.Equal(t, 1.0, customer["ID"])
				assert.EqualValues(t, "Smith", customer["NAME"])
				assert.EqualValues(t, 100.0, customer["DAILY_CAP"])
				assert.EqualValues(t, 1000.0, customer["OVERALL_CAP"])

			}
			{
				var customer = customers[1]
				assert.Equal(t, 2.0, customer["ID"])
				assert.EqualValues(t, "Kowalczyk", customer["NAME"])
				assert.EqualValues(t, 100.0, customer["DAILY_CAP"])
				assert.EqualValues(t, 1000.0, customer["OVERALL_CAP"])

			}
		}
	}

}

func TestDao_LoadUseCase10(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()
	var useCase7 = &UseCase7{}
	err := dao.Load(context, url.NewResource("test/use_case10.csv"), useCase7)
	if ! assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 2, len(useCase7.UseCases))

	for i := 0; i < 2; i++ {
		assert.Equal(t, i+1, useCase7.UseCases[i].Id)
		assert.Equal(t, fmt.Sprintf("use case %d", i+1), useCase7.UseCases[i].Description)
	}

	mydb, ok := useCase7.Setup["MyDb"]
	if assert.True(t, ok) {
		customers, ok := mydb["Customer"]
		if assert.True(t, ok) {
			assert.Equal(t, 2, len(customers))

			{
				var customer = customers[0]
				assert.Equal(t, 1.0, customer["ID"])
				assert.EqualValues(t, "Smith", customer["NAME"])
				assert.EqualValues(t, "200", customer["DAILY_CAP"])
				assert.EqualValues(t, "3000", customer["OVERALL_CAP"])

			}
			{
				var customer = customers[1]
				assert.Equal(t, 2.0, customer["ID"])
				assert.EqualValues(t, "Kowalczyk", customer["NAME"])
				assert.EqualValues(t, "100", customer["DAILY_CAP"])
				assert.EqualValues(t, "1000", customer["OVERALL_CAP"])

			}
		}
	}

}

func TestDao_LoadUseCase11(t *testing.T) {
	dao := neatly.NewDao(true, "", "", "", nil)
	var context = data.NewMap()
	var useCase7 = &UseCase7{}
	err := dao.Load(context, url.NewResource("test/use_case11.csv"), useCase7)

	assert.Nil(t, err)
	if ! assert.Nil(t, err) {
		return
	}
	for i := 0; i < 2; i++ {
		assert.Equal(t, i+1, useCase7.UseCases[i].Id)
		assert.Equal(t, fmt.Sprintf("use case %d", i+1), useCase7.UseCases[i].Description)
	}

	mydb, ok := useCase7.Setup["MyDb"]
	if assert.True(t, ok) {
		customers, ok := mydb["Customer"]
		if assert.True(t, ok) {
			assert.Equal(t, 2, len(customers))

			{
				var customer = customers[0]
				assert.Equal(t, 1.0, customer["ID"])
				assert.EqualValues(t, "Smith", customer["NAME"])
				assert.EqualValues(t, "200", customer["DAILY_CAP"])
				assert.EqualValues(t, "3000", customer["OVERALL_CAP"])

			}
			{
				var customer = customers[1]
				assert.Equal(t, 2.0, customer["ID"])
				assert.EqualValues(t, "Kowalczyk", customer["NAME"])
				assert.EqualValues(t, "100", customer["DAILY_CAP"])
				assert.EqualValues(t, "1000", customer["OVERALL_CAP"])

			}
		}
	}

}

func TestDao_LoadUseCase12(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()
	var useCase7 = &UseCase7{}
	err := dao.Load(context, url.NewResource("test/use_case12.csv"), useCase7)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(useCase7.UseCases))

	for i := 0; i < 2; i++ {
		assert.Equal(t, i+1, useCase7.UseCases[i].Id)
		assert.Equal(t, fmt.Sprintf("use case %d", i+1), useCase7.UseCases[i].Description)

	}

}

type UseCase13 struct {
	Prepare []struct {
		TagId string
	}
	Data map[string][]struct {
		Autoincrement []string
		Table         string
		Value         interface{}
	}
	Numbers struct {
		Seq []int
	}
}

func TestDao_LoadUseCase13(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()

	var useCase = &UseCase13{}
	err := dao.Load(context, url.NewResource("test/use_case13.csv"), &useCase)

	if assert.Nil(t, err) {
		assert.Equal(t, 2, len(useCase.Data["mydb1"]))
		assert.EqualValues(t, []string{"meta.USER", "meta.ACCOUNT"}, useCase.Data["mydb1"][0].Autoincrement)
		assert.EqualValues(t, []string{"meta.USER"}, useCase.Data["mydb1"][1].Autoincrement)

		assert.Equal(t, 3, len(useCase.Numbers.Seq))
	}

}

type UseCase14Action struct {
	Send struct {
		Udf      string
		Requests []struct {
			Method string
			URL    string
		}
	}
	Expect []struct {
		Code string
	}
}

type UseCase14 struct {
	Actions []*UseCase14Action
}

func TestDao_LoadUseCase14(t *testing.T) {

	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()

	var useCase = &UseCase14{}
	err := dao.Load(context, url.NewResource("test/use_case14.csv"), &useCase)

	if assert.Nil(t, err) {
		assert.EqualValues(t, 3, len(useCase.Actions))

		assert.EqualValues(t, "MyUdf", useCase.Actions[0].Send.Udf)
		assert.EqualValues(t, 2, len(useCase.Actions[0].Send.Requests))
		assert.EqualValues(t, "http://127.0.0.1/path1", useCase.Actions[0].Send.Requests[0].URL)
		assert.EqualValues(t, "http://127.0.0.1/path2", useCase.Actions[0].Send.Requests[1].URL)

		assert.EqualValues(t, "MyUdf", useCase.Actions[1].Send.Udf)
		if assert.EqualValues(t, 2, len(useCase.Actions[1].Send.Requests)) {
			assert.EqualValues(t, "http://127.0.0.1/path3", useCase.Actions[1].Send.Requests[0].URL)
			assert.EqualValues(t, "http://127.0.0.1/path4", useCase.Actions[1].Send.Requests[1].URL)
		}

		if assert.EqualValues(t, "MyUdf", useCase.Actions[2].Send.Udf) {
			assert.EqualValues(t, 2, len(useCase.Actions[2].Send.Requests))
			assert.EqualValues(t, "http://127.0.0.1/path5", useCase.Actions[2].Send.Requests[0].URL)
			assert.EqualValues(t, "http://127.0.0.1/path6", useCase.Actions[2].Send.Requests[1].URL)
		}

	}

}

func TestMissingReference(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()
	var document = make(map[string]interface{})
	err := dao.Load(context, url.NewResource("test/broken1.csv"), &document)
	assert.NotNil(t, err)
}

func TestBrokenJsonReference(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()
	var document = make(map[string]interface{})
	err := dao.Load(context, url.NewResource("test/broken2.csv"), &document)
	assert.NotNil(t, err)
}

func TestBrokenExternalReference(t *testing.T) {
	dao := neatly.NewDao(false, "", "", "", nil)
	var context = data.NewMap()
	var document = make(map[string]interface{})
	err := dao.Load(context, url.NewResource("test/broken3.csv"), &document)
	assert.NotNil(t, err)
}
