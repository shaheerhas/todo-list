package tasks

func findInSlice(s []string, key string) bool {
	for i := range s {
		if s[i] == key {
			return true
		}
	}
	return false
}

//func CacheCheckMiddleware(c *gin.Context) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// Get the ignoreCache parameter
//		ignoreCache := strings.ToLower(c.Query("ignoreCache")) == "true"
//
//		userID := strconv.Itoa(int(userId))
//		// See if we have a cached response
//
//		response, exists := cachey.Get(userID)
//		if !ignoreCache && exists && response != nil {
//			// If so, use it
//			log.Println("response", response)
//			userCache, exists2 := response.(UserCache)
//			if exists2 && userCache.requestStrings != nil {
//				reqBytes, responseExists := userCache.requestStringMap[c.Request.RequestURI]
//				if responseExists && reqBytes != nil {
//
//					log.Println("cache exists")
//					log.Println(reqBytes)
//					c.Data(http.StatusOK, "application/json", reqBytes)
//					c.Abort()
//					return
//				} else {
//					log.Println("before the storm")
//
//					userCache.requestStrings = append(userCache.requestStrings, c.Request.RequestURI)
//					c.Writer = &userCache
//					c.Next()
//				}
//			} else {
//				log.Println("else")
//			}
//
//		} else {
//			// If not, pass our cache writer to the next middleware
//			reqStrings := []string{c.Request.RequestURI}
//			userCache := UserCache{requestStrings: reqStrings, requestStringMap: map[string][]byte{}}
//			cachey.Set(userID, userCache, cache.DefaultExpiration)
//			bcw := &userCache
//			c.Writer = bcw
//			log.Println(cachey.Get(userID))
//			c.Next()
//		}
//
//	}
//}
