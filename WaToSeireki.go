package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {
	flag.Parse()

	log.Print("Start\r\n")

	// ファイルを読み込んで二次元配列に入れる
	records := readFile(flag.Arg(0))

	// 和暦を西暦に変換
	convRecords := setSeireki(records)

	// ファイルへ書き出す
	saveFile(flag.Arg(0), convRecords)

	log.Print("Finesh \r\n")

}

func readFile(fileName string) [][]string {
	// 入力ファイル準備
	infile, err := os.Open(fileName)
	failOnError(err)
	defer infile.Close()

	reader := csv.NewReader(transform.NewReader(infile, japanese.ShiftJIS.NewDecoder()))
	reader.Comma = '\t'

	//CSVファイルを２次元配列に展開
	readrecords := make([][]string, 0)
	for {
		record, err := reader.Read() // 1行読み出す
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}

		readrecords = append(readrecords, record)
	}

	return readrecords
}

func saveFile(filename string, saverecords [][]string) {
	//出力ファイル準備
	outDir, outfileName := filepath.Split(filename)
	pos := strings.LastIndex(outfileName, ".")
	outfile, err := os.Create(outDir + outfileName[:pos] + ".txt")
	failOnError(err)
	defer outfile.Close()

	writer := csv.NewWriter(transform.NewWriter(outfile, japanese.ShiftJIS.NewEncoder()))
	writer.Comma = '\t'
	writer.UseCRLF = true

	for _, out_record := range saverecords {
		writer.Write(out_record)
	}

	writer.Flush()
}

func setSeireki(r [][]string) [][]string {
	// 二次元配列から
	// 和暦Sxx.xx.xx形式を見つけて、西暦に変換

	for i, rec := range r {
		for j, cel := range rec {
			//文字列の長さは８桁か
			if len(cel) == 9 {
				//年号チェック
				if cel[0:1] == "M" || cel[0:1] == "T" || cel[0:1] == "S" || cel[0:1] == "H" {
					//月日チェック
					if cel[3:4] == "." && cel[6:7] == "." {
						Wa := cel[0:1]
						Y := cel[1:3]
						iY, _ := strconv.Atoi(Y)
						M := cel[4:6]
						D := cel[7:9]

						switch Wa {
						case "M":
							iY = 1900 + iY - 33
						case "T":
							iY = 1900 + iY + 11
						case "S":
							iY = 1900 + iY + 25
						case "H":
							iY = 1900 + iY + 88
						}

						r[i][j] = strconv.Itoa(iY) + "/" + M + "/" + D
					}
				}
			}
		}
	}

	return r

}
