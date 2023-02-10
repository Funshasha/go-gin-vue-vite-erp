package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
	"time"
)

type List struct {
	gorm.Model
	Name    string `gorm:"type:varchar(20); not null" json:"name" binding:"required"`
	State   string `gorm:"type:varchar(20); not null" json:"state" binding:"required"`
	Phone   string `gorm:"type:varchar(20); not null" json:"phone" binding:"required"`
	Email   string `gorm:"type:varchar(40); not null" json:"email" binding:"required"`
	Address string `gorm:"type:varchar(200); not null" json:"address" binding:"required"`
}

func main() {
	//链接到mysql数据库
	dsn := "root:lmy33063306@tcp(127.0.0.1:3306)/crud-list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 解决查表的时候自动添加复数的问题， 例如 user 变成 users
			SingularTable: true,
		},
	})

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("数据库打开 err", err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Second * 10) //10S

	db.AutoMigrate(&List{})

	//创建接口
	r := gin.Default()
	//test
	//r.GET("/", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "请求成功",
	//	})
	//})
	//var list List

	//增加
	r.POST("user/add", func(c *gin.Context) {
		var data List
		err := c.ShouldBindJSON(&data)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "添加失败",
				"data": gin.H{},
				"code": 400,
			})
		} else {

			//数据库操作
			db.Create(&data) //创建一条数据
			c.JSON(200, gin.H{
				"msg":  "添加成功",
				"data": data,
				"code": 200,
			})
		}
	})

	//删除 	通过对应id 进行删除
	r.DELETE("user/delete/:id", func(c *gin.Context) {
		//接收id
		var data []List
		id := c.Param("id")

		//判断 id 是否存在
		db.Where("id=?", id).Find(&data)
		if len(data) == 0 {
			c.JSON(200, gin.H{
				"msg":  "id 没有找到...",
				"code": 400,
			})
		} else {
			//操作数据库 删除
			db.Where("id=?", id).Delete(&data)
			c.JSON(200, gin.H{
				"msg":  "删除成功",
				"code": 200,
			})
		}
	})

	//修改
	r.PUT("user/update/:id", func(c *gin.Context) {
		var data List
		//接收 id
		id := c.Param("id")
		//判断 id 是否存在
		//db.Where("id = ?", id).Find(&data)
		db.Select("id").Where("id = ?", id).Find(&data)

		//判断id是否存在
		if data.ID == 0 {
			c.JSON(200, gin.H{
				"msg":  "用户id不存在",
				"code": 400,
			})
		} else {
			err := c.ShouldBindJSON(&data)
			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "修改失败",
					"code": 400,
				})
			} else {

				//修改数据库内容
				db.Where("id = ?", id).Updates(&data)

				c.JSON(200, gin.H{
					"msg":  "修改成功",
					"code": 200,
				})
			}
		}
	})

	//通过名字查询
	r.GET("/user/list/:name", func(c *gin.Context) {

		//获取路径参数
		name := c.Param("name")
		var dataList []List

		//查询数据库
		db.Where("name = ?", name).Find(&dataList)

		//判断是否查询到该数据
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询到数据",
				"code": 200,
				"data": dataList,
			})
		}
	})

	// 全部查询
	r.GET("user/list", func(c *gin.Context) {
		var dataList []List

		//查询全部数据，查询分页数据

		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))

		//判断是否是要分页
		if pageSize == 0 {
			pageSize = -1
		}

		if pageNum == 0 {
			pageNum = -1
		}

		offsetVal := (pageNum - 1) * pageSize
		if pageNum == -1 && pageSize == -1 {
			offsetVal = -1
		}
		var total int64
		//查询数据库
		//db.Model(dataList).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&dataList
		db.Model(dataList).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&dataList)

		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询到数据",
				"code": 200,
				"data": gin.H{
					"list":     dataList,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})
		}
	})

	//端口号
	PORT := "8989"
	r.Run(":" + PORT)
}
