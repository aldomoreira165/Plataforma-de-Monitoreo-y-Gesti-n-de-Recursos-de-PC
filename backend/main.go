package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Structs
type Ram struct {
	Libre float64 `json:"libre"`
	Usada float64 `json:"usada"`
}

type Cpu struct {
	Libre float64 `json:"libre"`
	Usada float64 `json:"usada"`
}

type RamHistorico struct {
	Porcentaje float64   `json:"porcentaje"`
	Tiempo     time.Time `json:"tiempo"`
}

type CPUHistorico struct {
	Porcentaje float64   `json:"porcentaje"`
	Tiempo     time.Time `json:"tiempo"`
}

type RamHistoricoTiempo struct {
	Porcentaje float64 `json:"porcentaje"`
	Tiempo     string  `json:"tiempo"`
}

type CPUHistoricoTiempo struct {
	Porcentaje float64 `json:"porcentaje"`
	Tiempo     string  `json:"tiempo"`
}

// credenciales necesarias para la conexión a la base de datos
const (
	DBHost     = "mySQL"
	DBPort     = "3306"
	DBUser     = "root"
	DBPassword = "root"
	DBName     = "proyectoSO1"
)

// Funciones
func ObtenerDataRam() (Ram, error) {
	cmd := exec.Command("sh", "-c", "cat /proc/ram_so1_1s2024")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return Ram{}, err
	}

	valores := strings.Split(string(out), ",")
	libre, _ := strconv.ParseFloat(strings.TrimSpace(valores[0]), 64)
	usada, _ := strconv.ParseFloat(strings.TrimSpace(valores[1]), 64)

	return Ram{
		Libre: libre,
		Usada: usada,
	}, nil
}

func ObtenerDataCpu() (Cpu, error) {
	cmd := exec.Command("mpstat")
	sumaPorcentaje := 0.0

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return Cpu{}, err
	}

	output := string(out)
	re := regexp.MustCompile(`\d+\.\d+`) // Encuentra todos los números decimales en la cadena
	matches := re.FindAllString(output, -1)

	//recorrer el array de matches y sumar los valores en la posicion 1, 3, 7, 8 y 9
	for i := 1; i < len(matches); i++ {
		if i == 1 || i == 3 || i == 7 || i == 8 || i == 9 {
			valor, _ := strconv.ParseFloat(matches[i], 64)
			sumaPorcentaje += valor
		}
	}

	sumaPorcentaje = math.Round(sumaPorcentaje*100) / 100
	usada := sumaPorcentaje
	libre := 100 - usada

	return Cpu{
		Libre: libre,
		Usada: usada,
	}, nil
}

func InsertarUsoHistoricoRAM(porcentaje float64, tiempo time.Time, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO uso_ram (porcentaje_usado, tiempo) VALUES (?, ?)", porcentaje, tiempo)
	if err != nil {
		return err
	}
	return nil
}

// Función para insertar datos de uso de CPU en la tabla uso_historico_cpu
func InsertarUsoHistoricoCPU(porcentaje float64, tiempo time.Time, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO uso_cpu (porcentaje_usado, tiempo) VALUES (?, ?)", porcentaje, tiempo)
	if err != nil {
		return err
	}
	return nil
}

