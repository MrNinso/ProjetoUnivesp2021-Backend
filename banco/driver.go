package banco

import (
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"
	"time"
)

type DriverBancoDados interface {
	CadastarUsuario(usuario objetos.Usuario) uint8

	Login(uemail, upassword string) string

	IsValidToken(uemail, utoken string) (bool, string)

	Logoff(uemail, token string) uint8

	//Atividades do Usuario
	ListarEspecialidades(page uint8) []objetos.Especialidade

	ListarMedicoPorEspecialiade(eid uint) []objetos.Medico

	ListarAgendamentosDoMedico(mid uint64, page uint8) []objetos.Agendamento

	ListarHospitais(page uint8) []objetos.Hospital

	MarcarConsulta(utoken string, mid uint64, data time.Time) uint8

	//Atividades Administrativas
	AdicionarHospital(hospital objetos.Hospital) uint8

	AdicionarMedico(medico objetos.Medico) uint8

	AdicionarEspecialidade(especialidade string) uint8
}
