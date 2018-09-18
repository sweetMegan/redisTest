package main

import (
	"github.com/go-redis/redis"
	"fmt"
	"encoding/json"
)

var client *redis.Client

type Person struct {
	Name string
	Age  int
}

func main() {
	connectRedisServer()
	transaction()
}
func connectRedisServer() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", //服务器地址
		Password: "",               //密码
		DB:       0,                //使用默认DB
	})
	//查看连接是否成功
	pong, err := client.Ping().Result()
	fmt.Println(pong, "===", err)
}

/**
*字符串类型
*不仅可以将string，int等基本类型已string类型保存
*结构体，数组，切片，map都可以通过序列化的方式以字符串方式保存，但是不建议，因为redis提供list，hash类型保存这些数据
*/
//存储字符串
func saveString() {
	//expiration 过期时间
	//0:表示永不过期
	err := client.Set("name", "sw", 0).Err()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("string success")
	}
}

//取字符串
func getString() {
	val, err := client.Get("name").Result()
	fmt.Println(val, err)
}

//存储结构体
//将结构体转字符串存储
func saveStructString() {
	p := Person{Name: "zhq", Age: 10}
	pbytes, _ := json.Marshal(p)
	err := client.Set("person", pbytes, 0).Err()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("struct success")
	}
}

//取结构体
//将结构体字符串
func getStructString() {
	val, err := client.Get("person").Bytes()
	if err != nil {
		fmt.Println(err)
	} else {
		person := &Person{}
		json.Unmarshal(val, person)
		fmt.Println(person)
	}
}

//存储数组
//将数组转字符串存储
func saveArrayString() {
	array := [2]Person{Person{"zhq", 10}, Person{"sw", 8}}
	//将数组序列化
	arrayBytes, _ := json.Marshal(array)
	err := client.Set("persons", arrayBytes, 0).Err()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("array success")
	}
}

//取数组
//取数组字符串
func getArrayString() {
	arrayBytes, err := client.Get("persons").Bytes()
	if err != nil {
		fmt.Println(err)
	} else {
		array := []Person{}
		//反序列化
		json.Unmarshal(arrayBytes, &array)
		fmt.Println(array[0].Name, "==", array[1].Name)
	}
}

//存储切片
func saveSliceString() {
	array := make([]Person, 2)
	array[0] = Person{"zhq2", 10}
	array[1] = Person{"sw", 8}
	//将数组序列化
	arrayBytes, _ := json.Marshal(array)
	err := client.Set("persons2", arrayBytes, 0).Err()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("slice success")
	}
}

//取切片
func getSliceString() {
	arrayBytes, err := client.Get("persons2").Bytes()
	if err != nil {
		fmt.Println(err)
	} else {
		array := []Person{}
		//反序列化
		json.Unmarshal(arrayBytes, &array)
		fmt.Println(array[0].Name, "==", array[1].Name)
	}
}

//存储map
//map转字符串存储
func saveMap() {
	response := make(map[string]string)
	response["msg"] = "this is map"
	response["data"] = "test map"
	responseByte, _ := json.Marshal(response)
	err := client.Set("map", responseByte, 0).Err()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("map success")
	}
}

//取map
//取map字符串
func getMap() {
	mapBytes, err := client.Get("map").Bytes()
	if err != nil {
		fmt.Println(err)
	} else {
		response := make(map[string]string)
		//反序列化
		json.Unmarshal(mapBytes, &response)
		fmt.Println(response["msg"], response["data"])
	}
}

/**
*list类型
*/
func setArray() {
	client.LPush("list", 1, 15, 4)
}
func getArray() {
	vals, err := client.LRange("list", 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(vals, "==", vals[0])
	}
}
func sortArray() {
	vals, err := client.Sort("list", &redis.Sort{Offset: 0, Count: 0, Order: "ASC"}).Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("排序结果:", vals)
	}
}

/**
*hash类型
*可以理解为map
*/
func setHash() {
	person := make(map[string]interface{})
	person["name"] = "zhq"
	person["age"] = 18
	person["birthday"] = "20170225"
	err := client.HMSet("user:003", person).Err()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("hash success")
	}
}
func getHash() {
	//获取固定属性值
	res, err := client.HGet("user:003", "name").Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
	//获取所有属性值
	person, err := client.HGetAll("user:003").Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(person, "====", person["name"])
	}
}

