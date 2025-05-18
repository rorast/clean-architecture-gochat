package routers

import (
	"clean-architecture-gochat/docs"
	"clean-architecture-gochat/infrastructure/mysql"
	"clean-architecture-gochat/interface/controllers"
	"clean-architecture-gochat/internal/domain/repositories"
	"clean-architecture-gochat/usecases/chat"
	"clean-architecture-gochat/usecases/user"
	"clean-architecture-gochat/usecases/websocket"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Swagger API 文檔
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 設置靜態資源
	r.Static("/web/asset", "web/asset")
	r.StaticFile("/favicon.ico", "web/asset/images/favicon.ico")
	r.LoadHTMLGlob("web/views/**/*")

	// 初始化依賴
	db := mysql.Connect()

	userRepo := repositories.NewUserRepository(db)
	contactRepo := repositories.NewContactRepository(db)
	userService := user.NewService(userRepo, contactRepo)
	userController := controllers.NewUserController(userService)
	indexController := controllers.NewIndexController()

	// 聊天相關依賴
	messageRepo := repositories.NewMessageRepository(db)
	groupRepo := repositories.NewGroupRepository(db)
	privateChatService := chat.NewPrivateChatService(messageRepo)
	groupChatService := chat.NewGroupChatService(groupRepo, messageRepo)
	connectionService := websocket.NewConnectionService()
	messageService := chat.NewMessageService(messageRepo)
	chatController := controllers.NewChatController(privateChatService, groupChatService, connectionService, messageService)

	// 首頁相關路由
	r.GET("/", indexController.GetIndex)
	r.GET("/index", indexController.GetIndex)
	r.GET("/toRegister", indexController.ToRegister)
	r.GET("/register", indexController.ToRegister)
	r.GET("/toChat", indexController.ToChat)

	// 用戶模組
	userGroup := r.Group("/user")
	{
		userGroup.GET("/list", userController.GetUserList)
		userGroup.POST("/create", userController.CreateUser)
		userGroup.DELETE("/delete", userController.DeleteUser)
		userGroup.POST("/updateUser", userController.UpdateUser)
		userGroup.POST("/login", userController.FindUserByNameAndPwd)
		userGroup.POST("/searchFriends", userController.SearchFriend)
		userGroup.POST("/find", userController.FindUserByID)
	}

	// 檔案上傳模組
	attachGroup := r.Group("/attach")
	{
		attachGroup.POST("/upload", userController.UploadFile)
	}

	// 聊天模組
	contactGroup := r.Group("/contact")
	{
		// 添加好友
		contactGroup.POST("/addFriend", userController.AddFriend)
		// 獲取群組列表
		contactGroup.GET("/groups", chatController.GetGroups)
		// 創建群組
		contactGroup.POST("/create-group", chatController.CreateCustomGroup)
		// 加入群組
		contactGroup.POST("/joinGroup", chatController.JoinGroup)
	}

	// 聊天模組路由
	chatGroup := r.Group("/chat")
	{
		chatGroup.GET("/ws", chatController.HandleWebSocket)
		chatGroup.POST("/private/send", chatController.SendPrivateMessage)
		chatGroup.GET("/private/history", chatController.GetPrivateHistory)

		// 群組相關路由
		chatGroup.POST("/group/create", chatController.CreateGroup)
		chatGroup.PUT("/group/:id", chatController.UpdateGroup)
		chatGroup.DELETE("/group/:id", chatController.DeleteGroup)
		chatGroup.GET("/group/:id", chatController.GetGroup)
		chatGroup.GET("/groups", chatController.GetGroups)
		chatGroup.GET("/group/:id/members", chatController.GetGroupMembers)
		chatGroup.POST("/group/:id/members", chatController.AddGroupMember)
		chatGroup.DELETE("/group/:id/members", chatController.RemoveGroupMember)
		chatGroup.POST("/group/send", chatController.SendGroupMessage)
		chatGroup.GET("/group/history", chatController.GetGroupHistory)
	}

	return r
}
