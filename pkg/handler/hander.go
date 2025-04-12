package handler

import (
	"ToDoGo/pkg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	//todo Настроить CORS
	router.Use(cors.Default())
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	api := router.Group("/api")
	{
		v1 := api.Group("/v1", h.userIdentity)
		{
			lists := v1.Group("/lists")
			{
				lists.POST("/", h.createList)      //Done
				lists.GET("/", h.getAllLists)      //Done
				lists.GET("/:id", h.getListById)   //Done
				lists.PUT("/:id", h.updateList)    //Done
				lists.DELETE("/:id", h.deleteList) //Done

				items := lists.Group(":id/items")
				{
					items.POST("/", h.createItem)         //Done
					items.GET("/", h.getAllItems)         //Done
					items.GET("/:item-id", h.getItemById) //Done
					items.PUT("/:item-id", h.updateItem)
					items.DELETE("/:item-id", h.deleteItem)
				}
			}
		}
		user := api.Group("/user")
		{
			user.GET("/", h.getUserInfo)                        // Get user info
			user.POST("/reset-request", h.passwordResetRequest) // Request reset user password
			user.POST("/reset-confirm/:token", h.passwordResetConfirm)
			user.DELETE("/delete", h.deleteUser) // Delete user
		}
	}
	return router
}
