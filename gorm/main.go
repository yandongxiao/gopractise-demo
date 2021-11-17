package main

import (
	"fmt"
	"log"

	"github.com/spf13/pflag"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 首先定义了一个 GORM 模型（Models）
// GORM 使用模型（Models）来映射一个数据库表。默认情况下，使用 ID 作为主键，
// 使用结构体名的 snake_cases 作为表名，使用字段名的 snake_case 作为列名，
// 并使用 CreatedAt、UpdatedAt、DeletedAt 字段追踪创建、更新和删除时间。
type Product struct {
	// Models 是标准的 Go struct，是数据库中的常见字段
	gorm.Model
	Code  string `gorm:"column:code"`
	Price uint   `gorm:"column:price"`
}

// TableName maps to mysql table name.
// 我们可以给 Models 添加 TableName 方法，来告诉 GORM 该 Models 映射到数据库中的哪张表。
// 如果没有指定表名，则 GORM 使用结构体名的蛇形复数作为表名。
// TODO(ydx): 例如：结构体名为 DockerInstance ，则表名为 dockerInstances。 这个举例可能是错误的。
func (p *Product) TableName() string {
	return "product"
}

var (
	host     = pflag.StringP("host", "H", "127.0.0.1:3306", "MySQL service host address")
	username = pflag.StringP("username", "u", "root", "Username for access to mysql service")
	password = pflag.StringP("password", "p", "root", "Password for access to mysql, should be used pair with password")
	database = pflag.StringP("database", "d", "test", "Database name to use")
	help     = pflag.BoolP("help", "h", false, "Print this help message")
)

func main() {
	// Parse command line flags
	pflag.CommandLine.SortFlags = false
	pflag.Usage = func() {
		pflag.PrintDefaults()
	}
	pflag.Parse()
	if *help {
		pflag.Usage()
		return
	}

	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		*username,
		*password,
		*host,
		*database,
		true,
		"Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 1. Auto migration for given models
	// 如果 Product 中有新增字段或索引，则相应地变更表结构
	db.AutoMigrate(&Product{})

	// 2. Insert the value into database
	if err := db.Create(&Product{Code: "D42", Price: 100}).Error; err != nil {
		log.Fatalf("Create error: %v", err)
	}
	PrintProducts(db)

	// 3. Find first record that match given conditions
	product := &Product{}
	if err := db.Where("code= ?", "D42").First(&product).Error; err != nil {
		log.Fatalf("Get product error: %v", err)
	}

	// 4. Update value in database, if the value doesn't have primary key, will insert it
	product.Price = 200
	if err := db.Save(product).Error; err != nil {
		log.Fatalf("Update product error: %v", err)
	}
	PrintProducts(db)

	// 5. Delete value match given conditions
	// 因为 Product 中有 gorm.DeletedAt 字段，所以，上述删除操作不会真正把记录从数据库表中删除掉，
	// 而是通过设置数据库 product 表 deletedAt 字段为当前时间的方法来删除。
	if err := db.Where("code = ?", "D42").Delete(&Product{}).Error; err != nil {
		log.Fatalf("Delete product error: %v", err)
	}
	PrintProducts(db)
}

// List products
func PrintProducts(db *gorm.DB) {
	products := make([]*Product, 0)
	var count int64
	d := db.Where("code like ?", "%D%").
		Offset(0).Limit(2).
		Order("id desc").Find(&products).
		Offset(-1).Limit(-1).
		Count(&count)
	if d.Error != nil {
		log.Fatalf("List products error: %v", d.Error)
	}

	log.Printf("totalcount: %d", count)
	for _, product := range products {
		log.Printf("\tcode: %s, price: %d\n", product.Code, product.Price)
	}
}