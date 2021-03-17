package banco

import "github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"

type DriverBancoDados interface {
	Login(lid uint, uemail, upassword string) string

	Logoff(uemail, token string) uint

	CadastarUsuario(lid uint, tokenAdmin string, usuario objetos.Usuario) uint

	AtualizarUsuario(lid uint, tokenAdmin string, usuario objetos.Usuario) uint

	ListarUsuarios(lid uint, tokenAdmin string) []objetos.Usuario

	//TODO ADD CRIAR LOJA

	CriarProduto(lid uint, token string, produto objetos.Produto) uint

	ListarProdutos(lid uint, token string, page uint8) []objetos.Produto

	AtualizarProduto(lid uint, token string, produto objetos.Produto) uint

	RegistrarOrdem(lid uint, token string, ordem objetos.Ordem) uint

	ListarEstoque(lid uint, token string, page uint8) []objetos.ItemEstoque
}
