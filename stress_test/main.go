package main

import (
	"crypto/md5"
	"fmt"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend-stress-test/util"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"log"
	"math/rand"
	"sync"
	"time"
)

var j jsoniter.API

type setupJob struct {
	url string
	index int
	size int
	wg *sync.WaitGroup
	job func(url string, i, size int, wg *sync.WaitGroup)
}

type usuarioJob struct {
	url string
	index int
	wg *sync.WaitGroup
}

func init() {
	c := jsoniter.Config{
		MarshalFloatWith6Digits: true,
		SortMapKeys:             false,
		UseNumber:               false,
		DisallowUnknownFields:   true,
		OnlyTaggedField:         true,
		CaseSensitive:           true,
	}.Froze()

	j = c
}

func main() {
	size := 1000
	usuariosSimuntaneos := 10

	log.Println(time.Now())
	//setup("http://localhost:8080/api/v1", size, usuariosSimuntaneos)
	log.Println(time.Now())

	usuarioPool := make(chan usuarioJob)
	var wg sync.WaitGroup

	for i := 0; i <  usuariosSimuntaneos; i++ {
		go func(jobs <- chan usuarioJob) {
			for job := range jobs {
				usuario(job.url, fmt.Sprint(job.index), job.wg)
			}
		}(usuarioPool)
	}

	wg.Add(usuariosSimuntaneos * size)
	log.Println(time.Now())
	for i := 0; i < (usuariosSimuntaneos * size); i++ {
		usuarioPool <- usuarioJob{
			url: "http://localhost:8080/api/v1",
			wg:  &wg,
			index: i,
		}
	}

	wg.Wait()
	log.Println(time.Now())
}

func cadastrarEspecialidade(url string, index, size int, wg *sync.WaitGroup) {
	eJson, err := j.Marshal(map[string]string{
		"nome": fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String()))),
	})

	if err != nil {
		panic(err)
	}

	_, _, _ = util.DoRequest("PUT",
		fmt.Sprint(url, "/adm/especialidades/add"),
		string(eJson), "", "",
	)

	wg.Done()
}

func cadastrarHospital(url string, i, size int, wg *sync.WaitGroup) {
	// Cria um hospital
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	h := objetos.Hospital{
		HNOME:            fmt.Sprint("Hospital ", uuid.New().String()),
		HUF:              "SP",
		HCIDADE:          "Uma cidade ai",
		HCEP:             "48003174",
		HENDERECO:        fmt.Sprint("Um endereço ai,", r.Int()),
		HCOMPLEMENTO:     fmt.Sprint(r.Int()),
		HTELEFONE:        33368200,
		HISPRONTOSOCORRO: r.Int() % 2 == 0,
	}

	hJson, err := j.Marshal(h)

	if err != nil {
		panic(err)
	}

	_, _, _ = util.DoRequest("PUT",
		fmt.Sprint(url, "/adm/hospital/add"),
		string(hJson), "", "",
	)

	// Cria os medicos
	for k := 1; k <= (size + 1)/10; k++ {
		eid := r.Intn(size + 1)

		if eid == 0 {
			eid++
		}

		m := objetos.Medico{
			MNOME:    uuid.New().String(),
			EID:      uint(eid),
			HID:      uint(i),
			MATIVADO: true,
		}

		mJson, err := j.Marshal(m)

		if err != nil {
			panic(err)
		}

		_, _, _ = util.DoRequest("PUT",
			fmt.Sprint(url, "/adm/hospital/medico/add"),
			string(mJson), "", "",
		)

	}
	wg.Done()
}

