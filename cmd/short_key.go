package cmd

import (
	"fit/config"
	"fit/models"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

var shortIds []int
var m sync.Mutex
var size = 1000

func init() {
	m.Lock()
	defer m.Unlock()
	readDb := config.GetReadDB()
	var shortIdObjs []models.ShortId
	if err := readDb.Find(&shortIdObjs).Error; err != nil || len(shortIdObjs) == 0 {
		// 데이터가 없으면 생성
		var shortUrlObj models.ShortUrl
		// 마지막 shortUrl 조회
		if err := readDb.Last(&shortUrlObj).Error; err != nil || shortUrlObj.ID == 0 {
			ids := genShortIds(0)
			shortIds = append(shortIds, ids...)
		} else {
			lastId := shortUrlObj.ID
			ids := genShortIds(int(lastId))
			shortIds = append(shortIds, ids...)
		}
	} else {
		// 데이터가 있으면 메모리에 넣음
		var ids []int
		for _, shortIdObj := range shortIdObjs {
			ids = append(ids, int(shortIdObj.ID))
		}
		ids = shuffleIds(ids)
		shortIds = append(shortIds, ids...)
		// 기존 데이터 제거
		writeDb := config.GetWriteDB()
		if err := writeDb.Delete(&shortIdObjs).Error; err != nil {
			log.Fatal(err)
		}
	}
}

// Teardown 서버 종료 시 호출
func Teardown() error {
	m.Lock()
	defer m.Unlock()
	writeDb := config.GetWriteDB()
	var shortIdObjs []models.ShortId
	for _, shortId := range shortIds {
		shortIdObjs = append(shortIdObjs, models.ShortId{ID: uint(shortId)})
	}
	if err := writeDb.Create(shortIdObjs).Error; err != nil {
		return err
	}
	return nil
}

// 배열에 인자를 섞어주기
func shuffleIds(ids []int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(ids) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		ids[i], ids[j] = ids[j], ids[i]
	}
	return ids
}

func IdToKey(n int64) (string, error) {
	return idToKey(n)
}

// idToKey 인덱스를 이용한 키 생성
func idToKey(n int64) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("유효하지 않은 입력입니다.")
	}
	var total int64 = 0
	length := 4
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for {
		combinations := int64(math.Pow(float64(len(characters)), float64(length)))
		if total+combinations >= n {
			break
		}
		total += combinations
		length++
	}
	n -= total + 1
	result := ""
	for length > 0 {
		result = string(characters[n%int64(len(characters))]) + result
		n /= int64(len(characters))
		length--
	}
	return result, nil
}

// genShortIds 10000개의 숫자를 랜덤으로 섞어 shortIds 에 넣음
func genShortIds(currId int) []int {
	var startIdx int
	if currId%size == 0 {
		startIdx = currId
	} else {
		startIdx = ((currId / size) + 1) * size
	}
	startIdx++
	ids := make([]int, size)
	for i := 0; i < size; i++ {
		ids[i] = startIdx + i
	}
	return shuffleIds(ids)
}

// GetShortKey shortIds 에서 하나를 꺼내서 반환
func GetShortKey() (string, error) {
	m.Lock()
	defer m.Unlock()
	shortId := shortIds[0]
	shortIds = shortIds[1:]
	if len(shortIds) == 0 {
		ids := genShortIds(shortId)
		shortIds = append(shortIds, ids...)
	}
	key, err := idToKey(int64(shortId))
	if err != nil {
		return "", err
	}
	return key, nil
}
