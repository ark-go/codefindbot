// Загружает конфедициальные данные из файла secret.env (например при разработке)
// файл должен находиться рядом с исполняемым файлом, наличие файла не проверяется и не вызывает ошибку.
// Загружает данные из OS Environment (переменные окружения) перезаписывая загруженные из файла.
// Необходимо подготовить простую структуру с требуемыми полями, поля должны быть type string и экспортируемыми
// функция jt.LoadEnvironment(s *struct) принимает параметр ссылку на структуру
// в которую будут вставленны данные если они будут найдены
package jt

import (
	"bufio"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// var SecretEnv *secretEnv

// func ReadSecret() error {
// 	log.Println(rootDir)
// 	secretfile := filepath.Join(rootDir, "secret.json")
// 	data, err := ioutil.ReadFile(secretfile)
// 	if err != nil {
// 		log.Println(err, secretfile)
// 		return err
// 	}

// 	// json data
// 	SecretEnv = &secretEnv{}

// 	// unmarshall it
// 	err = json.Unmarshal(data, SecretEnv)
// 	if err != nil {
// 		log.Println("error:", err)
// 		return err
// 	}
// 	ReadEnvironment(SecretEnv)
// 	return nil
// }

// на вход подаем s указатель на структуру
//	с экспортированными полями, имена полей должны соответствовать переменным
//	поля могут быть  string,int
//	поля будут заполняться, если будут найдены,
//	сначала из файла secret.env
//	а потом, поверх, из переменных окружения os
//
//	файл secret.env должен располагаться рядом с запускаемым файлом
//	# - является комментарием
//	значенем является все что после первого '=' кавычки не нужны, они будут считаны как кавычки в строке
//	ошибок в случае отстутствия файла не будет
//
//	переменные окружения будут перекрывать значения из файла
//
//	переменные окружения не устанавливаются из файла !!!
func LoadEnvironment(s interface{}) { // ждем указатель на структуру
	mapEnv := readFileSecret()       // сначала прочитаем из файла
	el := reflect.ValueOf(s).Elem()  // возвращает значение которое содержит интерфейс (хочет указатель &)
	if el.Kind() == reflect.Struct { // точно ли это структура?
		for i := 0; i < el.NumField(); i++ { // NumField кол-во полей, пройдемся по ним. если не структура то паника
			name := el.Type().Field(i).Name // имя поля
			//f := el.Field(i)
			f := el.FieldByName(name) // получим поле по имени
			if f.IsValid() {          // если есть поле
				// Значение можно изменить, только если оно
				// адресуемо и не был получено
				// с использованием неэкспортированных полей структуры. т.е. с большой буквы
				if f.CanSet() { // можно ли изменить поле, поле должно быть экспортировано иначе false
					// Изменить значение
					var valueStr string
					if mapEnv != nil { // было ли что прочитано из файла
						if val, ok := mapEnv[name]; ok {
							//f.SetString(val) // заполним поле
							valueStr = val
						}
					}
					val, exists := os.LookupEnv(name) // environment получим переменную, если есть
					if exists {
						//f.SetString(val) // заполним поле
						valueStr = val
					}

					if f.Kind() == reflect.Int {
						if val, err := strconv.Atoi(valueStr); err == nil {
							//if !f.OverflowInt(int64(val)) { //  потом
							f.SetInt(int64(val))
							//}
						}
					} else if f.Kind() == reflect.String {
						f.SetString(valueStr)
					}
				}
			}
		}
	}
}

// читаем файл построчно
func readFileSecret() map[string]string {
	// rootDir - тут локальный
	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		rootDir = "."
	}
	secretfile := filepath.Join(rootDir, "secret.env")
	file, err := os.Open(secretfile)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	mapEnv := map[string]string{}
	for scanner.Scan() {
		setValueToStruct(scanner.Text(), mapEnv)
	}

	if err := scanner.Err(); err != nil {
		return nil
	}
	return mapEnv
}

func setValueToStruct(secret string, m map[string]string) {
	secret = strings.Trim(secret, " ")
	// отбрасываем пустые
	if secret == "" {
		return
	}
	// отбрасываем комментарии
	if strings.HasPrefix(secret, "#") {
		return
	}
	// отбрасываем коментарии в строке inline
	arr1 := strings.SplitN(secret, "#", 2)
	str1 := arr1[0]
	// получаем данные, должно быть 2 штуки
	arr2 := strings.SplitN(str1, "=", 2)
	if len(arr2) < 2 {
		return
	}
	// заполняем map
	m[strings.Trim(arr2[0], " ")] = strings.Trim(arr2[1], " ")
}

func PrintEnvironment(s interface{}) {
	//v := reflect.ValueOf(s)         //  здесь указатель
	//v = reflect.Indirect(v)         // получим значение по указателю v, если это не указатель вернет v
	el := reflect.ValueOf(s).Elem() // возвращает значение которое содержит интерфейс
	if el.Kind() == reflect.Struct {
		for i := 0; i < el.NumField(); i++ { // NumField кол-во полей, если не структура то паника
			name := el.Type().Field(i).Name
			f := el.FieldByName(name)
			if f.IsValid() {
				log.Printf("% 20s | %s", name, el.Field(i).Interface())
			}
		}
	}
}

/* example execute
type secretEnv struct {
	IsSecret  string
	APP_ID    string
	APP_HASH  string
	BOT_TOKEN string
}
var SecretEnv *secretEnv
func init(){
	SecretEnv = &secretEnv{}
}
jt.LoadEnvironment(SecretEnv)

//?Можно выводить информацию при некоторых условиях
if SecretEnv.IsSecret != "" {
	log.Println("Загружены данные из secret.env")
	if SecretEnv.IsSecret == "1" {
		jt.PrintEnvironment(SecretEnv)
	}
}

secret.env
APP_ID=WSFSFASFASEF
IsSecret=0

*/
