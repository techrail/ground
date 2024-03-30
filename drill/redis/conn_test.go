package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/techrail/ground/cache"
)

func TestRedisConnection(t *testing.T) {
	c := cache.RedisConfig{}
	c.Enabled = true
	c.Url = ""
	c.OperationMode = cache.ModeAuto

	r := cache.CreateNewRedisClient(c)

	fmt.Println("-------Setting string--------")
	r.SetStringWithExpiry("name", "Techrail", 3*time.Second)
	fmt.Printf("-------Getting string-------")
	fmt.Println(r.Get("name"))

	var redislist [5]string
	redislist[0] = "Golang"
	redislist[1] = "Java"
	redislist[2] = "Ruby"
	redislist[3] = "Typescript"
	redislist[4] = "Python"

	redismap := make(map[string]string)
	redismap["OS"] = "Ubuntu"
	redismap["Container"] = "Docker"
	redismap["Orchestration"] = "K8S"
	redismap["Cloud"] = "AWS"
	redismap["DB"] = "PGSql"

	fmt.Println("-------Setting list contents-------")
	for i := 0; i < 5; i++ {
		r.SetListContents("TestList", redislist[i])
	}
	fmt.Println("-------Getting list contents-------")
	fmt.Println(r.GetListRange("TestList", 0, 4))

	fmt.Println("-------Setting hash contents-------")
	r.SetHash("TestMap", redismap)
	fmt.Println("-------Getting hash contents-------")
	fmt.Println(r.GetHashVals("TestMap"))

	fmt.Println("-------Setting set contents-------")
	fmt.Println(r.SetAdd("TestSet", "RedisConnectionTest"))
	fmt.Println("-------Getting set contents-------")
	fmt.Println(r.GetSetMembers("TestSet"))

	fmt.Println("-------Deleting list contents-------")
	for i := 0; i < 5; i++ {
		r.DeleteListElements("TestList")
	}

	fmt.Println("-------Deleting hash contents-------")
	r.DeleteHash("TestMap", "OS")
	r.DeleteHash("TestMap", "Container")
	r.DeleteHash("TestMap", "Orchestration")
	r.DeleteHash("TestMap", "Cloud")
	r.DeleteHash("TestMap", "DB")

	fmt.Println("-------Deleting set contents-------")
	r.DeleteSet("TestSet")
}
