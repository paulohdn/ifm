package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"ifm/apoio"
	"ifm/pkg"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var Tli_id, Leg_seq, Leg_ordem, Lin_seq int
var Fra_id, Leg_id, Pal_id, Lin_id int64
var QtdeControlaLegendas, QtdeControlaFrases, QtdeControlaLinguagensPV, QtdeControlaLinguagensSl, QtdeControlaLinguagensHE, QtdeControlaPalavras int
var IdR, QtR64 int64
var QtR int
var exc = make(map[string]int)
var contaPalavras = make(map[string]int)

// var controlaFrases = make(map[string]int64)
// var controlaLinguagensPV = make(map[string]int64)
// var controlaLinguagensSl = make(map[string]int64)
// var controlaLinguagensHE = make(map[string]int64)
var controlaPalavras = make(map[string]int)

func limpaString(s string) string {
	sr := ""
	for _, v := range s {
		_, encontrou := exc[string(v)]
		if !encontrou {
			sr = sr + string(v)
		}
	}
	return sr
}

func setaParamentrosLimpaString(fase string) {
	switch fase {
	case "arquivo":
		exc = make(map[string]int)
		exc["["] = 91
		exc["]"] = 93
		exc[";"] = 0
	case "palavras":
		exc = make(map[string]int)
		exc["."] = 0
		exc[","] = 0
		exc[":"] = 0
		exc["!"] = 0
		exc["?"] = 0
		exc["&"] = 0
		exc["\""] = 0
	case "chatGPT":
		exc["\""] = 0
		exc["\\"] = 0
		exc[":"] = 0
	case "pattern":
		exc["\\"] = 0
	}

}
func montaChave(valor int) string {
	s := strconv.Itoa(valor)
	sr := ""
	for i := 0; i < 4-len(s); i++ {
		sr = sr + "0"
	}
	return sr + s
}

const textThatHappenWhenIsNotPhrasalVerb string = "there is no phrasal verb"

func thereArePhrasalVerb(chatGPTResponse string) string {
	sr := ""
	found := strings.Contains(strings.ToLower(chatGPTResponse), textThatHappenWhenIsNotPhrasalVerb)
	if !found {
		setaParamentrosLimpaString("chatGPT")
		sr = limpaString(chatGPTResponse)
		lowercaseSr := strings.ToLower(sr)
		lowercaseSearch := "phrasal verb"
		sr = strings.Replace(lowercaseSr, lowercaseSearch, "", -1)
		sr = strings.TrimSpace(sr)
		// arrayWords := strings.Split(sr, " ")
		// if len(arrayWords) == 1 {
		// 	sr = ""
		// }
	}
	return sr
}

const textThatHappenWhenIsNotEnglishExpression string = "there is no english expression"

func thereAreEnglishExpression(chatGPTResponse string) string {
	sr := ""
	found := strings.Contains(strings.ToLower(chatGPTResponse), textThatHappenWhenIsNotEnglishExpression)
	if !found {
		setaParamentrosLimpaString("chatGPT")
		sr = limpaString(chatGPTResponse)
		lowercaseSr := strings.ToLower(sr)
		lowercaseSearch := "english expression"
		sr = strings.Replace(lowercaseSr, lowercaseSearch, "", -1)
		lowercaseSr = sr
		lowercaseSearch = "\\n message n/a"
		sr = strings.Replace(lowercaseSr, lowercaseSearch, "", -1)
		sr = strings.TrimSpace(sr)
		// arrayWords := strings.Split(sr, " ")
		// if len(arrayWords) == 1 {
		// 	sr = ""
		// }
	}
	return sr
}

const textThatHappenWhenIsNotEnglishSlang string = "there is no slang"

func thereAreEnglishSlang(chatGPTResponse string) string {
	sr := ""
	if len(chatGPTResponse) > 0 {
		found := strings.Contains(strings.ToLower(chatGPTResponse), textThatHappenWhenIsNotEnglishSlang)
		if !found {
			setaParamentrosLimpaString("chatGPT")
			sr = limpaString(chatGPTResponse)
			lowercaseSr := strings.ToLower(sr)
			lowercaseSearch := "english slang"
			sr = strings.Replace(lowercaseSr, lowercaseSearch, "", -1)
			lowercaseSr = sr
			lowercaseSearch = "\\n message n/a"
			sr = strings.Replace(lowercaseSr, lowercaseSearch, "", -1)
			sr = strings.TrimSpace(sr)
			// arrayWords := strings.Split(sr, " ")
			// if len(arrayWords) == 1 {
			// 	sr = ""
			// }
		}
	}
	return sr
}

