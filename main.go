package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func generateFibonacci(limit int) []int {
	fib := []int{0, 1}
	for fib[len(fib)-1] < limit {
		fib = append(fib, fib[len(fib)-1]+fib[len(fib)-2])
	}
	return fib
}

func findFibonacciSum(num int, fib []int) []int {
	result := []int{}
	lastUsedIndex := len(fib)

	if num == 1 {
        return []int{1}  // 1에 대해 인덱스 2를 반환 (두 번째 1)
    }

	for i := len(fib) - 1; i > 1 && num > 0; i-- {
        if fib[i] <= num && i != lastUsedIndex-1 {
            result = append([]int{i}, result...) // 결과 배열의 앞에 추가
            num -= fib[i]
            lastUsedIndex = i
        }
    }

    // 남은 1들을 처리합니다
    for num > 0 {
        result = append([]int{1}, result...) // 결과 배열의 앞에 추가
        num--
    }

    return result
}

func encryptNumber(num int, fib []int) string {
	indices := findFibonacciSum(num, fib)
	var encrypted []string
	for _, idx := range indices {
		if(idx==02) { 
			encrypted = append(encrypted, fmt.Sprintf("%02d", 1))
			continue
		 }
		encrypted = append(encrypted, fmt.Sprintf("%02d", idx))
	}
	return strings.Join(encrypted, "")
}

func encryptWordPositions(text, targetWord string) ([]string, error) {
    log.Printf("텍스트: %s", text)
    log.Printf("찾을 단어: %s", targetWord)
    if targetWord == "" {
        return nil, fmt.Errorf("대상 단어를 입력해주세요")
    }
   
    var results []string
    targetRunes := []rune(strings.ReplaceAll(targetWord, " ", "")) // 대상 단어에서 공백 제거
    pages := strings.Split(text, "___")
    fib := generateFibonacci(10000)
    
    targetIndex := 0
    for pageNum, page := range pages {
        log.Printf("페이지 %d 처리 중", pageNum+1)
        lines := strings.Split(page, "\n")
        for lineNum, line := range lines {
            log.Printf("줄 %d 처리 중: %s", lineNum+1, line)
            runes := []rune(line)
            for charNum, char := range runes {

                if char == targetRunes[targetIndex] {
                    encryptedPage := encryptNumber(pageNum+1, fib)
                    encryptedLine := encryptNumber(lineNum+1-pageNum, fib)
                    encryptedChar := encryptNumber(charNum+1, fib)
					log.Printf("%s-%s-%s",pageNum, lineNum,charNum)
                    result := fmt.Sprintf("%d번째 글자 '%c': %s-%s-%s", targetIndex+1, char, encryptedPage, encryptedLine, encryptedChar)
                    log.Printf("결과: %s", result)
                    results = append(results, result)
                    
                    targetIndex++
                    if targetIndex == len(targetRunes) {
                        return results, nil // 모든 글자를 찾았으면 종료
                    }
                    break // 다음 줄로 이동
                }
            }
            if targetIndex == len(targetRunes) {
                break // 모든 글자를 찾았으면 페이지 루프 종료
            }
        }
        if targetIndex == len(targetRunes) {
            break // 모든 글자를 찾았으면 전체 루프 종료
        }
    }
   
    if len(results) == 0 {
        return nil, fmt.Errorf("단어의 글자들을 찾을 수 없습니다")
    }
   
    return results, nil
}
func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.POST("/encrypt", func(c *gin.Context) {
		text := c.PostForm("text")
		targetWord := c.PostForm("targetWord")

		log.Printf("POST 요청 받음 - 텍스트: %s, 대상 단어: %s", text, targetWord)

		if targetWord == "" {
			c.JSON(http.StatusOK, gin.H{"error": "대상 단어를 입력해주세요"})
			return
		}

		encrypted, err := encryptWordPositions(text, targetWord)
		if err != nil {
			log.Printf("오류 발생: %v", err)
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		log.Printf("암호화 성공: %v", encrypted)
		c.JSON(http.StatusOK, gin.H{"result": encrypted})
	})

	log.Println("서버 시작 - http://localhost:8080")
	r.Run(":8080")
}