// Middleware para agregar los encabezados CORS
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {

	//abriendo la base de datos
	db, err := sql.Open("mysql", DBUser+":"+DBPassword+"@tcp("+DBHost+":"+DBPort+")/"+DBName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//Handlers
	http.HandleFunc("/api/ram", func(w http.ResponseWriter, r *http.Request) {

		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		fmt.Println("Accediendo GET /api/ram")

		ram, err := ObtenerDataRam()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(ram)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	http.HandleFunc("/api/cpu", func(w http.ResponseWriter, r *http.Request) {
		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		fmt.Println("Accediendo GET /api/cpu")

		cpu, err := ObtenerDataCpu()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(cpu)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	http.HandleFunc("/api/uso_historico_ram", func(w http.ResponseWriter, r *http.Request) {
		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		// Consultar la base de datos para obtener los datos de uso de RAM
		rows, err := db.Query("SELECT * FROM uso_ram")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Crear una estructura para almacenar los datos de uso de RAM
		var usosRam []RamHistoricoTiempo

		// Iterar sobre los resultados de la consulta y mapearlos a la estructura Ram
		for rows.Next() {
			var id int
			var tiempoBytes []byte
			var porcentajeUsado float64

			err := rows.Scan(&id, &porcentajeUsado, &tiempoBytes)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tiempoStr := string(tiempoBytes)
			usoRam := RamHistoricoTiempo{porcentajeUsado, tiempoStr}
			usosRam = append(usosRam, usoRam)
		}

		// Convertir los datos a formato JSON y escribirlos en la respuesta HTTP
		jsonData, err := json.Marshal(usosRam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	http.HandleFunc("/api/uso_historico_cpu", func(w http.ResponseWriter, r *http.Request) {
		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		fmt.Println("Accediendo GET /api/uso_historico_cpu")

		// Consultar la base de datos para obtener los datos de uso de CPU
		rows, err := db.Query("SELECT * FROM uso_cpu")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Crear una estructura para almacenar los datos de uso de RAM
		var usosCpu []CPUHistoricoTiempo

		// Iterar sobre los resultados de la consulta y mapearlos a la estructura Ram
		for rows.Next() {
			var id int
			var tiempoBytes []byte // Escaneamos el tiempo como []byte
			var porcentajeUsado float64

			err := rows.Scan(&id, &porcentajeUsado, &tiempoBytes)
			if err != nil {
				fmt.Println("Error escaneando fila:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tiempoStr := string(tiempoBytes)
			usoCpu := CPUHistoricoTiempo{porcentajeUsado, tiempoStr}
			usosCpu = append(usosCpu, usoCpu)
		}

		fmt.Println("Usos de cpu:", usosCpu)

		// Convertir los datos a formato JSON y escribirlos en la respuesta HTTP
		jsonData, err := json.Marshal(usosCpu)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	http.HandleFunc("/api/insertar_uso_ram", func(w http.ResponseWriter, r *http.Request) {
		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		fmt.Println("Accediendo POST /insertar_uso_ram")

		// Parsear el cuerpo de la solicitud JSON
		var data RamHistorico
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Insertar datos de uso de RAM en la base de datos
		err = InsertarUsoHistoricoRAM(data.Porcentaje, data.Tiempo, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	})

	http.HandleFunc("/api/insertar_uso_cpu", func(w http.ResponseWriter, r *http.Request) {
		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		fmt.Println("Accediendo POST /insertar_uso_cpu")

		// Parsear el cuerpo de la solicitud JSON
		var data CPUHistorico
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Insertar datos de uso de CPU en la base de datos
		err = InsertarUsoHistoricoCPU(data.Porcentaje, data.Tiempo, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	// eliminar datos de la tabla uso_ram y uso_cpu
	http.HandleFunc("/api/eliminar_uso_historico", func(w http.ResponseWriter, r *http.Request) {
		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		// Eliminar datos de la tabla uso_ram
		_, err := db.Exec("DELETE FROM uso_ram")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Eliminar datos de la tabla uso_cpu
		_, err = db.Exec("DELETE FROM uso_cpu")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/start_process", func(w http.ResponseWriter, r *http.Request) {
		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		var process *exec.Cmd

		cmd := exec.Command("sleep", "infinity")
		err := cmd.Start()
		if err != nil {
			fmt.Print(err)
			http.Error(w, "Error al iniciar el proceso", http.StatusInternalServerError)
			return
		}

		//obteniendo el pid del proceso
		process = cmd
		pid := process.Process.Pid
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, pid)
	})

	http.HandleFunc("/api/kill_process", func(w http.ResponseWriter, r *http.Request) {
		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		pidStr := r.URL.Query().Get("pid")
		if pidStr == "" {
			http.Error(w, "Se requiere el parámetro 'pid'", http.StatusBadRequest)
			return
		}

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			http.Error(w, "El parámetro 'pid' debe ser un número entero", http.StatusBadRequest)
			return
		}

		// Enviar señal SIGCONT al proceso con el PID proporcionado
		cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
		err = cmd.Run()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error al intentar terminar el proceso con PID %d", pid), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Proceso con PID %d ha terminado", pid)
	})

	http.HandleFunc("/api/resume_process", func(w http.ResponseWriter, r *http.Request) {
		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		pidStr := r.URL.Query().Get("pid")
		if pidStr == "" {
			http.Error(w, "Se requiere el parámetro 'pid'", http.StatusBadRequest)
			return
		}

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			http.Error(w, "El parámetro 'pid' debe ser un número entero", http.StatusBadRequest)
			return
		}

		// Enviar señal SIGCONT al proceso con el PID proporcionado
		cmd := exec.Command("kill", "-SIGCONT", strconv.Itoa(pid))
		err = cmd.Run()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error al reanudar el proceso con PID %d", pid), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Proceso con PID %d reanudado", pid)
	})

	http.HandleFunc("/api/stop_process", func(w http.ResponseWriter, r *http.Request) {
		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		pidStr := r.URL.Query().Get("pid")
		if pidStr == "" {
			http.Error(w, "Se requiere el parámetro 'pid'", http.StatusBadRequest)
			return
		}

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			http.Error(w, "El parámetro 'pid' debe ser un número entero", http.StatusBadRequest)
			return
		}

		// Enviar señal SIGSTOP al proceso con el PID proporcionado
		cmd := exec.Command("kill", "-SIGSTOP", strconv.Itoa(pid))
		err = cmd.Run()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error al detener el proceso con PID %d", pid), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Proceso con PID %d detenido", pid)
	})

	http.HandleFunc("/api/insertar_estado", func(w http.ResponseWriter, r *http.Request) {
		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		// Parsear el cuerpo de la solicitud JSON
		var estado struct {
			IDProceso int    `json:"id_proceso"`
			Estado    string `json:"estado"`
		}
		err := json.NewDecoder(r.Body).Decode(&estado)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validar los datos recibidos
		if estado.IDProceso == 0 || estado.Estado == "" {
			http.Error(w, "IDProceso y Estado son campos requeridos", http.StatusBadRequest)
			return
		}

		// Insertar datos en la tabla estado
		_, err = db.Exec("INSERT INTO estados (id_proceso, estado) VALUES (?, ?)", estado.IDProceso, estado.Estado)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	http.HandleFunc("/api/obtener_historial_estado", func(w http.ResponseWriter, r *http.Request) {
		// Agregar los encabezados CORS
		enableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		// Obtener el ID del proceso de la URL
		pidStr := r.URL.Query().Get("pid")
		if pidStr == "" {
			http.Error(w, "Se requiere el parámetro 'pid'", http.StatusBadRequest)
			return
		}

		// Convertir el ID del proceso a un número entero
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			http.Error(w, "El parámetro 'pid' debe ser un número entero", http.StatusBadRequest)
			return
		}

		// Consultar la base de datos para obtener el historial de estado del proceso
		rows, err := db.Query("SELECT id_proceso, estado FROM estados WHERE id_proceso = ?", pid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Crear una estructura para almacenar el historial de estado del proceso
		var historialEstado []struct {
			IDProceso int    `json:"id_proceso"`
			Estado    string `json:"estado"`
		}

		// Iterar sobre los resultados de la consulta y mapearlos a la estructura
		for rows.Next() {
			var idProceso int
			var estado string
			err := rows.Scan(&idProceso, &estado)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			historialEstado = append(historialEstado, struct {
				IDProceso int    `json:"id_proceso"`
				Estado    string `json:"estado"`
			}{idProceso, estado})
		}

		// Convertir los datos a formato JSON y escribirlos en la respuesta HTTP
		jsonData, err := json.Marshal(historialEstado)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	// Iniciar el servidor
	server := http.Server{
		Addr: ":8080",
	}

	error := server.ListenAndServe()

	if error != nil {
		panic(error)
	}
}
