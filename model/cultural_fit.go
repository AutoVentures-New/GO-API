package model

type CulturalFitQuestion struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	AnswerA string `json:"answer_a"`
	AnswerB string `json:"answer_b"`
}

var CulturalFit = []CulturalFitQuestion{
	{Id: 1, Title: "Você geralmente", AnswerA: "se mistura bem, ou", AnswerB: "prefere ser quieto e reservado?"},
	{Id: 2, Title: "Quando você tem trabalhos e procedimentos complexos na sua rotina, você prefere", AnswerA: "Se certificar que não vai errar, ou", AnswerB: "identificar oportunidades para simplificar?"},
	{Id: 3, Title: "Quando você nota que está com o prazo apertado, você prefere", AnswerA: "entregar na data, mesmo sabendo que poderia ter feito melhor, ou", AnswerB: "atrasar, mas garantir que ficou melhor? "},
	{Id: 4, Title: "Quando você tem um trabalho especial para fazer, você gosta de", AnswerA: "organizá-lo cuidadosamente antes de começar, ou", AnswerB: "descobrir o que é necessário à medida que avança?"},
	{Id: 5, Title: "Quando a empresa muda a metodologia de trabalho ou uma plataforma que você já tem domínio, você", AnswerA: "se preocupa porque terá que aprender de novo e pode errar, ou", AnswerB: "apóia as mudanças e rapidamente se adapta?"},
	{Id: 6, Title: "É mais confortável para você", AnswerA: "trabalhar em equipe ou sozinho, mas interagindo com outros, ou", AnswerB: "trabalhar sozinho sem interações, para não perder o foco?"},
	{Id: 7, Title: "Ao tomar uma decisão que afeta a empresa, é mais importante para você", AnswerA: "fazer uma boa análise dos fatos para não errar, ou", AnswerB: "considerar os sentimentos e opiniões das pessoas?"},
	{Id: 8, Title: "Se você presenciasse um conflito na empresa, sua atitude mais provável seria", AnswerA: "apoiar quem está certo, ou", AnswerB: "apaziguar os ânimos?"},
	{Id: 9, Title: "Se você visse um colega de trabalho fazendo algo que fere as normas da empresa, você", AnswerA: "manteria distância para não ser envolvido, ou", AnswerB: "buscaria os meios adequados para que sua conduta fosse corrigida?"},
	{Id: 10, Title: "Quando você inicia um grande projeto com vencimento em uma semana, você", AnswerA: "reserva um tempo para listar as coisas separadamente e a ordem de fazê-las, ou", AnswerB: "mergulha nele?"},
	{Id: 11, Title: "Se seu superior solicita uma demanda fora do que está estabelecido em suas atividades, você", AnswerA: "realiza de boa vontade e vê como oportunidade, ou", AnswerB: "acha injusto, pois não estava combinado antes?"},
	{Id: 12, Title: "Você prefere", AnswerA: "Seguir os métodos estabelecidos de fazer o trabalho, ou", AnswerB: "analisar o que pode ser melhorado na forma de fazer o trabalho?"},
	{Id: 13, Title: "Realizar um trabalho com prazos apertados ou sob pressão", AnswerA: "prejudica seu desempenho, pois não foi planejado, ou", AnswerB: "É comum e você se adapta para realizar o combinado?"},
	{Id: 14, Title: "Para você, o que define melhor Responsabilidae Social e Ambiental", AnswerA: "a empresa manter ações sociais corporativas que beneficiam a comunidade, ou", AnswerB: "empresa e colaboradores protegerem o meio ambiente e serem bons integrantes da comunidade?"},
	{Id: 15, Title: "Você prefere", AnswerA: "se adaptar a novas situações e desafios, muitas vezes imprevisiveis, ou", AnswerB: "fazer com ótima qualidade trabalhos que você já domina?"},
}

type CompanyCulturalFit struct {
	CompanyID int64                      `json:"company_id"`
	Answers   []CompanyCulturalFitAnswer `json:"answers"`
}

type CompanyCulturalFitAnswer struct {
	CulturalFitID int64  `json:"cultural_fit_id"`
	Answer        string `json:"answer"`
}
