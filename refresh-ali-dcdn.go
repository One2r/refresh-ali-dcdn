package main

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dcdn20180115 "github.com/alibabacloud-go/dcdn-20180115/client"
	cli "github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
	"time"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dcdn20180115.Client, _err error) {
	endpoint := "dcdn.aliyuncs.com"
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Endpoint:        &endpoint,
	}
	_result = &dcdn20180115.Client{}
	_result, _err = dcdn20180115.NewClient(config)
	return _result, _err
}

func main() {
	app := &cli.App{
		Usage:       "阿里云 DCDN 缓存刷新工具",
		Description: "阿里云 DCDN 缓存刷新工具",
		UsageText:   "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "accessKeyId",
				Aliases: []string{"id"},
				Usage:   "您的AccessKey ID",
			},
			&cli.StringFlag{
				Name:    "accessKeySecret",
				Aliases: []string{"secret"},
				Usage:   "您的AccessKey Secret",
			},
			&cli.StringFlag{
				Name:    "objectPath",
				Aliases: []string{"path"},
				Usage:   "刷新目录，e.g. http://www.a.com/,http://www.b.com",
			},
		},
		Action: func(c *cli.Context) error {
			if c.String("accessKeyId") == "" {
				return cli.Exit("Err: 请指定 accessKeyId 值", 1)
			}
			if c.String("accessKeySecret") == "" {
				return cli.Exit("Err: 请指定 accessKeySecret 值", 1)
			}
			if c.String("objectPath") == "" {
				return cli.Exit("Err: 请指定 objectPath 值", 1)
			}

			accessKeyId := c.String("accessKeyId")
			accessKeySecret := c.String("accessKeySecret")
			objectPath := strings.Replace(c.String("objectPath"), ",", "\r\n", -1)
			ObjectType := "Directory"

			client, _err := CreateClient(&accessKeyId, &accessKeySecret)

			if _err != nil {
				return cli.Exit("Err: "+_err.Error(), 1)
			}

			fmt.Println("")
			fmt.Println("创建刷新任务 ...")

			refreshDcdnObjectCachesRequest := &dcdn20180115.RefreshDcdnObjectCachesRequest{
				ObjectPath: &objectPath,
				ObjectType: &ObjectType,
			}

			task, _err := client.RefreshDcdnObjectCaches(refreshDcdnObjectCachesRequest)
			if _err != nil {
				return cli.Exit("Err: "+_err.Error(), 1)
			}
			refreshTaskId := *task.Body.RefreshTaskId

			describeDcdnRefreshTasksRequest := &dcdn20180115.DescribeDcdnRefreshTasksRequest{
				TaskId: &refreshTaskId,
			}

			fmt.Println("创建刷新任务成功，任务 ID: " + refreshTaskId)
			fmt.Println("")

			complete := false
			fmt.Println("等待刷新结果 ...")
			for ; ; {
				time.Sleep(time.Duration(10) * time.Second)
				res, _err := client.DescribeDcdnRefreshTasks(describeDcdnRefreshTasksRequest)
				if _err != nil {
					continue
				}

				if *res.Body.TotalCount > 0 {
					tasksList := *res.Body.Tasks
					tasks := tasksList.Task
					for i := range tasks {
						fmt.Println("刷新: " + *tasks[i].ObjectPath + " " + *tasks[i].Status + " " + *tasks[i].Process)
						if *tasks[i].Status == "Complete" {
							complete = true
						} else {
							complete = false
						}
					}

				}
				fmt.Println("")
				if complete {
					fmt.Println("Succ: 刷新完成")
					break
				} else {
					fmt.Println("刷新未完成，继续等待 ...")
				}
			}
			return nil
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
