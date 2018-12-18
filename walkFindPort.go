package main

import (
	"fmt"
	"os"
	"time"
	"flag"
	"strings"
	"strconv"
	"regexp"
	"github.com/soniah/gosnmp"
	"github.com/go-redis/redis"
)

var client *redis.Client = redis.NewClient(&redis.Options{
	Addr:	"localhost:6379",
	Password:	"",
	DB:	0,
})
//refer to src/net/parse.go
const big = 0xFFFFFF
func xtoi (s string) (n int,i int,ok bool){
	n = 0
	for i = 0;i < len(s);i++{
		if '0' <= s[i] && s[i] <= '9'{
			n *= 16
			n += int(s[i] - '0')
		} else if 'a' <= s[i] && s[i] <= 'f' {
			n *= 16
			n += int(s[i]-'a') + 10
		} else if 'A' <= s[i] && s[i] <= 'F' {
			n *= 16
			n += int(s[i]-'A') + 10
		} else {
			break
		}

		if n >= big {
			return 0,i,false
		}
	}

	if i == 0 {
		return 0,i,false

	}
	return n,i,true
}

func main() {
	flag.Usage = func(){
		fmt.Printf("Usage:\n")
		fmt.Printf("    mac  - the host mac address\n")
	}
	 
	flag.Parse()
		  
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	target := flag.Args()[0]
	dec_mac := make([]string,1,6)
	
	re := regexp.MustCompile(`[-:]`)
	for _,item := range re.Split(target,-1){
		n,_,_ := xtoi(item)
		n_str := strconv.Itoa(n)
		dec_mac = append(dec_mac,n_str)	
	}

	aim := strings.TrimLeft(strings.Join(dec_mac,"."),".")

	val,err := client.Get(aim).Result()
	if err == redis.Nil {
		fmt.Println("The result does not currently exist!  \nWait....")
		sws := [3]string{"192.168.1.1","192.168.1.2","192.168.1.3"}
		for _,sw := range sws {
			connectTarget(sw)
		}
		val,err := client.Get(aim).Result()
		if err == redis.Nil{
			fmt.Println("Sorry,I did my best!")
		} else {
			fmt.Printf("%s\n",val)
		}
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf("%s\n",val)
	}
}

func connectTarget(target string) {
	gosnmp.Default.Target = target
	gosnmp.Default.Community = "public"
	gosnmp.Default.Timeout = time.Duration(10 * time.Second)
	err := gosnmp.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		os.Exit(1)
	}

	defer gosnmp.Default.Conn.Close()

	oid := "1.3.6.1.2.1.17.7.1.2.2.1.2"   //switch model is H3C S5120-52P-LI
	var result []gosnmp.SnmpPDU
	result,err = gosnmp.Default.BulkWalkAll(oid)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}

	for _,pdu := range result {
		if pdu.Value != 48 {  //port 48 is trunk port,link to other switch
			trimPre := strings.Split(pdu.Name,".")[15:]
			ret := strings.Join(trimPre,".")
			key := target + " " + strconv.Itoa(pdu.Value.(int))
			err := client.Set(ret,key,time.Hour).Err() //preserve key for one hour
			if err != nil {
				panic(err)
			}
		}
	}
}