func setup(url string, size, usuariosSimuntaneos int) {
	setupPoolChan := make(chan setupJob)

	var wg sync.WaitGroup
	wg.Add(size)

	for i := 0; i < usuariosSimuntaneos; i++ {
		go func(jobs <- chan setupJob) {
			for j := range jobs {
				j.job(j.url, j.index, j.size, j.wg)
			}
		}(setupPoolChan)
	}

	for i := 1; i <= size; i++ {
		setupPoolChan <- setupJob{
			url:   url,
			index: i,
			size:  size,
			wg:    &wg,
			job:   cadastrarEspecialidade,
		}
	}

	wg.Wait()
	wg.Add(size)

	for i := 1; i <= size; i++ {
		setupPoolChan <- setupJob{
			url:   url,
			index: i,
			size:  size,
			wg:    &wg,
			job:   cadastrarHospital,
		}
	}
	wg.Wait()
}

func usuario(url string, index string,wg *sync.WaitGroup) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// O usuario se cadastrar
	u := objetos.Usuario{
		UNOME:        index,
		UEMAIL:       index + "@email.com",
		UPASSWORD:    uuid.New().String() + "@ABC_D",
		UCPF:         util.GerarCpf(),
		UUF:          "SP",
		UCIDADE:      "Uma cidade ai",
		UCEP:         "04014902",
		UENDERECO:    fmt.Sprint("Um endereço ai,", r.Int()),
		UCOMPLEMENTO: fmt.Sprint(r.Int()),
	}

	uJson, err := j.Marshal(u)

	if err != nil {
		panic(err)
	}

	_, _, _ = util.DoRequest("POST",
		fmt.Sprint(url, "/cadastrar/usuario"),
		string(uJson), "", "",
	)

	// o usuario realiza o Login
	loginMap := map[string]string {
		"email": u.UEMAIL,
		"password": u.UPASSWORD,
	}

	loginMapJson, err := j.Marshal(loginMap)

	if err != nil {
		panic(err)
	}

	_, _, tJson := util.DoRequest("POST",
		fmt.Sprint(url, "/login"),
		string(loginMapJson), "", "",
	)

	m := map[string]string{}

	if err = j.Unmarshal(tJson, &m); err != nil {
		panic(err)
	}

	u.UTOKEN = m["token"]

	//O usuario logado lista as especialidades
	_, _, eJson := util.DoRequest("GET",
		fmt.Sprint(url, "/usuario/especialidades/1"),
		"", u.UEMAIL, u.UTOKEN,
	)

	var e []objetos.Especialidade

	if err = j.Unmarshal(eJson, &e); err != nil {
		panic(err)
	}

	// O usuario logado lista os hospitais
	_, _, _ = util.DoRequest("GET",
		fmt.Sprint(url, "/usuario/hospitais/1"),
		"", u.UEMAIL, u.UTOKEN,
	)

	// O usuario logado lista os medicos de uma especialidade aleatoria
	i := r.Intn(len(e))
	_, _, mJson := util.DoRequest("GET",
		fmt.Sprint(url, "/usuario/hospital/medicos/", i),
		"", u.UEMAIL, u.UTOKEN,
	)

	var medicos []objetos.Medico
	if err = j.Unmarshal(mJson, &medicos); err != nil {
		panic(err)
	}

	//O usuario logado marca uma consulta com um medico aleatorio
	i = r.Intn(len(medicos))

	agendamentoMap := map[string]int64{
		"data": time.Now().Unix() + 10000,
	}

	agendamentoJson, err := j.Marshal(agendamentoMap)

	if err != nil {
		panic(err)
	}

	_, _, _ = util.DoRequest("PUT",
		fmt.Sprint(url, "/usuario/hospital/agenda/", i, "/add"),
		string(agendamentoJson), u.UEMAIL, u.UTOKEN,
	)

	logOffMap := map[string]string{
		"email": u.UEMAIL,
		"password": u.UPASSWORD,
	}

	logOffJson, err := j.Marshal(logOffMap)

	if err != nil {
		panic(err)
	}

	// O usuario realiza o logoff
	_, _, _ = util.DoRequest("POST",
		fmt.Sprint(url, "/api/v1/logoff"),
		string(logOffJson), u.UEMAIL, u.UTOKEN,
	)

	wg.Done()
}

