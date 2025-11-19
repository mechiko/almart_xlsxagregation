package main

import (
	"agregat/domain"
	"agregat/process"
	"agregat/reductor"
	"agregat/repo/configdb"
	"agregat/zaplog"
	"flag"
	"fmt"
	"os"

	"github.com/mechiko/dbscan"
	"github.com/mechiko/utility"
	"go.uber.org/zap"
)

var file = flag.String("file", "", "file to parse xlsx")

var Mode = "development"

func errMessageExit(title string, errDescription string) {
	utility.MessageBox(title, errDescription)
	os.Exit(-1)
}

func init() {
	flag.Parse()
}

func main() {
	var logsOutConfig = map[string][]string{
		"logger": {"stdout"},
	}
	zl, err := zaplog.New(logsOutConfig, true)
	if err != nil {
		errMessageExit("ошибка", err.Error())
	}

	lg, err := zl.GetLogger("logger")
	if err != nil {
		errMessageExit("ошибка", err.Error())
	}
	loger := lg.Sugar()
	loger.Info("agregat started")

	fileXLSX := *file
	if fileXLSX == "" {
		fileXLSX, err = utility.DialogOpenFile([]utility.FileType{utility.Excel, utility.All}, "", "in")
		if err != nil {
			errMessageExit("ошибка", err.Error())
		}
	}

	proc, err := process.New(fileXLSX)
	if err != nil {
		errMessageExit("запуск обработки с ошибкой ", err.Error())
	}

	err = proc.ReadXlsx(fileXLSX)
	if err != nil {
		errMessageExit("ошибка чтения файла", err.Error())
	}

	outDir := "."
	outDir, err = utility.DialogSelectDir(outDir)
	if err != nil {
		errMessageExit("ошибка выбора пути назначения", err.Error())
	}

	err = proc.ScanPalet()
	if err != nil {
		// errMessageExit("ошибка ScanRecords", err.Error())
		loger.Errorf("ошибка ScanPalet %v", err)
	}

	// записываем короба и палеты раньше пусть будут
	err = proc.Save(outDir)
	if err != nil {
		// errMessageExit("ошибка ScanRecords", err.Error())
		loger.Errorf("ошибка Save %v", err)
	}

	// запись отчетов нанесения в БД
	// err = proc.WriteUtilisation()
	// if err != nil {
	// 	errMessageExit("ошибка ScanRecords", err.Error())
	// }
	fileSuccessReport, err := proc.SuccessHtml(outDir)
	if err != nil {
		errMessageExit("ошибка формирования отчета", err.Error())
	}
	utility.OpenFileInShell(fileSuccessReport)

}

func readModel(loger *zap.SugaredLogger, repo domain.Repo, model *reductor.Model) (err error) {
	info := repo.Info(dbscan.Config)
	if info == nil {
		return fmt.Errorf("базы конфиг не найдено")
	}
	db, err := configdb.New(loger, info, dbscan.Config)
	if err != nil {
		return fmt.Errorf("error open config.db")
	}
	defer func() {
		if errclose := db.Close(); errclose != nil {
			if err != nil {
				err = fmt.Errorf("%v %w", errclose, err)
			} else {
				err = fmt.Errorf("%v", errclose)
			}
		}
	}()
	model.Inn, err = db.Key("inn")
	if err != nil {
		return fmt.Errorf("error open config.db")
	}
	model.Kpp, err = db.Key("kpp")
	if err != nil {
		return fmt.Errorf("error open config.db")
	}
	return nil
}
