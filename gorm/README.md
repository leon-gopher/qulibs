# Gorm 

基于 `github.com/jinzhu/gorm` 封装 *MySQL* 客户端，主要提供多实例管理功能，并集成 Trace 和 Metrics 服务。

## 快速开始

### 使用单实例

- 配置

	```yaml
	model:
	  driver: mysql
	  dsn: root:@tcp(localhost:3306)/test_db?charset=utf8&parseTime=True&loc=Local
	  max_open_conns: 512
	  max_idle_conns: 96
	  max_life_conns: 300 # second
	```

- 示例

	```golang
	package main

	import (
		"github.com//qulibs"
		mysql "github.com/leon-yc/ggs/internal/qulibs/gorm"
		"github.com/jinzhu/gorm"
	)

	var (
		MySQL *mysql.Client
	)

	func init() {
		config := &mysql.Config{
			Driver: "mysql",
			DSN:    "root:@tcp(localhost:3306)/test_db?charset=utf8&parseTime=True&loc=Local",
		}
		config.FillWithDefaults()

		logger := qulibs.NewLogger(qulibs.LogDebug)

		client, err := mysql.New(config, logger)
		if err != nil {
			panic(err.Error())
		}

		// migration tables
		if err := client.AutoMigrate(new(TestingModel)).Error; err != nil {
			panic(err.Error())
		}

		MySQL = client
	}

	// A TestingModel of gorm
	type TestingModel struct {
		*gorm.Model

		Subject string
	}

	func main() {
		logger := qulibs.NewLogger(qulibs.LogDebug)

		// create a record
		tm := &TestingModel{
			Subject: "Testing Model",
		}
		if err := MySQL.Save(tm).Error; err != nil {
			logger.Errorf("MySQL.Save(%#v): %v", tm, err)
			return
		}

		// query record from MySQL
		var tmpm TestingModel
		if err := MySQL.Where("id", tm.ID).First(&tmpm).Error; err != nil {
			logger.Errorf("MySQL.First(%T): %v", tmpm, err)
			return
		}

		logger.Infof("Retrieved record: %#v", tmpm)
	}
	```

### 使用多实例

- 配置

	```yaml
	model:
	  first:
	    driver: mysql
	    dsn: root:@tcp(localhost:3306)/test_db1?charset=utf8&parseTime=True&loc=Local
	    max_open_conns: 512
	    max_idle_conns: 96
	    max_life_conns: 300 # second
	  second:
	    driver: mysql
	    dsn: root:@tcp(localhost:3306)/test_db2?charset=utf8&parseTime=True&loc=Local
	    max_open_conns: 512
	    max_idle_conns: 96
	    max_life_conns: 300 # second
	```


- 示例

	```golang
	package main

	import (
		"github.com/leon-yc/ggs/internal/qulibs"
		mysql "github.com/leon-yc/ggs/internal/qulibs/gorm"
		"github.com/jinzhu/gorm"
	)

	// MySQL clients manager
	var (
		MySQLMgr *mysql.Manager
	)

	func init() {
		config := &mysql.Config{
			Driver: "mysql",
			DSN:    "root:@tcp(localhost:3306)/testdb?charset=utf8&parseTime=True&loc=Local",
		}

		mgrconfig := &mysql.ManagerConfig{
			"qtt": config,
		}

		MySQLMgr = mysql.NewManager(mgrconfig)


		client, err := MySQLMgr.NewClient(name, logger)
		if err != nil {
			panic(err.Error())
		}

		// migration tables
		if err := client.AutoMigrate(new(TestingModel)).Error; err != nil {
			panic(err.Error())
		}
	}

	// A TestingModel of gorm
	type TestingModel struct {
		*gorm.Model

		Subject string
	}

	func main() {
		name := "qtt"
		logger := qulibs.NewLogger(qulibs.LogDebug)

		client, err := MySQLMgr.NewClient(name, logger)
		if err != nil {
			logger.Errorf("MySQLMgr.NewClient(%s, ?): %v", name, err)
			return
		}

		// create a record
		tm := &TestingModel{
			Subject: "Testing Model",
		}
		if err := client.Save(tm).Error; err != nil {
			logger.Errorf("client.Save(%#v): %v", tm, err)
			return
		}

		// query record from MySQL
		var tmpm TestingModel
		if err := client.Where("id", tm.ID).First(&tmpm).Error; err != nil {
			logger.Errorf("client.First(%T): %v", tmpm, err)
			return
		}

		logger.Infof("Retrieved record: %#v", tmpm)
	}
	```