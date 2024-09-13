package main

import (
	"GoHTMX/views"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Grocery struct {
	Title    string `form:"title" binding:"required"`
	Quantity string `form:"quantity" binding:"required"`
}

var data = newData()

func newItem(title string, quantity, id int) views.GroceryItem {
	return views.GroceryItem{
		Id:       id,
		Title:    title,
		Quantity: quantity,
	}
}

func newData() []views.GroceryItem {
	return []views.GroceryItem{
		newItem("Chicken", 1, 1),
		newItem("Beef", 1, 2),
	}
}

func createHomeTemplate(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "", views.Home(data))
}

func getGroceryList(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "", views.GroceryList(data))
}

func addItemToList(ctx *gin.Context) {
	var grocery Grocery

	if err := ctx.ShouldBind(&grocery); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	title := grocery.Title

	if err := validateTitle(title); err != nil {
		ctx.HTML(http.StatusOK, "", views.Validation(true))
		return
	}

	quantity, _ := strconv.Atoi(grocery.Quantity)
	id := data[len(data)-1].Id + 1
	newGrocery := newItem(title, quantity, id)

	data = append(data, newGrocery)

	ctx.HTML(http.StatusOK, "", views.Grocery(newGrocery))
	ctx.HTML(http.StatusOK, "", views.Validation(false))
}

func deleteItemFromList(ctx *gin.Context) {
	index := -1
	id, _ := strconv.Atoi(ctx.Param("id"))

	for i, item := range data {
		if item.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no ID found"})
		return
	}

	data = append(data[:index], data[index+1:]...)
}

func validateTitle(title string) error {
	for _, item := range data {
		if item.Title == title {
			return errors.New("title already exists")
		}
	}

	return nil
}



var infiniteScrollData = newInfiniteScrollData()

func newInfiniteScrollData() []int {
	numbers := make([]int, 301)

	for i := 0; i <= 300; i++ {
		numbers[i] = i
	}

	return numbers
}

func createInfiniteScrollDemo(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "", views.InfiniteScroll())
}

func getMoreStuff(ctx *gin.Context) {
	start, _ := strconv.Atoi(ctx.Param("start"))
	end := start + 11

	if start == len(infiniteScrollData) - 1 {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "no results found"})
		return
	}

	if end > len(infiniteScrollData) {
		end = len(infiniteScrollData)
	}

	data := append([]int{}, infiniteScrollData[start:end]...)

	ctx.HTML(http.StatusOK, "", views.InfiniteScrollItems(data))
}

func main() {
	router := gin.Default()
	router.HTMLRender = &TemplRender{}
	router.Static("/styles", "./styles")
	router.SetTrustedProxies(nil)

	router.GET("/", createHomeTemplate)
	router.GET("/groceries", getGroceryList)
	router.POST("/groceries", addItemToList)
	router.DELETE("/groceries/:id", deleteItemFromList)

	router.GET("/infinite-scroll", createInfiniteScrollDemo)
	router.GET("/infinite-scroll/:start", getMoreStuff)

	router.Run(":8080")
}