// funções e tipos para ordenação por valor
type Pair struct {
	Key   string
	Value int
}
type PairList []Pair

func (p PairList) Len() int      { return len(p) }
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }
func (p PairList) Less(i, j int) bool {
	return montaChave(p[i].Value)+p[i].Key > montaChave(p[j].Value)+p[j].Key
}

var phraseEnglish, phrasePortuguese string
var insertStmtLegenda, insertStmtPalavra, updateStmtPalavra, insertStmtDob_Pal, insertStmtFrase, updateStmtFrase *sql.Stmt
var insertStmtDob_Fra, insertStmtLinguagem, insertStmtLin_Fra, updateStmtLinguagem *sql.Stmt
var selectStmtLinguagem, selectStmtFrase, selectStmtPalavra, selectStmtDetObraAudivisual *sql.Stmt
var db *sql.DB
var err error

func openDB() error {
	dbConnectionString := ""
	dbConnectionString = "usr_ifm:usr_ifm_Senha@tcp(mysql247.umbler.com:41890)/bd_ifm"
	dbConnectionString = "root:@tcp(localhost:3306)/db_ifm"
	db, err = sql.Open("mysql", dbConnectionString)
	if err != nil {
		fmt.Println("Falha ao conectar ao banco de dados:", err)
		return err
	}
	return nil
}
func preparaInsertUpdateDeleteConsulta() {
	insertStmtLegenda, err = db.Prepare("INSERT INTO legenda (leg_seq, dob_id, leg_seq_tempo) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do insert legenda:", err)
		return
	}
	insertStmtPalavra, err = db.Prepare("INSERT INTO palavra (pal_nm, pal_ocor, pal_sig_pt) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do insert da palavra:", err)
		return
	}
	updateStmtPalavra, err = db.Prepare("UPDATE palavra SET pal_ocor = pal_ocor + 1 WHERE pal_id = ?")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do update da palavra:", err)
		return
	}
	insertStmtDob_Pal, err = db.Prepare("INSERT INTO dob_pal (dob_id, pal_id, leg_id, leg_seq, leg_ordem, pal_ordem) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do insert da dob_pal:", err)
		return
	}
	insertStmtFrase, err = db.Prepare("INSERT INTO frase (fra_nm, fra_ocor, fra_sig_pt) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do insert da frase:", err)
		return
	}
	updateStmtFrase, err = db.Prepare("UPDATE frase SET fra_ocor = fra_ocor + 1 WHERE fra_id = ?")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do update da frase:", err)
		return
	}
	insertStmtDob_Fra, err = db.Prepare("INSERT INTO dob_fra (dob_id, fra_id, leg_id, leg_seq, leg_ordem) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do insert da dob_fra:", err)
		return
	}
	insertStmtLinguagem, err = db.Prepare("INSERT INTO linguagem (tli_id, lin_ocor, lin_texto) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do insert da linguagem:", err)
		return
	}
	updateStmtLinguagem, err = db.Prepare("UPDATE linguagem SET lin_ocor = lin_ocor + 1 WHERE lin_id = ?")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do update da linguagem:", err)
		return
	}
	insertStmtLin_Fra, err = db.Prepare("INSERT INTO lin_fra (lin_id, dob_id, fra_id, leg_id, leg_seq, leg_ordem, lin_seq) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do insert da lin_fra:", err)
		return
	}
	selectStmtFrase, err = db.Prepare("SELECT fra_id, fra_ocor FROM frase WHERE fra_nm = ?")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do select da frase:", err)
		return
	}
	selectStmtDetObraAudivisual, err = db.Prepare("SELECT dob_id FROM det_obra_audivisual WHERE obr_id = ? and dob_temp = ? and dob_seq_temp = ?")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do select da det_obra_audivisual:", err)
		return
	}
	selectStmtLinguagem, err = db.Prepare("SELECT lin_id, lin_ocor FROM linguagem WHERE tli_id = ? and lin_texto = ?")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do select da linguagem:", err)
		return
	}
	selectStmtPalavra, err = db.Prepare("SELECT pal_id, pal_ocor FROM palavra WHERE pal_nm = ?")
	if err != nil {
		fmt.Println("Falha na declaração da preparação do select da palavra:", err)
		return
	}
}
func fechaInsertUpdateDeleteConsulta() {
	insertStmtLegenda.Close()
	insertStmtPalavra.Close()
	updateStmtPalavra.Close()
	insertStmtDob_Pal.Close()
	insertStmtFrase.Close()
	updateStmtFrase.Close()
	insertStmtDob_Fra.Close()
	insertStmtLinguagem.Close()
	updateStmtLinguagem.Close()
	insertStmtLin_Fra.Close()
	selectStmtFrase.Close()
	selectStmtLinguagem.Close()
	selectStmtPalavra.Close()
	selectStmtDetObraAudivisual.Close()
	db.Close()
}
func insertLegenda(leg_seq int, dob_id int64, leg_seq_tempo string) int64 {
	result, err := insertStmtLegenda.Exec(leg_seq, dob_id, leg_seq_tempo)
	if err != nil {
		fmt.Println("Falha ao inserir a linha na legenda:", err)
		return 0
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Falha no retorno do último ID inserido:", err)
		return 0
	}
	return lastInsertID
}
func insertPalavra(pal_nm string, pal_ocor int64, pal_sig_pt string) int64 {
	result, err := insertStmtPalavra.Exec(pal_nm, pal_ocor, pal_sig_pt)
	if err != nil {
		fmt.Println("Falha ao inserir a linha na palavra:", err)
		return 0
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Falha no retorno do último ID inserido:", err)
		return 0
	}
	return lastInsertID
}
func updatePalavra(pal_id int64) {
	_, err := updateStmtPalavra.Exec(pal_id)
	if err != nil {
		fmt.Println("Falha ao atualiza a linha na palavra:", err)
		return
	}
}
func insertDob_Pal(dob_id, pal_id, leg_id int64, leg_seq, leg_ordem int, pal_ordem int) {
	_, err := insertStmtDob_Pal.Exec(dob_id, pal_id, leg_id, leg_seq, leg_ordem, pal_ordem)
	if err != nil {
		fmt.Println("Falha ao inserir a linha na dob_pal:", err)
		return
	}
}
func insertFrase(fra_nm string, fra_ocor int64, fra_sig_pt string) int64 {
	result, err := insertStmtFrase.Exec(fra_nm, fra_ocor, fra_sig_pt)
	if err != nil {
		fmt.Println("Falha ao inserir a linha na frase:", err)
		return 0
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Falha no retorno do último ID inserido na frase:", err)
		return 0
	}
	return lastInsertID
}
func updateFrase(fra_id int64) {
	_, err := updateStmtFrase.Exec(fra_id)
	if err != nil {
		fmt.Println("Falha ao atualiza a linha na frase:", err)
		return
	}
}
func insertDob_Fra(dob_id, fra_id, leg_id int64, leg_seq, leg_ordem int) {
	_, err := insertStmtDob_Fra.Exec(dob_id, fra_id, leg_id, leg_seq, leg_ordem)
	if err != nil {
		fmt.Println("Falha ao inserir a linha na dob_fra:", err)
		return
	}
}
func insertLinguagem(tli_id, lin_ocor int, lin_texto string) int64 {
	result, err := insertStmtLinguagem.Exec(tli_id, lin_ocor, lin_texto)
	if err != nil {
		fmt.Println("Falha ao inserir a linha na linguagem:", err)
		return 0
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Falha no retorno do último ID inserido na linguagem:", err)
		return 0
	}
	return lastInsertID
}
func updateLinguagem(lin_id int64) {
	_, err := updateStmtLinguagem.Exec(lin_id)
	if err != nil {
		fmt.Println("Falha ao atualiza a linha na linguagem:", err)
		return
	}
}
func insertLin_Fra(lin_id, dob_id, fra_id, leg_id int64, leg_seq, leg_ordem, lin_seq int) {
	_, err := insertStmtLin_Fra.Exec(lin_id, dob_id, fra_id, leg_id, leg_seq, leg_ordem, lin_seq)
	if err != nil {
		fmt.Println("Falha ao inserir a linha na lin_fra:", err)
		return
	}
}
func consultaOcorFrase(fra_nm string) (fra_idR int64) {
	var fra_id int64
	var fra_ocor int
	rows, err := selectStmtFrase.Query(fra_nm)
	if err != nil {
		fmt.Println("Falha ao executar a instrução SELECT para frase:", err)
		return 0
	}
	for rows.Next() {
		err := rows.Scan(&fra_id, &fra_ocor)
		if err != nil {
			fmt.Println("Falha ao recuperar os dados do registro da frase:", err)
			return 0
		}
	}
	return fra_id
}
func consultaOcorLinguagem(tli_id int, lin_texto string) (lin_idR int64) {
	var lin_id int64
	var lin_ocor int
	rows, err := selectStmtLinguagem.Query(tli_id, lin_texto)
	if err != nil {
		fmt.Println("Falha ao executar a instrução SELECT para linguagem:", err)
		return 0
	}
	for rows.Next() {
		err := rows.Scan(&lin_id, &lin_ocor)
		if err != nil {
			fmt.Println("Falha ao recuperar os dados do registro da linguagem:", err)
			return 0
		}
	}
	return lin_id
}

