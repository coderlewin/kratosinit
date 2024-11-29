package main

import (
	"github.com/coderlewin/kratosinit/internal/data/gorm_gen/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

const _dsn = "root:12345678@tcp(127.0.0.1:3306)/gininit?charset=utf8mb4&parseTime=True"

func main() {
	autoGen()
}

func autoGen() {
	db, err := gorm.Open(mysql.Open(_dsn))
	if err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		ModelPkgPath: "internal/data/gorm_gen/entity",
		OutPath:      "internal/data/gorm_gen/dal",
		Mode:         gen.WithQueryInterface | gen.WithDefaultQuery,
	})

	g.UseDB(db)

	g.ApplyInterface(func() {}, entity.User{})

	g.GenerateAllTable(logicDeleteConfig()...)

	g.Execute()
}

// logicDeleteConfig 逻辑删除
func logicDeleteConfig() []gen.ModelOpt {
	return []gen.ModelOpt{
		gen.FieldGORMTag("is_delete", func(tag field.GormTag) field.GormTag {
			tag.Append("softDelete", "flag")
			return tag
		}),
		gen.FieldType("is_delete", "soft_delete.DeletedAt"),
	}
}
