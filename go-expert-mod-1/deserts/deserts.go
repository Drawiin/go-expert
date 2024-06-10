package deserts

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"html/template" 
)

type Desert struct {
	Name     string
	Calories int
}

type DesertController struct{}

func (c DesertController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lenght, err := strconv.Atoi(r.URL.Query().Get("length"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("length is required"))
		return
	}
	desert := randomDesert(lenght)
	fmt.Println()

	templates := []string{"deserts/deserts.html", "deserts/deserts_list.html"}
	templateError := template.Must(template.New("deserts.html").ParseFiles(templates...)).Execute(w, desert)
	if templateError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(templateError.Error()))
		return
	}
}

func randomDesert(lenght int) []Desert {
	deserts := []Desert{
		{Name: "Brigadeiro", Calories: 100},
		{Name: "Brigadeiro Light", Calories: 50},
		{Name: "Tiramissu", Calories: 200},
		{Name: "Bolo de chocolate", Calories: 300},
		{Name: "Bolo de aveia", Calories: 250},
		{Name: "Bolo de milho", Calories: 197},
	}
	result := []Desert{}
	for i := 0; i < lenght; i++ {
		result = append(result, deserts[rand.Int()%len(deserts)])
	}

	return result
}