func consultaOcorPalavra(pal_nm string) (pal_idR int64) {
	var pal_id int64
	var pal_ocor int64
	rows, err := selectStmtPalavra.Query(pal_nm)
	if err != nil {
		fmt.Println("Falha ao executar a instrução SELECT para palavra:", err)
		return 0
	}
	for rows.Next() {
		err := rows.Scan(&pal_id, &pal_ocor)
		if err != nil {
			fmt.Println("Falha ao recuperar os dados do registro da palavra:", err)
			return 0
		}
	}
	return pal_id
}
func consultaDetObraAudivisual(obr_id int, dob_temp int, dob_seq_temp int) (dob_indR int64) {
	var dob_id int64
	rows, err := selectStmtDetObraAudivisual.Query(obr_id, dob_temp, dob_seq_temp)
	if err != nil {
		fmt.Println("Falha ao executar a instrução SELECT para det_obra_audivisual:", err)
		return 0
	}
	for rows.Next() {
		err := rows.Scan(&dob_id)
		if err != nil {
			fmt.Println("Falha ao recuperar os dados do registro da det_obra_audivisual:", err)
			return 0
		}
	}
	return dob_id
}
func checkAndReplace(inputLine string) string {
	if strings.Contains(inputLine, "www.") || strings.Contains(inputLine, "http://") || strings.Contains(inputLine, "https://") {
		return "."
	}
	inputLine = strings.Replace(inputLine, "\ufeff", "", -1)
	return inputLine
}
func retiraCaracteresFrase(inputLine string) string {
	runeLeading := []rune{' ', '\'', '_', '-', '(', '.', ' '}
	for _, r := range runeLeading {
		inputLine = apoio.RemoveLeadingChar(inputLine, r)
	}
	runeTrailing := []rune{' ', '\'', '_', '-', ')', '.'}
	for _, r := range runeTrailing {
		inputLine = apoio.RemoveTrailingChar(inputLine, r)
	}
	inputLine = apoio.RemoveFirstOrLastDoubleQuote(inputLine)
	inputLine = apoio.ReplaceSubstring(inputLine, "H6Y", "HEY")
	inputLine = apoio.ReplacePontos(inputLine)
	return inputLine
}
func retiraCaracteresPalavra(inputWord string) string {
	runeLeading := []rune{' ', '\'', '_', '-', '(', '.', ' '}
	for _, r := range runeLeading {
		inputWord = apoio.RemoveLeadingChar(inputWord, r)
	}
	runeTrailing := []rune{' ', '\'', '_', '-', ')', '.'}
	for _, r := range runeTrailing {
		inputWord = apoio.RemoveTrailingChar(inputWord, r)
	}
	inputWord = apoio.RemoveFirstOrLastDoubleQuote(inputWord)
	return inputWord
}
func main() {
	openDB()
	if db == nil {
		return
	}
	preparaInsertUpdateDeleteConsulta()
	defer fechaInsertUpdateDeleteConsulta()

	var DOB_ID int64 = 1 // detalhe da obra audivisual (det_obra_audivisual) // 1 = obra = 1, Friends, Temp 1
	for k := 6; k < 11; k++ {

		fn := ""
		rootDir := fmt.Sprintf("./subtitles%d/", k)

		var srtFiles []string

		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Check if the file has .srt extension
			if !info.IsDir() && filepath.Ext(info.Name()) == ".srt" {
				srtFiles = append(srtFiles, info.Name())
			}

			return nil
		})

		if err != nil {
			fmt.Println("Error:", err)
		}

		for _, fileName := range srtFiles {
			fn = rootDir + fileName
			season, episode := apoio.ReturnSeasonEpisode(fileName)
			obr_id := 1
			DOB_ID = consultaDetObraAudivisual(obr_id, season, episode)
			fmt.Printf("%s, dob_id: %d, season: %d, episode: %d\n", fn, DOB_ID, season, episode)
			InsereLegenda(fn, DOB_ID)
		}

		// fn := "./subtitles/friends.s01e01.720p.bluray.x264-psychd-resumo.teste-srt"
		// InsereLegenda(fn, DOB_ID)
	}
}

