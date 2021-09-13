// Загружает конфедициальные данные из файла secret.env (например при разработке)
// файл должен находиться рядом с исполняемым файлом, наличие файла не проверяется и не вызывает ошибку.
// Загружает данные из OS Environment (переменные окружения) перезаписывая загруженные из файла.
// Необходимо подготовить простую структуру с требуемыми полями, поля должны быть type string и экспортируемыми
// функция jt.LoadEnvironment(s *struct) принимает параметр ссылку на структуру
// в которую будут вставленны данные если они будут найдены
package jt

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"reflect"
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
func LoadEnvironment(s interface{}) {
	mapEnv := readFileSecret()
	v := reflect.ValueOf(s) //  здесь указатель
	v = reflect.Indirect(v) // получим значение
	typeOfS := v.Type()
	el := reflect.ValueOf(s).Elem() // возвращает значение которое содержит интерфейс
	if el.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {

			f := el.FieldByName(typeOfS.Field(i).Name)
			if f.IsValid() {
				// Значение можно изменить, только если оно
				// адресуемо и не был получено
				// с использованием неэкспортированных полей структуры. т.е. с большой буквы
				if f.CanSet() {
					// Изменить значение
					if mapEnv != nil {
						if val, ok := mapEnv[typeOfS.Field(i).Name]; ok {
							f.SetString(val)
						}
					}
					val, exists := os.LookupEnv(typeOfS.Field(i).Name) // environment
					if exists {
						f.SetString(val)
					}
				}
			}

			// log.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
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
				log.Printf("%s\t\t: %s", name, el.Field(i).Interface())
			}
			// log.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
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
//? Можно выводить информацию при некоторых условиях
if SecretEnv.IsSecret != "" {
	log.Println("Загружены данные из secret.env")
	if SecretEnv.IsSecret == "1" {
		jt.PrintEnvironment(SecretEnv)
	}
}
*/
