package config

// "gitee.com/dark.H/go-remote-repl/cli"
// "gitee.com/dark.H/go-remote-repl/datas"
// "gitee.com/dark.H/go-remote-repl/lpc"
// "gitee.com/dark.H/go-remote-repl/remote"

const (
	RandomeMode = 2
	FlowMode    = 0
	ManueMode   = 1
)

var (
	// RouteSys   = &RouteSystem{}
	UsedRoutes = ""
)

// func SetHandle() {
// 	remote.RegistedService("4. show config", func(m *remote.ApiServerMan) {
// 		m.SendMsg(UsedRoutes)
// 	})
// 	remote.RegistedService("use config file", func(m *remote.ApiServerMan) {
// 		// m.FunctionInput()
// 		buf, err := m.FunctionInput("use a config", datas.G{})
// 		fmt.Println(buf)
// 		if err != nil {
// 			log.Println(err)
// 			m.SendMsg(err.Error())
// 		}

// 		UsedRoutes = buf
// 	})
// 	remote.RegistedClientHandler("use a config", func(a *remote.ApiConn, args string) {
// 		file := lpc.FilePathInput()
// 		var err error
// 		if file, err = filepath.Abs(file); err != nil {
// 			return
// 		}
// 		if buf, err := ioutil.ReadFile(file); err == nil {
// 			a.SendMsg(string(buf))
// 		}
// 		a.Close()
// 	})
// 	// fmt.Println(remote.RegistedLocalFunc)
// }

//
// route : kd["forward"] = from -> to /manue mode
// route : kd["connected"] = server -> config /random/flow
//
//

// func GenerateRoute(buf string) string {
// 	lines := strings.Split(strings.TrimSpace(buf), "\n")
// 	length := len(lines)
// 	if length < 3 {
// 		return ""
// 	}
// 	count, err := strconv.Atoi(lines[0])
// 	if err != nil {
// 		return ""
// 	}
// 	lines = lines[1:]
// 	length--
// 	start := "proxy://" + lines[0]
// 	// count := lines[0]
// 	choosed := []string{}

// 	chooseNum := 0
// 	used := map[string]int{}
// 	end := ""
// 	for n, i := range lines {
// 		if strings.TrimSpace(i) == "" {
// 			continue
// 		}
// 		if strings.HasPrefix(i, "[") {
// 			start = "proxy://" + i[1:]
// 			used[i] = n
// 		} else if strings.HasPrefix(i, "]") {
// 			end = "proxy://" + i[1:]
// 			used[i] = n
// 		}
// 	}

// 	for _, b := range used {
// 		lines = append(lines[:b], lines[b+1:]...)
// 		length--
// 	}

// 	for {
// 		if chooseNum >= count {
// 			break
// 		}
// 		i := rand.Int() % len(lines)
// 		l := lines[i]
// 		if _, ok := used[l]; ok {
// 			continue
// 		}
// 		choosed = append(choosed, "proxy://"+l)
// 		chooseNum++
// 		used[l] = i
// 	}
// 	if end != "" {
// 		choosed = append(choosed, end)
// 	}
// 	return strings.Join(append([]string{start}, choosed...), CHAIN)
// }

// func StartAPIServer() {
// 	if !lpc.IfOuterStart() {
// 		go cli.StartOuterAPIServer()
// 		// os.Exit(0)
// 	} else {
// 		fmt.Println("first ... kill old !!!!")
// 		lpc.KillAPI()
// 		time.Sleep(1 * time.Second)
// 		go cli.StartOuterAPIServer()
// 	}
// 	if !lpc.IfLocalStart() {
// 		ln, err := lpc.UnixListener()
// 		if err != nil {
// 			log.Fatal("unix start err:", err)
// 		}
// 		cli.RunningAPIListener(ln)
// 	}
// }

// type RouteSystem struct {
// 	LastID    int
// 	RouteMode int
// }

// func GetRoute(key string) (c datas.Gi) {
// 	switch RouteSys.RouteMode {
// 	case FlowMode:
// 		c = RouteSys.FlowRoute()
// 	case RandomeMode:
// 		c = RouteSys.RandomeRoute()
// 	case ManueMode:
// 		c = RouteSys.ManuRoute(key)
// 	}
// 	return
// }

// func (rsys *RouteSystem) RandomeRoute() (config datas.Gi) {
// 	routes, routesCount := GetConnected()
// 	if routesCount < 0 {
// 		return nil
// 	}
// 	if rand.Float32() > 0.5 {
// 		return nil
// 	}
// 	one := rand.Int() % routesCount
// 	c := 0
// 	for _, serverConf := range routes {
// 		if c == one {
// 			config = serverConf.(datas.Gi)
// 			// key = serverHost
// 			break
// 		}
// 		c++
// 	}
// 	return
// }

// func (rsys *RouteSystem) FlowRoute() (config datas.Gi) {
// 	routes, routesCount := GetConnected()
// 	if routesCount <= 0 {
// 		return nil
// 	}
// 	one := (rsys.LastID + 1) % routesCount
// 	c := 0
// 	for _, serverConf := range routes {
// 		if c == one {
// 			config = serverConf.(datas.Gi)
// 			// key = serverHost
// 			break
// 		}
// 		c++
// 	}
// 	rsys.LastID = (rsys.LastID + 1) % routesCount
// 	return
// }

// func (rsys *RouteSystem) ManuRoute(key string) (config datas.Gi) {
// 	config = GetMapConf(key)
// 	return
// }

// func GetServers() (servers []string) {
// 	for k, v := range datas.MemDB.Kd {
// 		if strings.HasPrefix(k, "c:") {
// 			servers = append(servers, v["server"].(string))
// 		}
// 	}
// 	return
// }

// func GetCanConnectedServer() (server []string) {
// 	for k, v := range datas.MemDB.Kd {
// 		if addr, ok := v["server"]; ok {
// 			if IsMe(addr.(string)) {
// 				continue
// 			}
// 			server = append(server, k)
// 		}
// 	}
// 	return
// }

// // get connected server in db
// func GetConnected() (routes datas.Gi, c int) {
// 	routes = make(datas.Gi)
// 	for k, v := range datas.MemDB.Kd {
// 		if strings.HasPrefix(k, "c:") {
// 			routes[k[2:]] = v
// 		}
// 	}
// 	c = len(routes)
// 	return
// }

// // add server conf to db
// func AddConnected(key string, conf datas.Gi) {
// 	datas.MemDB.Kd["c:"+key] = conf
// }

// // for manu used
// func AddMapConf(key string, conf datas.Gi) {
// 	datas.MemDB.Kd["m:"+key] = conf
// }

// func GetMapConf(key string) (conf datas.Gi) {
// 	conf, _ = datas.MemDB.Kd["m:"+key]
// 	return
// }

// func DelConf(key string) {
// 	datas.Locker.Lock()
// 	defer datas.Locker.Unlock()
// 	delete(datas.MemDB.Kd, "m:"+key)
// 	delete(datas.MemDB.Kd, "c:"+key)
// }

// func DeleteConnected() {
// 	datas.Locker.Lock()
// 	defer datas.Locker.Unlock()
// 	deletedKeys := []string{}
// 	for k := range datas.MemDB.Kd {
// 		if strings.HasPrefix(k, "c:") || strings.HasPrefix(k, "m:") {
// 			deletedKeys = append(deletedKeys, k)
// 		}
// 	}
// 	for _, k := range deletedKeys {
// 		delete(datas.MemDB.Kd, k)
// 	}
// }

/*

Tls 生成器
Conf 链接机制

*/
