package models

import "sync"

type IPRateLimiter struct{
	limiters map[string]*RateLimiter
	mutex sync.Mutex
}

func NewIpLimiter() *IPRateLimiter{ 
	return &IPRateLimiter{
		limiters: make(map[string]*RateLimiter),
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *RateLimiter{ 
	i.mutex.Lock()

	defer i.mutex.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists { 
		limiter = NewRateLimiter(10, 0.15)
		i.limiters[ip] = limiter
	}

	return limiter
}