//hash数组
//将hash存入list当中，并通过age字段排序输出
func sortHashList() {
	//person1
	person := make(map[string]interface{})
	person["id"] = "001"
	person["name"] = "zhq"
	person["age"] = 18
	person["birthday"] = "20170225"
	//保存person1
	savePerson(person)
	//person2
	person2 := make(map[string]interface{})
	person2["id"] = "002"
	person2["name"] = "sw"
	person2["age"] = 17
	person2["birthday"] = "20170401"
	//保存person2
	savePerson(person2)
	//person3
	person3 := make(map[string]interface{})
	person3["id"] = "003"
	person3["name"] = "sweet"
	person3["age"] = 23
	person3["birthday"] = "20170225"
	//保存person1
	savePerson(person3)
	//将p1,p2,p3存入list
	err := client.LPush("persons", fmt.Sprintf("user:%s", person["id"].(string)), fmt.Sprintf("user:%s", person2["id"].(string)),
		fmt.Sprintf("user:%s", person3["id"].(string))).Err()
	if err != nil {
		fmt.Println(err)
	} else {
		//遍历list
		vals, err := client.LRange("persons", 0, -1).Result()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(vals, "**", vals[0])
		}
		//list按age排序
		vals, err = client.Sort("persons", &redis.Sort{By: "*->age", Order: "asc"}).Result()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(vals, "$$$", vals[0])
		}
	}
}

//
func savePerson(person map[string]interface{}) {
	//person1保存为hash
	err := client.HMSet(fmt.Sprintf("user:%s", person["id"].(string)), person).Err()
	if err != nil {
		fmt.Println(err)
	}
}

/**
*集合 set
*set与list区别
*set无序、键值唯一不可以重复
*可以理解为value为空的Hash
*/
func setSet() {
	err := client.SAdd("key1", "a", "b", "c").Err()
	client.SAdd("key2", "c", "d", "e").Err()
	client.SAdd("key3", "a", "c", "f").Err()
	if err != nil {
		fmt.Println(err)
	} else {
		//取set中所有元素
		res, _ := client.SMembers("key3").Result()
		fmt.Println(res) //[c a f]
		//求两个set的差值
		res, _ = client.SDiff("key1", "key2").Result()
		fmt.Println(res) //[b a]
		//求两个set的差值并将结果保存在新的set中
		client.SDiffStore("key", "key1", "key2")
		res, _ = client.SMembers("key").Result()
		fmt.Println("差集", res) // [b a]
		//求多个set的并集
		res, _ = client.SInter("key1", "key2", "key3").Result()
		fmt.Println(res) //[c]
		//求多个set的并集并将结果保存在新的set中
		client.SInterStore("key", "key1", "key2", "key3")
		res, _ = client.SMembers("key").Result()
		fmt.Println("并集：", res) //[c]
		//求多个set的合集
		res, _ = client.SUnion("key1", "key2", "key3").Result()
		fmt.Println(res)
		//求两个set的合集并将结果保存在新的set中
		client.SUnionStore("key", "key1", "key2", "key3")
		res, _ = client.SMembers("key").Result()
		fmt.Println("合集", res)
		//查看set中是否包含某元素
		//exist, _ := client.SIsMember("key1", "hello").Result()
		exist, _ := client.SIsMember("key1", "a").Result()
		if exist {
			fmt.Println("有")
		} else {
			fmt.Println("没有")
		}
		//将key1中的元素移动到key2中
		success, _ := client.SMove("key1", "key2", "a").Result()
		if success {
			res, _ = client.SMembers("key1").Result()
			fmt.Println(res) //[b c]
			res, _ = client.SMembers("key2").Result()
			fmt.Println(res) //[d a e c]
		} else {
			fmt.Println("出错了")
		}
		//获取并删除set中随机的N个元素
		v, _ := client.SPopN("key2", 2).Result()
		fmt.Println(v) //[a d]
		res, _ = client.SMembers("key2").Result()
		fmt.Println(res) //[e c]
		//随机获取set中N个元素,与pop不同，SRandMemberN不删除元素
		res, _ = client.SRandMemberN("key1", 10).Result()
		fmt.Println(res) //[b c]
		//删除一个元素
		rem, _ := client.SRem("key3", "c", "d").Result()
		fmt.Println(rem) //1 即1条数据被修改 因为key3中没有d元素
		res, _ = client.SMembers("key3").Result()
		fmt.Println(res) //[a f]
	}
}

//迭代器遍历Set
func func_SSCAN() {
	vals, cursor, _ := client.SScan("key", 0, "1*", 4).Result()
	fmt.Println("cursor:", cursor) //0
	fmt.Println("vals:", vals)     //vals: [1 11 12 13 14]
}

//消息订阅
func registMessage(ch chan string) {
	go func() {
		for true {
			message, _ := client.Subscribe("sw").ReceiveMessage()
			fmt.Println(message)
			if message.Payload == "my love"  {
				ch <- "finish"
			}
		}
	}()
}
func transaction() {
	client.Watch(func(tx *redis.Tx) error {
		err := tx.Set("age",20,0).Err()
		err = tx.Set("name","sw",0).Err()
		return err
	},"age")
}
