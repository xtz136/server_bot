package commands

// func LoginUser(ctx Context) {
// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 	tokenList := viper.GetString("systems." + ctx.Name + ".list_user_system")
// 	tokenMake := viper.GetString("systems." + ctx.Name + ".make_token")

// 	var accountAndWhy []string
// 	for {
// 		ctx.Sender <- "请回复登录账号，同时回复申请原因(用空格隔开)。例如：\n\nxiaoming 火箭升空演示，需要测试推进系统"
// 		accountAndWhy = strings.SplitN(<-ctx.Reply, " ", 2)
// 		if len(accountAndWhy) != 2 {
// 			ctx.Sender <- "回复有误，请重新输入"
// 		}
// 		break
// 	}

// 	account := accountAndWhy[0]
// 	why := accountAndWhy[1]

// 	if account == "" || why == "" {
// 		ctx.MakeTalkEnd(ctx.Sender, "回复有误，本次服务结束")
// 		return
// 	}

// 	// 找到这个用户绑定的所有系统
// 	resp, err := http.Get(tokenList + "&account=" + account)
// 	if err != nil {
// 		ctx.Log.Error().Err(err).Str("url", tokenList).Msg("list user system")
// 		ctx.MakeTalkEnd(ctx.Sender, "")
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// 让客服选一个
// 	userSystem := UserSystem{}
// 	json.NewDecoder(resp.Body).Decode(&userSystem)

// 	// 不需要处理，结束
// 	if len(userSystem.Data) == 0 {
// 		ctx.MakeTalkEnd(ctx.Sender, "这个用户没有绑定系统，本次服务结束")
// 		return
// 	}

// 	// 让客服选择要处理的IP
// 	query := "请选择，回复序号，格式[用户名][系统名][部门名]：\n\n"
// 	for i, item := range userSystem.Data {
// 		query += fmt.Sprintf("%d. [%s][%s][%s]\n", i, item.UserAlias, item.SystemAlias, item.DanweiParentAlias)
// 	}
// 	ctx.Sender <- query
// 	var index int

// 	for {
// 		answer := <-ctx.Reply
// 		i, err := strconv.Atoi(answer)
// 		if err != nil || i >= len(userSystem.Data) {
// 			ctx.Sender <- "选择错误，请重新选择!"
// 		} else {
// 			index = i
// 			break
// 		}
// 	}

// 	// 请求token
// 	{
// 		params := url.Values{
// 			"account":     {account},
// 			"system_id":   {userSystem.Data[index].SystemID},
// 			"danwei_id":   {userSystem.Data[index].DanweiID},
// 			"quanxian_id": {userSystem.Data[index].QuanxianID},
// 			"why":         {why},
// 		}
// 		resp, err := http.Get(tokenMake + "&" + params.Encode())
// 		if err != nil {
// 			ctx.Log.Error().Err(err).Str("url", tokenMake).Msg("make token")
// 			ctx.MakeTalkEnd(ctx.Sender, "")
// 			return
// 		}
// 		token := Token{}
// 		json.NewDecoder(resp.Body).Decode(&token)
// 		ctx.MakeTalkEnd(ctx.Sender, fmt.Sprintf("访问下方链接(%s)即可登录,本次服务结束\n%s", token.Message, token.Data))
// 		defer resp.Body.Close()

// 		ctx.Log.Info().
// 			Str("account", params.Get("account")).
// 			Str("danwei_id", params.Get("danwei_id")).
// 			Str("quanxian_id", params.Get("quanxian_id")).
// 			Str("why", params.Get("why")).
// 			Msg("make token")
// 	}
// }
