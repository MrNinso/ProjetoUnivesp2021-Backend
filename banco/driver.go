package banco

import (
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"
	"time"
)

type DriverBancoDados interface {
	NewConn(host, port, username, password, database string) (DriverBancoDados, error)

	CadastarUsuario(usuario objetos.Usuario) uint8

	Login(uemail, upassword string) string

	IsValidToken(uemail, utoken string) (bool, string)

	Logoff(uemail, token string) uint8

	//Atividades do Usuario
	ListarHospitais() []objetos.Hospital

	ListarHospitaisPorPlanoConvenio(cpid uint64) []uint //HID

	ListarEspecialidadesHospital(hid uint) []uint //EID

	ListarEspecialidades() []objetos.Especialidade

	ListarConvenios() []objetos.Convenio

	ListarPlanosConvenio(cid uint) []objetos.Planos

	ListarMedicoPorEspecialiade(eid uint) []objetos.Medico

	ListarAgendamentos(utoken string) []objetos.Agendamento

	ListarAgendamentosDoMedico(mid uint64) []objetos.Agendamento

	MarcarConsulta(utoken string, did uint64, mid uint64, data time.Time) uint8

	FavoritarHospital(utoken string, hid uint) uint8

	ListarHospitaisFavoritos(utoken string) []objetos.Hospital

	AdicionarDependete(utoken string, dependete objetos.Dependente) uint8

	ListarDependentes(utoken string) []objetos.Dependente

	RemoverDependente(utoken string, did uint64) uint8

	//Atividades Administrativas
	AdicionarHospital(hospital objetos.Hospital) uint8

	AdicionarConvenioHospital(cpid uint64, hid uint) uint8

	RemoverConvenioHospital(cpid uint64, hid uint) uint8

	AdicionarMedico(medico objetos.Medico) uint8

	AdicionarEspecialidade(especialidade string) uint8

	AdicionarConvenio(nome string) uint8

	AdicionarPlanoConvenio(cid uint64, nome string) uint8
}