func InsereLegenda(fname string, DOB_ID int64) {
	readFile, err := os.Open(fname)
	apoio.Check(err)
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)
	linha := ""
	var contaLinha int = 1
	var qtdeContaPalavras int

	//var encontrado bool
	saltos := 2
	achouLinha := false
	var palavras []string
	var tempoLegenda string = ""
	var linhaAux string
	// var  chatGPTResponse, phrasalVerb, EnglishExpression, EnglishSlang string
	Leg_id = 0
	Leg_seq = 0
	Leg_ordem = 0
	for fileScanner.Scan() {
		linha = fileScanner.Text()
		linha = checkAndReplace(linha)
		valorLinha, err := strconv.Atoi(linha)
		if err == nil && contaLinha == valorLinha {
			Leg_seq = contaLinha
			contaLinha++
			achouLinha = true
			saltos = 2
		}
		if saltos == 1 {
			tempoLegenda = strings.Trim(linha, " ")
			if Leg_seq == 0 {
				fmt.Println(Leg_seq)
			}
			Leg_id = insertLegenda(Leg_seq, DOB_ID, tempoLegenda)
			Leg_ordem = 0
		}
		if (achouLinha || err != nil) && saltos > 0 || linha == "" {
			achouLinha = false
			saltos--
			continue
		}
		Leg_ordem++
		setaParamentrosLimpaString("arquivo")
		linha = limpaString(linha)
		linha = retiraCaracteresFrase(linha)
		linhaAux = strings.ToUpper(linha)
		palavras = strings.Fields(linhaAux)

		phraseEnglish = linha
		phrasePortuguese = linha
		phrasePortuguese = pkg.TranslateTextGoogle(linha)
		if apoio.LeftString(linha, 1) == "-" {
			phraseEnglish = " " + phraseEnglish
			phrasePortuguese = " " + phrasePortuguese
		}
		Fra_id = consultaOcorFrase(phraseEnglish)
		if Fra_id == 0 {
			Fra_id = insertFrase(phraseEnglish, 1, phrasePortuguese)
		} else {
			updateFrase(Fra_id)
		}
		insertDob_Fra(DOB_ID, Fra_id, Leg_id, Leg_seq, Leg_ordem)
		controlaPalavras = make(map[string]int)
		setaParamentrosLimpaString("palavras")
		for par_ordem, v := range palavras {
			v = limpaString(v)
			v = retiraCaracteresPalavra(v)
			qtdeContaPalavras = contaPalavras[v]
			QtdeControlaPalavras = controlaPalavras[v]
			if v == "-" || v == "" {
			} else {
				contaPalavras[v] = qtdeContaPalavras + 1
				controlaPalavras[v] = QtdeControlaPalavras + 1
				Pal_id = consultaOcorPalavra(v)
				if Pal_id == 0 {
					wordPortuguese := v
					wordPortuguese = pkg.TranslateTextGoogle(v)
					Pal_id = insertPalavra(v, 1, wordPortuguese)
				} else {
					updatePalavra(Pal_id)
				}
				insertDob_Pal(DOB_ID, Pal_id, Leg_id, Leg_seq, Leg_ordem, par_ordem+1)
			}
		}

		// chatGPTResponse = pkg.ChapGPTExample2("There are phrasal verb in sentense below:\n\"" + linha + "\"? Return just the phrasal verb or the message: \"There is no phrasal verb\".")
		// phrasalVerb = thereArePhrasalVerb(chatGPTResponse)
		// if len(phrasalVerb) > 0 {
		// 	phraseEnglish = phrasalVerb
		// 	phrasePortuguese = "" //pkg.TranslateTextGoogle(phrasalVerb)
		// 	wDestPhrasalVerbsBuffer.WriteString(phraseEnglish + ";" + phrasePortuguese + "\r\n")
		// 	Tli_id = 1
		// 	Lin_seq = 1
		// 	Lin_id = consultaOcorLinguagem(Tli_id, phraseEnglish)
		// 	if Lin_id == 0 {
		// 		Lin_id = insertLinguagem(Tli_id, 1, phraseEnglish)
		// 	} else {
		// 		updateLinguagem(Lin_id)
		// 	}
		// 	insertLin_Fra(Lin_id, DOB_ID, Fra_id, Leg_id, Leg_seq, Leg_ordem, Lin_seq)
		// }
		// chatGPTResponse = pkg.ChapGPTExample2("There are English expression in sentense below:\n\"" + linha + "\"? Return just the English expression or the message: \"There is no English expression\".")
		// EnglishExpression = thereAreEnglishExpression(chatGPTResponse)
		// if len(EnglishExpression) > 0 {
		// 	phraseEnglish = EnglishExpression
		// 	phrasePortuguese = "" //pkg.TranslateTextGoogle(EnglishExpression)
		// 	wDestEnglishExpressionsBuffer.WriteString(phraseEnglish + ";" + phrasePortuguese + "\r\n")
		// 	Tli_id = 2
		// 	Lin_seq = 1
		// 	Lin_id = consultaOcorLinguagem(Tli_id, phraseEnglish)
		// 	if Lin_id == 0 {
		// 		Lin_id = insertLinguagem(Tli_id, 1, phraseEnglish)
		// 	} else {
		// 		updateLinguagem(Lin_id)
		// 	}
		// 	insertLin_Fra(Lin_id, DOB_ID, Fra_id, Leg_id, Leg_seq, Leg_ordem, Lin_seq)
		// }
		// chatGPTResponse = pkg.ChapGPTExample2("There are slang in sentense below:\n\"" + linha + "\"?") //  Return just the slang or the message: \"There is no slang\".
		// EnglishSlang = thereAreEnglishSlang(chatGPTResponse)
		// if len(EnglishSlang) > 0 {
		// 	phraseEnglish = EnglishExpression
		// 	phrasePortuguese = "" //pkg.TranslateTextGoogle(EnglishExpression)
		// 	wDestEnglishSlangsBuffer.WriteString(phraseEnglish + ";" + phrasePortuguese + "\r\n")
		// 	Tli_id = 3
		// 	Lin_seq = 1
		// 	Lin_id = consultaOcorLinguagem(Tli_id, phraseEnglish)
		// 	if Lin_id == 0 {
		// 		Lin_id = insertLinguagem(Tli_id, 1, phraseEnglish)
		// 	} else {
		// 		updateLinguagem(Lin_id)
		// 	}
		// 	insertLin_Fra(Lin_id, DOB_ID, Fra_id, Leg_id, Leg_seq, Leg_ordem, Lin_seq)
		// }

	}
	qtdePalavras := 0

	// ordenando pela chave
	keys := make([]string, 0, len(contaPalavras))
	for k := range contaPalavras {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// ordenando pela qtde
	p := make(PairList, len(contaPalavras))
	i := 0
	for k, v := range contaPalavras {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	for _, k := range p {
		qtdePalavras += k.Value
	}

	fmt.Println("Total de falas:", contaLinha-1)
	fmt.Println("Total de palavras:", qtdePalavras)
	fmt.Println("Total de palavras únicas:", len(contaPalavras))

	readFile.Close()
}
