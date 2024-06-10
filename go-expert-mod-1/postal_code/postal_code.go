package postal_code

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Adress struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type PostalCodeController struct{}

func (c PostalCodeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postalCode := r.URL.Query().Get("cep")
	if postalCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cep is required"))
		return
	}

	adress, err := getAdress(postalCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adress)
}


func getAdress(cep string) (*Adress, error) {
	response, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return &Adress{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return &Adress{}, err
	}

	adress := Adress{}
	err = json.Unmarshal(body, &adress)
	if err != nil {
		return &Adress{}, err
	}

	return &adress, nil
}
