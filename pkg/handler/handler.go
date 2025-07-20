package handler

import (
	"todo-app/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service.Authorization
	service.ToDoList
	service.ToDoItem
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		Authorization: services.Authorization,
		ToDoList:      services.ToDoList,
		ToDoItem:      services.ToDoItem,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList) // Создать список
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListByID)
			lists.PUT("/:id", h.updateListByID)
			lists.DELETE("/:id", h.deleteListByID)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem) // Создать задачу
				items.GET("/", h.getAllItems)
			}
		}

		items := api.Group("/items")
		{
			items.GET("/:id", h.getItemByID)
			items.PUT("/:id", h.updateItemsByID)
			items.DELETE("/:id", h.deleteItemsByID)
		}
	}
	return router
}
