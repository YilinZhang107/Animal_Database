/*
* @Author: Oatmeal107
* @Date:   2023/6/12 16:15
 */

package Init

import (
	"Animal_database/api/v1"
	"Animal_database/config"
	"Animal_database/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() {
	r := gin.Default()
	//ginServer.Use(favicon.New("./webIcon.jpg")) //设置网站图标
	// 强制日志颜色化
	gin.ForceConsoleColor()

	//跨域
	r.Use(middleware.Cors())

	//todo 引入swagger

	//加载静态页面
	r.LoadHTMLGlob("./resources/templates/*")
	//加载静态文件路径,将 /static 路径映射到本地文件系统中的 ./static 目录，用于提供静态文件服务。
	r.StaticFS("/static", http.Dir("./resources/static"))

	//主页
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"msg": "后台传来的数据"})
	})

	apiv1 := r.Group("/api/v1")
	{
		//测试
		apiv1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})

		//test
		//用户
		apiv1.POST("user/register", v1.UserRegister) // 用户注册
		apiv1.POST("user/login", v1.UserLogin)       // 用户登陆
		//Record用于展示主页地图
		apiv1.POST("record/getByArea", v1.GetByArea) // 按地区查询现存记录

		// 需要登陆保护
		authed := apiv1.Group("/")
		authed.Use(middleware.JWT()) //jwt验证
		{
			//user
			authed.GET("user/selfInfo", v1.SelfInfo) // 获取自身信息
			authed.POST("user/avatar", v1.UserUpdateAvatar)
			authed.POST("user/changePassword", v1.UserChangePassword)
			authed.POST("user/updateEmail", v1.UserUpdateEmail)
			authed.POST("user/findUsers", v1.FindUsers)             // 按条件查看用户, 需要等级验证, 用于之后修改等级
			authed.POST("user/updateUserGrade", v1.UpdateUserGrade) // 修改用户等级
			authed.POST("user/deleteUser", v1.DeleteUser)           // 删除用户

			//UnreviewedRecord
			authed.POST("record/uploadByExcel", v1.UploadRecord)    // 上传记录
			authed.GET("record/getURecord", v1.GetUnreviewedRecord) // 查询待审批记录
			authed.POST("record/review", v1.ReviewURecord)          // 审批待审批记录

			//Record
			authed.GET("record/getRecord", v1.GetRecord) // 查询现存记录

		}

	}

	r.Run(config.ServerPort)

}
