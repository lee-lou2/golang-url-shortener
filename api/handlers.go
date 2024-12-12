package api

import (
	"encoding/json"
	"fit/cmd"
	"fit/config"
	"fit/models"
	"fit/pkg/utils"
	"log"
	"net/http"
	"regexp"
	"time"
)

// createShortUrlHandler 단축 URL 생성 핸들러
func createShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Url string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 유효성 검사
	regexPattern := `^(?i)(https?:\/\/)?(www\.)?([a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5})(:[0-9]{1,5})?(\/.*)?(\?.*)?$`
	re := regexp.MustCompile(regexPattern)
	if !re.MatchString(req.Url) {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	// 데이터 존재 여부 확인
	hashedUrl := utils.SHA256(req.Url)
	readDb := config.GetReadDB()
	var shortUrlObj models.ShortUrl
	if err := readDb.Select("short_key").Where("hashed_url = ?", hashedUrl).Take(&shortUrlObj).Error; err == nil {
		_, _ = w.Write([]byte(shortUrlObj.ShortKey))
		return
	}

	// 데이터 생성
	writeDb := config.GetWriteDB()
	shortKey, err := cmd.GetShortKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := writeDb.Create(&models.ShortUrl{
		ShortKey:  shortKey,
		Url:       req.Url,
		HashedUrl: hashedUrl,
	}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("short_key: %s, url: %s", shortKey, req.Url)

	// 캐시에 저장
	cache := config.GetCache()
	cache.Set(shortKey, req.Url, 1*time.Hour)
	_, _ = w.Write([]byte(shortKey))
}

// redirectShortUrlHandler 단축 URL 리다이렉트 핸들러
func redirectShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	// 파라미터 조회
	shortKey := r.PathValue("short_key")
	if shortKey == "" {
		http.Error(w, "short_key is required", http.StatusBadRequest)
		return
	}

	// 1. 캐시 조회
	cache := config.GetCache()
	url := cache.Get(shortKey)
	if url != "" {
		log.Printf("cache hit: %s/%s", "https://f-it.kr", shortKey)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	// 2. 데이터베이스 조회
	readDb := config.GetReadDB()
	var shortUrlObj models.ShortUrl
	if err := readDb.Select("url").Where("short_key = ?", shortKey).Take(&shortUrlObj).Error; err != nil || shortUrlObj.Url == "" {
		http.Error(w, "short_key is not found", http.StatusNotFound)
		return
	}

	// 캐시 저장
	cache.Set(shortKey, shortUrlObj.Url, 1*time.Hour)
	log.Printf("cache miss: %s/%s", "https://f-it.kr", shortKey)
	http.Redirect(w, r, shortUrlObj.Url, http.StatusFound)
}
