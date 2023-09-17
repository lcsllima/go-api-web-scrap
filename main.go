package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/PuerkitoBio/goquery"
)

type scrapedcontent struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Content    string  `json:"content"`
	RecordTime float64 `json:"record_time"`
}

// var scrapedcontents = []scrapedcontent{
// 	{ID: "1", Title: "grasshoper", Content: "grasshoper content", RecordTime: 56.99},
// 	{ID: "2", Title: "ant", Content: "ant content", RecordTime: 17.99},
// 	{ID: "3", Title: "bee", Content: "bee content", RecordTime: 39.99},
// }

var scrapedcontents = []scrapedcontent{
	{ID: "1", Title: "grasshoper", Content: "", RecordTime: 0},
	{ID: "2", Title: "ant", Content: "", RecordTime: 0},
	{ID: "3", Title: "bee", Content: "", RecordTime: 0},
}

func scrapeWikipedia(title string) (string, error) {
	url := "https://pt.wikipedia.org/wiki/" + title
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	// Obtém o conteúdo (html mesmo) do id bodyContent
	pageContent := doc.Find("#bodyContent").Text()
	return pageContent, nil
}

// Atualiza os dados dos insetos
func atualizarDados(c *gin.Context) {
	insetos := []string{"Gafanhoto", "Formiga", "Abelha"}

	for i, inseto := range insetos {
		content, err := scrapeWikipedia(inseto)
		if err != nil {
			log.Printf("Erro ao obter dados para %s: %v", inseto, err)
			continue
		}

		// Atualiza os dados em scrapedcontents
		scrapedcontents[i].Title = inseto
		scrapedcontents[i].Content = content
		scrapedcontents[i].RecordTime = float64(time.Now().UnixNano()) / 1e9
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dados atualizados com sucesso"})
}

func getScrapedContent(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, scrapedcontents)
}


func main() {
	router := gin.Default()

	router.GET("/insect-content", getScrapedContent)
	router.GET("/atualizar-dados", atualizarDados)
	router.Run("localhost:8080")
}
